package modules

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type LocalSearchResult struct {
	File string
	Line string
}

func searchFile(path, query string) []LocalSearchResult {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var results []LocalSearchResult
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".json" {
		var data any
		dec := json.NewDecoder(f)
		if err := dec.Decode(&data); err != nil {
			return nil
		}
		raw, _ := json.Marshal(data)
		rawStr := string(raw)
		if strings.Contains(strings.ToLower(rawStr), strings.ToLower(query)) {
			if arr, ok := data.([]any); ok {
				for _, item := range arr {
					b, _ := json.Marshal(item)
					if strings.Contains(strings.ToLower(string(b)), strings.ToLower(query)) {
						results = append(results, LocalSearchResult{File: path, Line: string(b)})
					}
				}
			} else {
				results = append(results, LocalSearchResult{File: path, Line: rawStr})
			}
		}
		return results
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			results = append(results, LocalSearchResult{
				File: path,
				Line: fmt.Sprintf("%d: %s", lineNum, line),
			})
		}
	}
	return results
}

func formatBytes(b int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	val := float64(b)
	i := 0
	for val >= 1024 && i < len(units)-1 {
		val /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d %s", b, units[i])
	}
	return fmt.Sprintf("%.2f %s", val, units[i])
}

func RunLocalDBSearch() {
	dir := AskInput("Путь к папке с базой данных")
	if dir == "" {
		PrintError("Путь не указан")
		return
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		PrintError("Папка не существует: " + dir)
		return
	}

	query := AskInput("Строка для поиска")
	if query == "" {
		PrintError("Запрос не может быть пустым")
		return
	}

	var totalFiles int64
	var totalSize int64
	var filePaths []string
	allowedExts := map[string]bool{".csv": true, ".txt": true, ".json": true, "": true}

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if allowedExts[ext] {
			filePaths = append(filePaths, path)
			totalSize += info.Size()
			totalFiles++
		}
		return nil
	})

	fmt.Printf("  ├── %s: %d\n", StyleKey.Render("files"), totalFiles)
	fmt.Printf("  ├── %s: %s\n", StyleKey.Render("size"), StyleVal.Render(formatBytes(totalSize)))
	fmt.Printf("  └── %s: %s\n\n", StyleKey.Render("query"), StyleVal.Render(query))

	pb := NewProgressBar("поиск", int(totalFiles))
	pb.Render()

	var mu sync.Mutex
	var wg sync.WaitGroup
	var allResults []LocalSearchResult
	var processed int64
	sem := make(chan struct{}, 8)

	for _, path := range filePaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			results := searchFile(p, query)
			mu.Lock()
			allResults = append(allResults, results...)
			mu.Unlock()

			atomic.AddInt64(&processed, 1)
			pb.Set(int(atomic.LoadInt64(&processed)))
		}(path)
	}

	wg.Wait()
	pb.Done()

	if len(allResults) == 0 {
		PrintWarn("Ничего не найдено")
		return
	}

	byFile := make(map[string][]string)
	for _, r := range allResults {
		byFile[r.File] = append(byFile[r.File], r.Line)
	}

	for file, lines := range byFile {
		fmt.Printf("\n  %s %s\n", StyleDim.Render("·"), StyleKey.Render(filepath.Base(file)))
		PrintDim("    " + file)
		for i, line := range lines {
			if i >= 20 {
				PrintWarn(fmt.Sprintf("    ... и ещё %d совпадений", len(lines)-20))
				break
			}
			highlighted := highlightQuery(line, query)
			fmt.Println("    " + highlighted)
		}
	}

	PrintSuccess(fmt.Sprintf("Найдено: %d совпадений", len(allResults)))
}

func highlightQuery(line, query string) string {
	lower := strings.ToLower(line)
	lowerQ := strings.ToLower(query)
	idx := strings.Index(lower, lowerQ)
	if idx == -1 {
		return StyleDim.Render(TruncateStr(line, 150))
	}
	before := line[:idx]
	match := line[idx : idx+len(query)]
	after := line[idx+len(query):]
	return StyleDim.Render(TruncateStr(before, 50)) +
		StyleError.Render(match) +
		StyleDim.Render(TruncateStr(after, 80))
}
