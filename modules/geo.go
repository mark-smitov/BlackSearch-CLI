package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func geoFetch(rawURL string, headers map[string]string) ([]byte, error) {
	client := &http.Client{Timeout: 12 * time.Second}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; BlackSearch/4.0)")
	req.Header.Set("Accept", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d от %s", resp.StatusCode, req.Host)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("пустой ответ")
	}
	return body, nil
}

func RunGeoSearch() {
	fmt.Printf("\n  %s\n  %s\n",
		StyleInfo.Render("1. Адрес → Координаты"),
		StyleInfo.Render("2. Координаты → Адрес"),
	)

	choice := AskInput("Выберите режим (1/2)")

	switch choice {
	case "1":
		geoForward()
	case "2":
		geoReverse()
	default:
		PrintError("Неверный выбор")
	}
}

// ─── Адрес → Координаты ─────────────────────────────────────────────────────

func geoForward() {
	address := AskInput("Введите адрес")
	if address == "" {
		PrintError("Адрес не введён")
		return
	}

	sp := NewSpinner("гео")
	sp.Start()

	q := url.QueryEscape(address)

	// 1. Nominatim
	body, err := geoFetch(
		"https://nominatim.openstreetmap.org/search?q="+q+"&format=json&limit=5&addressdetails=1",
		map[string]string{"Accept-Language": "ru,en"},
	)
	if err == nil {
		var arr []map[string]any
		if json.Unmarshal(body, &arr) == nil && len(arr) > 0 {
			sp.Stop()
			PrintInfo("(nominatim.openstreetmap.org)")
			printGeoForwardResults(arr)
			return
		}
	}
	nominatimErr := err

	// 2. Photon (komoot)
	body, err = geoFetch(
		"https://photon.komoot.io/api/?q="+q+"&limit=5&lang=ru",
		nil,
	)
	if err == nil {
		var fc map[string]any
		if json.Unmarshal(body, &fc) == nil {
			if features, ok := fc["features"].([]any); ok && len(features) > 0 {
				sp.Stop()
				PrintInfo("(photon.komoot.io)")
				printPhotonForwardResults(features)
				return
			}
		}
	}
	photonErr := err

	// 3. geoapify (бесплатный ключ — можно добавить в .env)
	apiKey := getEnvFallback("GEOAPIFY_KEY", "")
	if apiKey != "" {
		body, err = geoFetch(
			"https://api.geoapify.com/v1/geocode/search?text="+q+"&limit=5&apiKey="+apiKey,
			nil,
		)
		if err == nil {
			var fc map[string]any
			if json.Unmarshal(body, &fc) == nil {
				if features, ok := fc["features"].([]any); ok && len(features) > 0 {
					sp.Stop()
					PrintInfo("(api.geoapify.com)")
					printGeoapifyForwardResults(features)
					return
				}
			}
		}
	}

	sp.Stop()

	// Показываем реальные ошибки вместо "ничего не найдено"
	if nominatimErr != nil {
		PrintWarn("Nominatim: " + nominatimErr.Error())
	} else {
		PrintWarn("Nominatim: результатов не найдено")
	}
	if photonErr != nil {
		PrintWarn("Photon:    " + photonErr.Error())
	} else {
		PrintWarn("Photon:    результатов не найдено")
	}
	if apiKey == "" {
		PrintInfo("Совет: добавьте GEOAPIFY_KEY= в .env (бесплатно на geoapify.com)")
	}
}

func printGeoForwardResults(arr []map[string]any) {
	for i, item := range arr {
		if i > 0 {
			Divider()
		}
		var items []struct{ k, v string }
		for _, key := range []string{"display_name", "lat", "lon", "type", "class"} {
			if v, ok := item[key].(string); ok && v != "" {
				items = append(items, struct{ k, v string }{key, v})
			}
		}
		lat, _ := item["lat"].(string)
		lon, _ := item["lon"].(string)
		if lat != "" && lon != "" {
			items = append(items, struct{ k, v string }{
				"google_maps", "https://www.google.com/maps?q=" + lat + "," + lon,
			})
		}
		printKVItems(items)
	}
}

func printPhotonForwardResults(features []any) {
	for i, f := range features {
		if i > 0 {
			Divider()
		}
		feat, ok := f.(map[string]any)
		if !ok {
			continue
		}
		props, _ := feat["properties"].(map[string]any)
		var lat, lon string
		if geom, ok := feat["geometry"].(map[string]any); ok {
			if coords, ok := geom["coordinates"].([]any); ok && len(coords) == 2 {
				lon = fmt.Sprintf("%v", coords[0])
				lat = fmt.Sprintf("%v", coords[1])
			}
		}

		var items []struct{ k, v string }
		// Собираем display_name вручную
		nameParts := []string{}
		for _, key := range []string{"name", "street", "housenumber", "city", "state", "country"} {
			if props != nil {
				if v, ok := props[key].(string); ok && v != "" {
					nameParts = append(nameParts, v)
				}
			}
		}
		if len(nameParts) > 0 {
			items = append(items, struct{ k, v string }{"display_name", strings.Join(nameParts, ", ")})
		}
		if lat != "" {
			items = append(items, struct{ k, v string }{"lat", lat})
			items = append(items, struct{ k, v string }{"lon", lon})
			items = append(items, struct{ k, v string }{
				"google_maps", "https://www.google.com/maps?q=" + lat + "," + lon,
			})
		}
		printKVItems(items)
	}
}

