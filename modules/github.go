package modules

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func RunGitHubSearch() {
	fmt.Printf("\n  %s\n  %s\n  %s\n",
		StyleInfo.Render("1. Пользователь / Организация"),
		StyleInfo.Render("2. Репозитории"),
		StyleInfo.Render("3. Поиск по коду"),
	)

	choice := AskInput("Выберите (1-3)")
	token := os.Getenv("GITHUB_TOKEN")

	headers := map[string]string{
		"User-Agent": "BlackSearch/4.0",
		"Accept":     "application/vnd.github.v3+json",
	}
	if token != "" {
		headers["Authorization"] = "Bearer " + token
	} else {
		PrintWarn("Нет GITHUB_TOKEN — лимит 60 req/h")
	}

	var query string
	var apiURL string

	switch choice {
	case "1":
		query = AskInput("Логин пользователя или организации")
		if query == "" {
			PrintError("Логин не введён")
			return
		}
		apiURL = fmt.Sprintf("https://api.github.com/users/%s", query)
	case "2":
		query = AskInput("Поисковый запрос")
		if query == "" {
			PrintError("Запрос не введён")
			return
		}
		apiURL = fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&per_page=15", strings.ReplaceAll(query, " ", "+"))
	case "3":
		query = AskInput("Поиск по коду")
		if query == "" {
			PrintError("Запрос не введён")
			return
		}
		apiURL = fmt.Sprintf("https://api.github.com/search/code?q=%s&per_page=10", strings.ReplaceAll(query, " ", "+"))
	default:
		PrintError("Неверный выбор")
		return
	}

	sp := NewSpinner("github")
	sp.Start()
	data, err := FetchJSON(apiURL, headers, 12*time.Second)
	sp.Stop()

	if err != nil {
		PrintError("Ошибка: " + err.Error())
		return
	}

	if errMsg, ok := data["message"].(string); ok && errMsg != "" {
		PrintError("GitHub API: " + errMsg)
		return
	}

	switch choice {
	case "1":
		printGitHubUser(data)
	case "2":
		printGitHubRepos(data)
	case "3":
		PrintMapTreeNested(data, 0)
	}
}

func printGitHubUser(data map[string]any) {
	fields := []struct{ key, label string }{
		{"name", "name"},
		{"login", "login"},
		{"type", "type"},
		{"company", "company"},
		{"location", "location"},
		{"email", "email"},
		{"blog", "blog"},
		{"bio", "bio"},
		{"public_repos", "public_repos"},
		{"public_gists", "public_gists"},
		{"followers", "followers"},
		{"following", "following"},
		{"created_at", "created_at"},
		{"updated_at", "updated_at"},
		{"html_url", "html_url"},
		{"avatar_url", "avatar_url"},
	}

	var items []struct{ k, v string }
	for _, f := range fields {
		if v, ok := data[f.key]; ok && v != nil {
			str := fmt.Sprintf("%v", v)
			if str != "" && str != "0" && str != "<nil>" {
				items = append(items, struct{ k, v string }{f.label, str})
			}
		}
	}

	for i, item := range items {
		last := i == len(items)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(item.k), StyleVal.Render(TruncateStr(item.v, 100)))
	}
}

func printGitHubRepos(data map[string]any) {
	count, _ := data["total_count"].(float64)
	fmt.Printf("  %s%s: %s\n\n",
		StyleDim.Render("┌── "),
		StyleKey.Render("total"),
		StyleVal.Render(fmt.Sprintf("%.0f", count)),
	)

	items, ok := data["items"].([]any)
	if !ok || len(items) == 0 {
		PrintWarn("Репозитории не найдены")
		return
	}

	for i, item := range items {
		repo, ok := item.(map[string]any)
		if !ok {
			continue
		}
		last := i == len(items)-1
		connector := "├"
		if last {
			connector = "└"
		}
		name := fmt.Sprintf("%v", repo["full_name"])
		stars := fmt.Sprintf("★ %.0f", func() float64 { v, _ := repo["stargazers_count"].(float64); return v }())
		lang := fmt.Sprintf("%v", repo["language"])
		if lang == "<nil>" {
			lang = "—"
		}
		desc := ""
		if d, ok := repo["description"].(string); ok && d != "" {
			desc = " · " + TruncateStr(d, 60)
		}
		url := fmt.Sprintf("%v", repo["html_url"])

		fmt.Printf("  %s── %s  %s  [%s]%s\n",
			StyleDim.Render(connector),
			StyleBold.Render(name),
			StyleSuccess.Render(stars),
			StyleInfo.Render(lang),
			StyleDim.Render(desc),
		)
		fmt.Printf("  %s   %s\n",
			StyleDim.Render("│"),
			StyleVal.Render(url),
		)
		if !last {
			fmt.Println(StyleDim.Render("  │"))
		}
	}
}
