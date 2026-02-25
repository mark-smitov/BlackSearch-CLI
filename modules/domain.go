package modules

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func cleanDomain(input string) string {
	d := strings.ToLower(input)
	for _, prefix := range []string{"https://", "http://"} {
		d = strings.TrimPrefix(d, prefix)
	}
	if idx := strings.Index(d, "/"); idx != -1 {
		d = d[:idx]
	}
	return strings.TrimSpace(d)
}

func RunDomainAnalysis() {
	input := AskInput("Введите домен или URL")
	if input == "" || !strings.Contains(input, ".") {
		PrintError("Некорректный домен")
		return
	}

	domain := cleanDomain(input)

	type DomainResult struct {
		BaseDomain string
		Whois      map[string]string
		DNS        map[string][]string
		IP         string
		GeoData    map[string]any
	}
	result := DomainResult{BaseDomain: domain, DNS: make(map[string][]string), Whois: make(map[string]string)}

	pb := NewProgressBar("домен", 5)
	pb.Render()

	raw, err := whois.Whois(domain)
	if err == nil {
		parsed, err := whoisparser.Parse(raw)
		if err == nil {
			if parsed.Registrar.Name != "" {
				result.Whois["registrar"] = parsed.Registrar.Name
			}
			if parsed.Domain.CreatedDate != "" {
				result.Whois["created"] = parsed.Domain.CreatedDate
			}
			if parsed.Domain.ExpirationDate != "" {
				result.Whois["expires"] = parsed.Domain.ExpirationDate
			}
			if parsed.Domain.UpdatedDate != "" {
				result.Whois["updated"] = parsed.Domain.UpdatedDate
			}
			if parsed.Administrative.Email != "" {
				result.Whois["admin_email"] = parsed.Administrative.Email
			}
			if parsed.Registrant.Name != "" {
				result.Whois["registrant"] = parsed.Registrant.Name
			}
			if parsed.Registrant.Organization != "" {
				result.Whois["organization"] = parsed.Registrant.Organization
			}
		}
	}
	pb.Inc()

	var mu sync.Mutex
	var wg sync.WaitGroup

	dnsTypes := []string{"A", "AAAA", "MX", "NS", "TXT"}
	for _, qt := range dnsTypes {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			records := resolveDNS(domain, t)
			mu.Lock()
			if len(records) > 0 {
				result.DNS[t] = records
			}
			mu.Unlock()
		}(qt)
	}
	wg.Wait()
	pb.Inc()
	pb.Inc()
	pb.Inc()

	ips, err := net.LookupHost(domain)
	if err == nil && len(ips) > 0 {
		result.IP = ips[0]
		geo, err := FetchJSON(
			fmt.Sprintf("https://ipinfo.io/%s/json", result.IP),
			map[string]string{"User-Agent": "BlackSearch/4.0"},
			8*time.Second,
		)
		if err == nil {
			result.GeoData = geo
		}
	}
	pb.Inc()
	pb.Done()

	merged := make(map[string]any)
	merged["domain"] = domain
	if result.IP != "" {
		merged["ip"] = result.IP
	}
	for k, v := range result.Whois {
		merged[k] = v
	}
	for t, records := range result.DNS {
		merged["dns_"+strings.ToLower(t)] = strings.Join(records, " | ")
	}
	for _, key := range []string{"country", "region", "city", "org", "timezone", "loc"} {
		if v, ok := result.GeoData[key]; ok {
			merged["geo_"+key] = fmt.Sprintf("%v", v)
		}
	}

	printMergedTree(merged)
}

func printKVList(items []struct{ k, v string }) {
	for i, item := range items {
		last := i == len(items)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(item.k), StyleVal.Render(item.v))
	}
}
