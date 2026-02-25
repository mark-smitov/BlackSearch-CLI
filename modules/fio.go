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

type OfDataResponse struct {
	Meta struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
	Data json.RawMessage `json:"data"`
}

func ofDataRequest(endpoint, params string) ([]map[string]any, error) {
	baseURL := fmt.Sprintf("https://api.ofdata.ru%s?key=%s&%s", endpoint, os.Getenv("OFDATA_KEY"), params)
	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ofResp OfDataResponse
	if err := json.Unmarshal(body, &ofResp); err != nil {
		return nil, err
	}
	if ofResp.Meta.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", ofResp.Meta.Msg)
	}

	var list []map[string]any
	if err := json.Unmarshal(ofResp.Data, &list); err == nil {
		return filterRecords(list), nil
	}

	var single map[string]any
	if err := json.Unmarshal(ofResp.Data, &single); err == nil && len(single) > 0 {
		return filterRecords([]map[string]any{single}), nil
	}

	return nil, fmt.Errorf("данные не найдены")
}

func filterRecords(records []map[string]any) []map[string]any {
	var result []map[string]any
	for _, rec := range records {
		filtered := make(map[string]any)
		for k, v := range rec {
			if v != nil {
				filtered[k] = v
			}
		}
		if len(filtered) > 0 {
			result = append(result, filtered)
		}
	}
	return result
}

func printRecordsTree(records []map[string]any) {
	for i, rec := range records {
		if i > 0 {
			Divider()
		}
		PrintMapTree(rec)
	}
}

func RunFIOSearch() {
	fio := AskInput("Введите ФИО")
	if fio == "" {
		PrintError("ФИО не может быть пустым")
		return
	}

	sp := NewSpinner("поиск")
	sp.Start()
	encoded := url.QueryEscape(fio)
	records, err := ofDataRequest("/v2/search", "by=founder-name&obj=org&query="+encoded)
	sp.Stop()

	if err != nil {
		PrintError(err.Error())
		return
	}
	if len(records) == 0 {
		PrintWarn("Данные не найдены: " + fio)
		return
	}

	printRecordsTree(records)
}

func RunINNSearch() {
	inn := AskInput("Введите ИНН")
	if inn == "" || !IsDigits(inn) {
		PrintError("ИНН должен содержать только цифры")
		return
	}

	sp := NewSpinner("поиск")
	sp.Start()
	records, err := ofDataRequest("/v2/person", "inn="+inn)
	sp.Stop()

	if err != nil {
		PrintError(err.Error())
		return
	}
	if len(records) == 0 {
		PrintWarn("Данные не найдены: " + inn)
		return
	}

	printRecordsTree(records)
}

func RunOGRNSearch() {
	ogrn := AskInput("Введите ОГРН")
	if ogrn == "" || !IsDigits(ogrn) {
		PrintError("ОГРН должен содержать только цифры")
		return
	}

	sp := NewSpinner("поиск")
	sp.Start()
	records, err := ofDataRequest("/v2/inspections", "ogrn="+ogrn)
	sp.Stop()

	if err != nil {
		PrintError(err.Error())
		return
	}
	if len(records) == 0 {
		PrintWarn("Данные не найдены: " + ogrn)
		return
	}

	printRecordsTree(records)
}

func RunCompanySearch() {
	name := AskInput("Введите название компании")
	if name == "" {
		PrintError("Название не может быть пустым")
		return
	}

	sp := NewSpinner("поиск")
	sp.Start()
	encoded := url.QueryEscape(name)
	records, err := ofDataRequest("/v2/search", "by=name&obj=org&query="+encoded)
	sp.Stop()

	if err != nil {
		PrintError(err.Error())
		return
	}
	if len(records) == 0 {
		PrintWarn("Данные не найдены: " + name)
		return
	}

	printRecordsTree(records)
}

