package modules

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type PhoneResult struct {
	Number    string
	Ninjas    map[string]any
	Numverify map[string]any
	Abstract  map[string]any
}

func cleanPhone(phone string) string {
	re := regexp.MustCompile(`[^\d+]`)
	phone = re.ReplaceAllString(phone, "")
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}
	return phone
}

func phoneNinjas(phone string) (map[string]any, error) {
	url := fmt.Sprintf("https://api.api-ninjas.com/v1/validatephone?number=%s", phone)
	return FetchJSON(url, map[string]string{"X-Api-Key": os.Getenv("API_NINJAS_KEY")}, 10*time.Second)
}

func phoneNumverify(phone string) (map[string]any, error) {
	url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s", os.Getenv("NUMVERIFY_KEY"), phone)
	return FetchJSON(url, nil, 10*time.Second)
}

func phoneAbstract(phone string) (map[string]any, error) {
	url := fmt.Sprintf("https://phonevalidation.abstractapi.com/v1/?api_key=%s&phone=%s", os.Getenv("ABSTRACT_PHONE_KEY"), phone)
	return FetchJSON(url, nil, 10*time.Second)
}

func RunPhoneAnalysis() {
	raw := AskInput("Введите номер телефона")
	if raw == "" {
		PrintError("Номер не введён")
		return
	}

	phone := cleanPhone(raw)

	result := PhoneResult{Number: phone}
	var mu sync.Mutex
	var wg sync.WaitGroup

	pb := NewProgressBar("телефон", 3)
	pb.Render()

	wg.Add(3)

	go func() {
		defer wg.Done()
		data, err := phoneNinjas(phone)
		mu.Lock()
		if err == nil {
			result.Ninjas = data
		}
		mu.Unlock()
		pb.Inc()
	}()

	go func() {
		defer wg.Done()
		data, err := phoneNumverify(phone)
		mu.Lock()
		if err == nil {
			result.Numverify = data
		}
		mu.Unlock()
		pb.Inc()
	}()

	go func() {
		defer wg.Done()
		data, err := phoneAbstract(phone)
		mu.Lock()
		if err == nil {
			result.Abstract = data
		}
		mu.Unlock()
		pb.Inc()
	}()

	wg.Wait()
	pb.Done()

	merged := mergePhoneSources(result)
	printMergedTree(merged)

	locQuery := extractLocation(result.Ninjas, result.Numverify, result.Abstract)
	if locQuery != "" {
		geoURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
			strings.ReplaceAll(locQuery, " ", "+"))
		geo, err := FetchJSON(geoURL, map[string]string{"User-Agent": "BlackSearch/4.0"}, 8*time.Second)
		if err == nil {
			printMergedTree(geo)
		}
	}
}

func mergePhoneSources(r PhoneResult) map[string]any {
	merged := make(map[string]any)
	for _, src := range []map[string]any{r.Ninjas, r.Numverify, r.Abstract} {
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
	return merged
}

func printMergedTree(data map[string]any) {
	keys := sortedKeys(data)
	valid := filterValidKeys(keys, data)
	for i, k := range valid {
		v := data[k]
		val := fmt.Sprintf("%v", v)
		last := i == len(valid)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n",
			StyleDim.Render(prefix),
			StyleKey.Render(k),
			StyleVal.Render(TruncateStr(val, 120)),
		)
	}
}

func extractLocation(sources ...map[string]any) string {
	for _, src := range sources {
		if src == nil {
			continue
		}
		if loc, ok := src["location"].(string); ok && loc != "" {
			return loc
		}
		if country, ok := src["country_name"].(string); ok && country != "" {
			if loc, ok := src["location"].(string); ok {
				return loc + ", " + country
			}
			return country
		}
	}
	return ""
}