func printGeoapifyForwardResults(features []any) {
	for i, f := range features {
		if i > 0 {
			Divider()
		}
		feat, ok := f.(map[string]any)
		if !ok {
			continue
		}
		props, _ := feat["properties"].(map[string]any)
		var items []struct{ k, v string }
		for _, key := range []string{"formatted", "lat", "lon", "country", "city", "postcode"} {
			if props != nil {
				if v, ok := props[key]; ok {
					items = append(items, struct{ k, v string }{key, fmt.Sprintf("%v", v)})
				}
			}
		}
		if lat, ok := props["lat"]; ok {
			if lon, ok := props["lon"]; ok {
				items = append(items, struct{ k, v string }{
					"google_maps", fmt.Sprintf("https://www.google.com/maps?q=%v,%v", lat, lon),
				})
			}
		}
		printKVItems(items)
	}
}

// ─── Координаты → Адрес ─────────────────────────────────────────────────────

func geoReverse() {
	lat := AskInput("Широта (latitude)")
	lon := AskInput("Долгота (longitude)")
	if lat == "" || lon == "" {
		PrintError("Координаты не введены")
		return
	}

	sp := NewSpinner("гео")
	sp.Start()

	// 1. BigDataCloud — полностью бесплатно, без ключа
	body, err := geoFetch(
		fmt.Sprintf("https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=ru",
			url.QueryEscape(lat), url.QueryEscape(lon)),
		nil,
	)
	if err == nil {
		var data map[string]any
		if json.Unmarshal(body, &data) == nil && data["city"] != nil {
			sp.Stop()
			PrintInfo("(bigdatacloud.net)")
			printReverseResult(lat, lon, data)
			return
		}
	}
	bdcErr := err

	// 2. Nominatim
	body, err = geoFetch(
		fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json&addressdetails=1&accept-language=ru",
			url.QueryEscape(lat), url.QueryEscape(lon)),
		map[string]string{"Accept-Language": "ru,en"},
	)
	if err == nil {
		var data map[string]any
		if json.Unmarshal(body, &data) == nil {
			sp.Stop()
			PrintInfo("(nominatim.openstreetmap.org)")
			var items []struct{ k, v string }
			if v, ok := data["display_name"].(string); ok {
				items = append(items, struct{ k, v string }{"display_name", v})
			}
			items = append(items, struct{ k, v string }{
				"google_maps", "https://www.google.com/maps?q=" + lat + "," + lon,
			})
			items = append(items, struct{ k, v string }{
				"openstreetmap", "https://www.openstreetmap.org/?mlat=" + lat + "&mlon=" + lon,
			})
			if addr, ok := data["address"].(map[string]any); ok {
				for _, key := range []string{"country", "state", "county", "city", "town", "village", "suburb", "road", "house_number", "postcode"} {
					if v, ok := addr[key]; ok {
						items = append(items, struct{ k, v string }{key, fmt.Sprintf("%v", v)})
					}
				}
			}
			printKVItems(items)
			return
		}
	}
	nominatimErr := err

	sp.Stop()
	if bdcErr != nil {
		PrintWarn("BigDataCloud: " + bdcErr.Error())
	}
	if nominatimErr != nil {
		PrintWarn("Nominatim:    " + nominatimErr.Error())
	} else {
		PrintWarn("Ничего не найдено по этим координатам")
	}
}

func printReverseResult(lat, lon string, data map[string]any) {
	var items []struct{ k, v string }
	for _, key := range []string{"locality", "city", "principalSubdivision", "countryName", "postcode"} {
		if v, ok := data[key].(string); ok && v != "" {
			items = append(items, struct{ k, v string }{key, v})
		}
	}
	if v, ok := data["localityInfo"].(map[string]any); ok {
		if admin, ok := v["administrative"].([]any); ok && len(admin) > 0 {
			if a, ok := admin[0].(map[string]any); ok {
				if name, ok := a["name"].(string); ok {
					items = append(items, struct{ k, v string }{"region", name})
				}
			}
		}
	}
	items = append(items, struct{ k, v string }{
		"google_maps", "https://www.google.com/maps?q=" + lat + "," + lon,
	})
	items = append(items, struct{ k, v string }{
		"openstreetmap", "https://www.openstreetmap.org/?mlat=" + lat + "&mlon=" + lon,
	})
	printKVItems(items)
}

func printKVItems(items []struct{ k, v string }) {
	for i, item := range items {
		last := i == len(items)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(item.k), StyleVal.Render(item.v))
	}
}

func getEnvFallback(key, fallback string) string {
	if v := getEnv(key); v != "" {
		return v
	}
	return fallback
}

func getEnv(key string) string {
	
	
	return strings.TrimSpace(os.Getenv(key))
}