func RunGoogleDork() {
	fio := AskInput("Введите ФИО, имя или запрос")
	if fio == "" {
		PrintError("Запрос не может быть пустым")
		return
	}

	queries := map[string]string{
		"Суд":       fmt.Sprintf(`"%s" (site:sudrf.ru OR site:kad.arbitr.ru OR site:bsr.sudrf.ru)`, fio),
		"Реестры":   fmt.Sprintf(`"%s" (site:egrul.nalog.ru OR site:rosreestr.ru)`, fio),
		"Налоги":    fmt.Sprintf(`"%s" (site:nalog.ru OR site:fns.gov.ru)`, fio),
		"Соцсети":   fmt.Sprintf(`"%s" (site:vk.com OR site:facebook.com OR site:instagram.com)`, fio),
		"Медиа":     fmt.Sprintf(`"%s" (site:rbc.ru OR site:ria.ru OR site:tass.ru OR site:kommersant.ru)`, fio),
		"Профили":   fmt.Sprintf(`"%s" (site:linkedin.com OR site:hh.ru OR site:habr.com)`, fio),
		"Документы": fmt.Sprintf(`"%s" filetype:pdf OR filetype:doc OR filetype:xlsx`, fio),
		"Должности": fmt.Sprintf(`"%s" ("генеральный директор" OR "руководитель" OR "депутат" OR "учредитель")`, fio),
	}

	categories := []string{"Суд", "Реестры", "Налоги", "Соцсети", "Медиа", "Профили", "Документы", "Должности"}
	for i, category := range categories {
		last := i == len(categories)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		query := queries[category]
		encoded := url.QueryEscape(query)
		gURL := "https://www.google.com/search?q=" + encoded
		fmt.Printf("  %s%s\n      %s\n",
			StyleDim.Render(prefix),
			StyleKey.Render(category),
			StyleVal.Render(gURL),
		)
		if !last {
			fmt.Println()
		}
	}

	filename := "dorks_" + strings.ReplaceAll(fio, " ", "_") + ".txt"
	var sb strings.Builder
	sb.WriteString("Google Dorks: " + fio + "\n\n")
	for _, cat := range categories {
		query := queries[cat]
		encoded := url.QueryEscape(query)
		gURL := "https://www.google.com/search?q=" + encoded
		sb.WriteString(cat + ":\n" + gURL + "\n\n")
	}

	if err := WriteFile(filename, sb.String()); err == nil {
		PrintSuccess("Сохранено: " + filename)
	}
}

func RunBINSearch() {
	bin := AskInput("Введите BIN (6-8 цифр)")
	if bin == "" || !IsDigits(bin) || len(bin) < 6 || len(bin) > 8 {
		PrintError("Некорректный BIN: 6-8 цифр")
		return
	}

	sp := NewSpinner("bin")
	sp.Start()
	data, err := FetchJSON(
		"https://lookup.binlist.net/"+bin,
		map[string]string{"Accept-Version": "3", "User-Agent": "Mozilla/5.0"},
		10*time.Second,
	)
	sp.Stop()

	if err != nil {
		PrintError("Ошибка: " + err.Error())
		return
	}

	items := []struct{ key, label string }{
		{"scheme", "scheme"},
		{"type", "type"},
		{"brand", "brand"},
		{"prepaid", "prepaid"},
	}
	for _, f := range items {
		if v, ok := data[f.key]; ok && v != nil {
			fmt.Printf("  ├── %s: %s\n", StyleKey.Render(f.label), StyleVal.Render(fmt.Sprintf("%v", v)))
		}
	}

	if bank, ok := data["bank"].(map[string]any); ok {
		fmt.Printf("  ├── %s\n", StyleKey.Render("bank"))
		bankFields := []struct{ key, label string }{
			{"name", "name"},
			{"url", "url"},
			{"phone", "phone"},
			{"city", "city"},
		}
		for i, k := range bankFields {
			last := i == len(bankFields)-1
			p := "│   ├── "
			if last {
				p = "│   └── "
			}
			if v, ok := bank[k.key]; ok && v != nil {
				fmt.Printf("  %s%s: %s\n", StyleDim.Render(p), StyleKey.Render(k.label), StyleVal.Render(fmt.Sprintf("%v", v)))
			}
		}
	}

	if country, ok := data["country"].(map[string]any); ok {
		fmt.Printf("  └── %s\n", StyleKey.Render("country"))
		cFields := []struct{ key, label string }{
			{"name", "name"},
			{"alpha2", "alpha2"},
			{"currency", "currency"},
			{"numeric", "numeric"},
		}
		for i, k := range cFields {
			last := i == len(cFields)-1
			p := "    ├── "
			if last {
				p = "    └── "
			}
			if v, ok := country[k.key]; ok && v != nil {
				fmt.Printf("  %s%s: %s\n", StyleDim.Render(p), StyleKey.Render(k.label), StyleVal.Render(fmt.Sprintf("%v", v)))
			}
		}
	}
}
