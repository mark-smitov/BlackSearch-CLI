package modules

import (
	"fmt"
	"net"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func getMXRecords(domain string) []string {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil
	}
	var result []string
	for _, mx := range mxRecords {
		result = append(result, strings.TrimSuffix(mx.Host, "."))
	}
	return result
}

func smtpCheck(email, mailServer string) bool {
	conn, err := net.DialTimeout("tcp", mailServer+":25", 6*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, mailServer)
	if err != nil {
		return false
	}
	defer c.Close()
	if err = c.Hello("check.example.com"); err != nil {
		return false
	}
	if err = c.Mail("check@example.com"); err != nil {
		return false
	}
	if err = c.Rcpt(email); err != nil {
		return false
	}
	return true
}

func RunEmailAnalysis() {
	email := AskInput("Введите email")
	if email == "" {
		PrintError("Email не введён")
		return
	}
	if !emailRegex.MatchString(email) {
		PrintError("Некорректный формат email: " + email)
		return
	}

	domain := strings.Split(email, "@")[1]

	pb := NewProgressBar("email", 4)
	pb.Render()

	mxRecords := getMXRecords(domain)
	pb.Inc()

	type EmailResult struct {
		Email     string
		MXRecords []string
		SMTPOK    bool
		Sources   map[string]map[string]any
	}
	result := EmailResult{
		Email:     email,
		MXRecords: mxRecords,
		Sources:   make(map[string]map[string]any),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if len(result.MXRecords) > 0 {
			ok := smtpCheck(email, result.MXRecords[0])
			mu.Lock()
			result.SMTPOK = ok
			mu.Unlock()
		}
		pb.Inc()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://api.api-ninjas.com/v1/validateemail?email=%s", email)
		data, err := FetchJSON(url, map[string]string{"X-Api-Key": os.Getenv("API_NINJAS_KEY")}, 10*time.Second)
		mu.Lock()
		if err == nil {
			result.Sources["ninjas"] = data
		}
		mu.Unlock()
		pb.Inc()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		url := fmt.Sprintf("https://emailvalidation.abstractapi.com/v1/?api_key=%s&email=%s", os.Getenv("ABSTRACT_EMAIL_KEY"), email)
		data, err := FetchJSON(url, nil, 10*time.Second)
		mu.Lock()
		if err == nil {
			result.Sources["abstract"] = data
		}
		mu.Unlock()
		pb.Inc()
	}()

	wg.Wait()
	pb.Done()

	merged := make(map[string]any)
	merged["email"] = email
	merged["domain"] = domain
	if len(mxRecords) > 0 {
		merged["mx"] = strings.Join(mxRecords, " | ")
	}
	smtpStr := "нет"
	if result.SMTPOK {
		smtpStr = "да"
	}
	merged["smtp_ok"] = smtpStr

	for _, src := range []map[string]any{result.Sources["ninjas"], result.Sources["abstract"]} {
		if src == nil {
			continue
		}
		for k, v := range src {
			if _, exists := merged[k]; !exists {
				str := fmt.Sprintf("%v", v)
				if str != "" && str != "<nil>" && str != "false" && str != "0" {
					merged[k] = v
				}
			}
		}
	}

	printMergedTree(merged)
}
