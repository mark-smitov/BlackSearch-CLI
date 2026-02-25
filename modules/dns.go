package modules

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func resolveDNS(domain, qtype string) []string {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second
	m := new(dns.Msg)

	typeMap := map[string]uint16{
		"A":     dns.TypeA,
		"MX":    dns.TypeMX,
		"NS":    dns.TypeNS,
		"TXT":   dns.TypeTXT,
		"AAAA":  dns.TypeAAAA,
		"CNAME": dns.TypeCNAME,
		"SOA":   dns.TypeSOA,
	}

	t, ok := typeMap[qtype]
	if !ok {
		return nil
	}

	m.SetQuestion(dns.Fqdn(domain), t)
	r, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		r, _, err = c.Exchange(m, "1.1.1.1:53")
		if err != nil {
			return nil
		}
	}

	var results []string
	for _, ans := range r.Answer {
		switch v := ans.(type) {
		case *dns.A:
			results = append(results, v.A.String())
		case *dns.AAAA:
			results = append(results, v.AAAA.String())
		case *dns.MX:
			results = append(results, fmt.Sprintf("%d %s", v.Preference, strings.TrimSuffix(v.Mx, ".")))
		case *dns.NS:
			results = append(results, strings.TrimSuffix(v.Ns, "."))
		case *dns.TXT:
			results = append(results, strings.Join(v.Txt, " "))
		case *dns.CNAME:
			results = append(results, strings.TrimSuffix(v.Target, "."))
		case *dns.SOA:
			results = append(results, fmt.Sprintf("%s | %s | Serial:%d", strings.TrimSuffix(v.Ns, "."), strings.TrimSuffix(v.Mbox, "."), v.Serial))
		}
	}
	return results
}

func RunDNSAnalysis() {
	domain := AskInput("Введите домен")
	if domain == "" {
		PrintError("Домен не введён")
		return
	}
	domain = cleanDomain(domain)

	types := []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA"}

	pb := NewProgressBar("dns", len(types))
	pb.Render()

	var mu sync.Mutex
	var wg sync.WaitGroup
	results := make(map[string][]string)

	for _, qt := range types {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			records := resolveDNS(domain, t)
			mu.Lock()
			if len(records) > 0 {
				results[t] = records
			}
			mu.Unlock()
			pb.Inc()
		}(qt)
	}
	wg.Wait()
	pb.Done()

	if len(results) == 0 {
		PrintWarn("DNS записи не найдены")
		return
	}

	orderedTypes := []string{"A", "AAAA", "MX", "NS", "TXT", "CNAME", "SOA"}
	type entry struct {
		t       string
		records []string
	}
	var found []entry
	for _, t := range orderedTypes {
		if records, ok := results[t]; ok {
			found = append(found, entry{t, records})
		}
	}

	for i, item := range found {
		last := i == len(found)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		if len(item.records) == 1 {
			fmt.Printf("  %s%s: %s\n",
				StyleDim.Render(prefix),
				StyleKey.Render(fmt.Sprintf("%-6s", item.t)),
				StyleVal.Render(item.records[0]),
			)
		} else {
			fmt.Printf("  %s%s: [%d]\n",
				StyleDim.Render(prefix),
				StyleKey.Render(fmt.Sprintf("%-6s", item.t)),
				len(item.records),
			)
			childPrefix := "│   "
			if last {
				childPrefix = "    "
			}
			for j, rec := range item.records {
				childLast := j == len(item.records)-1
				recPrefix := childPrefix + "├── "
				if childLast {
					recPrefix = childPrefix + "└── "
				}
				fmt.Printf("  %s%s\n", StyleDim.Render(recPrefix), StyleVal.Render(rec))
			}
		}
	}
}
