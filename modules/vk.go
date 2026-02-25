package modules

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func RunVKSearch() {
	target := AskInput("Введите ID, ник или ссылку vk.com/...")
	if target == "" {
		PrintError("ID не введён")
		return
	}

	target = strings.TrimPrefix(target, "https://vk.com/")
	target = strings.TrimPrefix(target, "http://vk.com/")
	target = strings.TrimPrefix(target, "vk.com/")
	target = strings.TrimSpace(target)

	sp := NewSpinner("vk")
	sp.Start()

	token := os.Getenv("VK_TOKEN")
	var url string
	if token != "" {
		url = fmt.Sprintf(
			"https://api.vk.com/method/users.get?user_ids=%s&fields=photo_max_orig,city,bdate,status,counters,contacts,education,universities,schools,career&v=5.131&access_token=%s",
			target, token,
		)
	} else {
		PrintWarn("Нет VK_TOKEN — counters и закрытые поля недоступны")
		url = fmt.Sprintf(
			"https://api.vk.com/method/users.get?user_ids=%s&fields=photo_max_orig,city,bdate,status,contacts,education,universities,schools,career&v=5.131",
			target,
		)
	}
	data, err := FetchJSON(url, nil, 10*time.Second)
	sp.Stop()

	if err != nil {
		PrintError("Ошибка запроса: " + err.Error())
		return
	}

	if errMsg, ok := data["error"].(map[string]any); ok {
		PrintError(fmt.Sprintf("VK API: %v", errMsg["error_msg"]))
		return
	}

	response, ok := data["response"].([]any)
	if !ok || len(response) == 0 {
		PrintWarn("Пользователь не найден или профиль скрыт")
		return
	}

	user, ok := response[0].(map[string]any)
	if !ok {
		PrintError("Неверный формат ответа")
		return
	}

	name := fmt.Sprintf("%v %v", user["first_name"], user["last_name"])
	profileURL := fmt.Sprintf("https://vk.com/id%v", user["id"])

	items := []struct{ k, v string }{
		{"name", name},
		{"id", fmt.Sprintf("%v", user["id"])},
		{"url", profileURL},
	}

	if city, ok := user["city"].(map[string]any); ok {
		items = append(items, struct{ k, v string }{"city", fmt.Sprintf("%v", city["title"])})
	}
	if bdate, ok := user["bdate"].(string); ok && bdate != "" {
		items = append(items, struct{ k, v string }{"bdate", bdate})
	}
	if status, ok := user["status"].(string); ok && status != "" {
		items = append(items, struct{ k, v string }{"status", status})
	}
	if photo, ok := user["photo_max_orig"].(string); ok && photo != "" {
		items = append(items, struct{ k, v string }{"photo", photo})
	}
	if deactivated, ok := user["deactivated"].(string); ok && deactivated != "" {
		items = append(items, struct{ k, v string }{"deactivated", deactivated})
	}
	if closed, ok := user["is_closed"].(bool); ok {
		val := "open"
		if closed {
			val = "closed"
		}
		items = append(items, struct{ k, v string }{"profile", val})
	}

	hasCounters := user["counters"] != nil
	for i, item := range items {
		last := i == len(items)-1 && !hasCounters
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(item.k), StyleVal.Render(item.v))
	}

	if counters, ok := user["counters"].(map[string]any); ok {
		fmt.Printf("  └── %s\n", StyleKey.Render("counters"))
		ckeys := sortedKeys(counters)
		for i, k := range ckeys {
			last := i == len(ckeys)-1
			prefix := "    ├── "
			if last {
				prefix = "    └── "
			}
			val := fmt.Sprintf("%v", counters[k])
			if val == "0" || val == "<nil>" {
				continue
			}
			fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(k), StyleVal.Render(val))
		}
	}
}
