package modules

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type SiteCheck struct {
	Name string
	URL  string
}

var nickSites = []SiteCheck{
	{"VKontakte", "https://vk.com/{nick}"},
	{"Telegram", "https://t.me/{nick}"},
	{"GitHub", "https://github.com/{nick}"},
	{"GitLab", "https://gitlab.com/{nick}"},
	{"Twitter/X", "https://twitter.com/{nick}"},
	{"Instagram", "https://www.instagram.com/{nick}/"},
	{"TikTok", "https://www.tiktok.com/@{nick}"},
	{"YouTube", "https://www.youtube.com/@{nick}"},
	{"Reddit", "https://www.reddit.com/user/{nick}"},
	{"Pinterest", "https://www.pinterest.com/{nick}/"},
	{"Tumblr", "https://{nick}.tumblr.com/"},
	{"Medium", "https://medium.com/@{nick}"},
	{"Twitch", "https://www.twitch.tv/{nick}"},
	{"Kick", "https://kick.com/{nick}"},
	{"LinkedIn", "https://www.linkedin.com/in/{nick}"},
	{"NPM", "https://www.npmjs.com/~{nick}"},
	{"PyPI", "https://pypi.org/user/{nick}/"},
	{"HackerNews", "https://news.ycombinator.com/user?id={nick}"},
	{"StackOverflow", "https://stackoverflow.com/users/1/{nick}"},
	{"Habr", "https://habr.com/ru/users/{nick}/"},
	{"dev.to", "https://dev.to/{nick}"},
	{"Replit", "https://replit.com/@{nick}"},
	{"Codepen", "https://codepen.io/{nick}"},
	{"Steam", "https://steamcommunity.com/id/{nick}"},
	{"ProductHunt", "https://www.producthunt.com/@{nick}"},
	{"Dribbble", "https://dribbble.com/{nick}"},
	{"Behance", "https://www.behance.net/{nick}"},
	{"SoundCloud", "https://soundcloud.com/{nick}"},
}

type NickCheckResult struct {
	Site   string
	URL    string
	Found  bool
	Error  bool
	Status int
}

func checkNick(nick string, site SiteCheck, client *http.Client) NickCheckResult {
	rawURL := strings.ReplaceAll(site.URL, "{nick}", nick)
	result := NickCheckResult{Site: site.Name, URL: rawURL}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		result.Error = true
		return result
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		result.Error = true
		return result
	}
	defer resp.Body.Close()
	result.Status = resp.StatusCode

	if resp.StatusCode == 200 || resp.StatusCode == 301 || resp.StatusCode == 302 {
		result.Found = true
	}
	return result
}

func RunNickSearch() {
	nick := AskInput("Введите ник")
	if nick == "" {
		PrintError("Ник не введён")
		return
	}

	nick = strings.TrimPrefix(nick, "@")

	client := &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 2 {
				return fmt.Errorf("stop redirect")
			}
			return nil
		},
	}

	pb := NewProgressBar("ник", len(nickSites))
	pb.Render()

	var mu sync.Mutex
	var wg sync.WaitGroup
	results := make([]NickCheckResult, 0, len(nickSites))
	sem := make(chan struct{}, 15)

	for _, site := range nickSites {
		wg.Add(1)
		go func(s SiteCheck) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			r := checkNick(nick, s, client)
			mu.Lock()
			results = append(results, r)
			pb.Inc()
			mu.Unlock()
		}(site)
	}

	wg.Wait()
	pb.Done()

	var found []NickCheckResult
	for _, r := range results {
		if r.Found {
			found = append(found, r)
		}
	}

	sort.Slice(found, func(i, j int) bool { return found[i].Site < found[j].Site })

	if len(found) == 0 {
		PrintWarn("Аккаунты не найдены")
	} else {
		for i, r := range found {
			last := i == len(found)-1
			prefix := "├── "
			if last {
				prefix = "└── "
			}
			fmt.Printf("  %s%s %-20s %s\n",
				StyleDim.Render(prefix),
				StyleSuccess.Render("✓"),
				StyleBold.Render(r.Site),
				StyleVal.Render(r.URL),
			)
		}
	}

	filename := fmt.Sprintf("nick_%s_%s.txt", nick, time.Now().Format("20060102-1504"))
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("nick: %s\ndate: %s\n\n", nick, time.Now().Format("02.01.2006 15:04")))
	sb.WriteString(fmt.Sprintf("found (%d):\n", len(found)))
	for _, r := range found {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", r.Site, r.URL))
	}

	if err := WriteFile(filename, sb.String()); err == nil {
		PrintSuccess("Сохранено: " + filename)
	}
}
