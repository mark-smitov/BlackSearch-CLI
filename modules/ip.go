package modules

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func isPrivateIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	privateRanges := []string{
		"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16",
		"127.0.0.0/8", "::1/128", "fc00::/7",
	}
	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil && network.Contains(parsed) {
			return true
		}
	}
	return false
}

func RunIPAnalysis() {
	ip := AskInput("Введите IP-адрес")
	if ip == "" {
		PrintError("IP не введён")
		return
	}
	if net.ParseIP(ip) == nil {
		PrintError("Некорректный IP: " + ip)
		return
	}

	type sourceFunc struct {
		name string
		url  string
	}
	sources := []sourceFunc{
		{"ip-api.com", fmt.Sprintf("http://ip-api.com/json/%s?fields=66846719", ip)},
		{"ipinfo.io", fmt.Sprintf("https://ipinfo.io/%s/json", ip)},
		{"ipwho.is", fmt.Sprintf("https://ipwho.is/%s", ip)},
		{"ipwhois.app", fmt.Sprintf("https://ipwhois.app/json/%s", ip)},
		{"OTX AlienVault", fmt.Sprintf("https://otx.alienvault.com/api/v1/indicator/ip/%s/general", ip)},
		{"Abstract API", fmt.Sprintf("https://ipgeolocation.abstractapi.com/v1/?api_key=%s&ip_address=%s", os.Getenv("ABSTRACT_IP_KEY"), ip)},
	}

	pb := NewProgressBar("ip", len(sources))
	pb.Render()

	sourcesData := make(map[string]map[string]any)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)
		go func(s sourceFunc) {
			defer wg.Done()
			data, err := FetchJSON(s.url, map[string]string{"User-Agent": "BlackSearch/4.0"}, 8*time.Second)
			mu.Lock()
			if err == nil && len(data) > 0 {
				sourcesData[s.name] = data
			}
			mu.Unlock()
			pb.Inc()
		}(src)
	}

	wg.Wait()
	pb.Done()

	merged := mergeIPSources(sourcesData)
	if isPrivateIP(ip) {
		merged["type"] = "private"
	}
	merged["ip"] = ip
	printMergedTree(merged)
}

func mergeIPSources(sources map[string]map[string]any) map[string]any {
	priority := []string{"ip-api.com", "ipinfo.io", "ipwho.is", "ipwhois.app", "Abstract API", "OTX AlienVault"}
	merged := make(map[string]any)
	for _, name := range priority {
		src, ok := sources[name]
		if !ok {
			continue
		}
		for k, v := range src {
			if _, exists := merged[k]; !exists {
				str := fmt.Sprintf("%v", v)
				if str != "" && str != "<nil>" && str != "0" && str != "false" && str != "[]" && str != "map[]" {
					merged[k] = v
				}
			}
		}
	}
	return merged
}
