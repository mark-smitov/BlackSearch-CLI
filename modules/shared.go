package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	pb "blacksearch/progress_bar"
)

var (
	StyleSuccess = pb.StyleSuccess
	StyleError   = pb.StyleError
	StyleWarn    = pb.StyleWarn
	StyleInfo    = pb.StyleInfo
	StyleDim     = pb.StyleDim
	StyleBold    = pb.StyleBold
	StyleKey     = pb.StyleKey
	StyleVal     = pb.StyleVal
)

func PrintSuccess(msg string)       { pb.PrintSuccess(msg) }
func PrintError(msg string)         { pb.PrintError(msg) }
func PrintWarn(msg string)          { pb.PrintWarn(msg) }
func PrintInfo(msg string)          { pb.PrintInfo(msg) }
func PrintKV(key, val string)       { pb.PrintKV(key, val) }
func PrintSection(title string)     { pb.PrintSection(title) }
func PrintHeader(title string)      { pb.PrintHeader(title) }
func AskInput(prompt string) string { return pb.AskInput(prompt) }
func Divider()                      { pb.Divider() }
func PrintDim(s string)             { pb.PrintDim(s) }

func NewSpinner(label string) *pb.Spinner {
	return pb.NewSpinner(label)
}

func NewProgressBar(label string, total int) *pb.ProgressBar {
	return pb.NewProgressBar(label, total)
}

func PrintTree(items []pb.TreeItem) { pb.PrintTree(items) }

func FetchJSONArray(url string, headers map[string]string, timeout time.Duration) ([]map[string]any, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body[:min(len(body), 120)])))
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("пустой ответ от сервера")
	}
	if body[0] == '<' {
		return nil, fmt.Errorf("сервер вернул HTML вместо JSON (возможно rate limit или блокировка)")
	}
	var result []map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FetchJSON(url string, headers map[string]string, timeout time.Duration) (map[string]any, error) {
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func PrintMapTree(data map[string]any) {
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

func filterValidKeys(keys []string, data map[string]any) []string {
	var result []string
	for _, k := range keys {
		v := data[k]
		str := fmt.Sprintf("%v", v)
		if str == "" || str == "<nil>" || str == "0" || str == "false" || str == "[]" || str == "map[]" {
			continue
		}
		result = append(result, k)
	}
	return result
}

func PrintMapTreeNested(data map[string]any, depth int) {
	keys := sortedKeys(data)
	indent := strings.Repeat("│   ", depth)

	for i, k := range keys {
		v := data[k]
		last := i == len(keys)-1
		prefix := indent + "├── "
		childIndent := indent + "│   "
		if last {
			prefix = indent + "└── "
			childIndent = indent + "    "
		}

		switch val := v.(type) {
		case map[string]any:
			fmt.Printf("  %s%s\n", StyleDim.Render(prefix), StyleKey.Render(k))
			printNestedMap(val, childIndent)
		case []any:
			if len(val) == 0 {
				continue
			}
			fmt.Printf("  %s%s: [%d]\n", StyleDim.Render(prefix), StyleKey.Render(k), len(val))
			for j, item := range val {
				if j >= 5 {
					fmt.Printf("  %s%s\n", StyleDim.Render(childIndent+"└── "), StyleDim.Render(fmt.Sprintf("... и ещё %d", len(val)-5)))
					break
				}
				if m, ok := item.(map[string]any); ok {
					printNestedMap(m, childIndent)
				}
			}
		default:
			str := fmt.Sprintf("%v", val)
			if str == "" || str == "<nil>" || str == "0" || str == "false" {
				continue
			}
			fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(k), StyleVal.Render(TruncateStr(str, 100)))
		}
	}
}

func printNestedMap(data map[string]any, indent string) {
	keys := sortedKeys(data)
	for i, k := range keys {
		v := data[k]
		last := i == len(keys)-1
		prefix := indent + "├── "
		if last {
			prefix = indent + "└── "
		}
		str := fmt.Sprintf("%v", v)
		if str == "" || str == "<nil>" || str == "0" || str == "false" {
			continue
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(k), StyleVal.Render(TruncateStr(str, 100)))
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TruncateStr(s string, n int) string {
	return truncateStr(s, n)
}

func truncateStr(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "…"
}

func IsDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func WriteFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
