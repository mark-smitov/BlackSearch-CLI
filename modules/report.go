package modules

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func RunReport() {
	fmt.Printf("\n  %s\n  %s\n  %s\n  %s\n",
		StyleInfo.Render("1. HTML отчёт"),
		StyleInfo.Render("2. TXT отчёт"),
		StyleInfo.Render("3. MD отчёт"),
		StyleInfo.Render("4. PDF-ready HTML"),
	)

	format := AskInput("Выберите формат (1/2/3/4)")
	subject := AskInput("Тема / объект расследования")
	analyst := AskInput("Аналитик")

	PrintSection("Разделы отчёта")
	PrintInfo("Пустое название раздела = завершить")

	var sections []struct{ Title, Content string }

	for {
		title := AskInput("Название раздела")
		if title == "" || strings.ToLower(title) == "done" {
			break
		}
		content := AskInput("Содержимое раздела")
		sections = append(sections, struct{ Title, Content string }{title, content})
	}

	if len(sections) == 0 {
		PrintWarn("Отчёт пуст")
		return
	}

	now := time.Now()
	dateStr := now.Format("02.01.2006 15:04")
	id := now.Format("20060102-150405")

	sp := NewSpinner("отчёт")
	sp.Start()

	var filename string
	var content string
	var err error

	switch format {
	case "2":
		filename = fmt.Sprintf("report_%s.txt", id)
		content = buildTXTReport(subject, analyst, dateStr, sections)
	case "3":
		filename = fmt.Sprintf("report_%s.md", id)
		content = buildMDReport(subject, analyst, dateStr, sections)
	case "4":
		filename = fmt.Sprintf("report_%s_print.html", id)
		content = buildHTMLReport(subject, analyst, dateStr, id, sections, true)
	default:
		filename = fmt.Sprintf("report_%s.html", id)
		content = buildHTMLReport(subject, analyst, dateStr, id, sections, false)
	}

	err = os.WriteFile(filename, []byte(content), 0644)
	sp.Stop()

	if err != nil {
		PrintError("Ошибка записи: " + err.Error())
		return
	}

	PrintSuccess("Отчёт: " + filename)
}

func buildHTMLReport(subject, analyst, date, id string, sections []struct{ Title, Content string }, printMode bool) string {
	var sectionHTML strings.Builder
	for i, s := range sections {
		sectionHTML.WriteString(fmt.Sprintf(`
    <div class="section">
      <div class="section-num">%02d</div>
      <div class="section-body">
        <h2 class="section-title">%s</h2>
        <div class="section-content">%s</div>
      </div>
    </div>`, i+1, s.Title, strings.ReplaceAll(s.Content, "\n", "<br>")))
	}

	printCSS := ""
	if printMode {
		printCSS = `@media print { body { background: white !important; color: #222 !important; } .cover { background: #f5f5f5 !important; } }`
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Отчёт: %s</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;600;700&family=JetBrains+Mono:wght@400;500&display=swap');
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { background: #070707; color: #c8c8c8; font-family: 'Inter', sans-serif; min-height: 100vh; }
  .cover {
    min-height: 100vh; display: flex; flex-direction: column;
    justify-content: center; align-items: flex-start;
    padding: 80px 10%%; border-bottom: 1px solid #1a1a1a;
    position: relative; overflow: hidden;
  }
  .cover::before {
    content: ''; position: absolute; top: 0; left: 0; right: 0; bottom: 0;
    background: radial-gradient(ellipse at 80%% 50%%, #0d1a2e 0%%, transparent 60%%);
    pointer-events: none;
  }
  .cover-eyebrow { font-size: 10px; letter-spacing: 4px; color: #2a4a6a; text-transform: uppercase; margin-bottom: 24px; }
  .cover-title { font-size: 52px; font-weight: 700; letter-spacing: -2px; color: #fff; line-height: 1.1; margin-bottom: 16px; max-width: 720px; }
  .cover-sub { font-size: 16px; color: #333; margin-bottom: 56px; }
  .cover-meta { display: flex; gap: 48px; }
  .meta-item { display: flex; flex-direction: column; gap: 5px; }
  .meta-label { font-size: 9px; color: #2a2a2a; text-transform: uppercase; letter-spacing: 2px; }
  .meta-value { font-size: 14px; color: #777; }
  .report-id { position: absolute; top: 40px; right: 10%%; font-family: 'JetBrains Mono', monospace; font-size: 10px; color: #1a1a1a; letter-spacing: 2px; }
  .content { max-width: 760px; margin: 0 auto; padding: 80px 40px; }
  .section { display: flex; gap: 48px; margin-bottom: 64px; padding-bottom: 64px; border-bottom: 1px solid #0f0f0f; }
  .section:last-child { border-bottom: none; margin-bottom: 0; padding-bottom: 0; }
  .section-num { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #222; min-width: 36px; padding-top: 5px; }
  .section-body { flex: 1; }
  .section-title { font-size: 20px; font-weight: 600; color: #e0e0e0; margin-bottom: 20px; letter-spacing: -0.5px; }
  .section-content { font-size: 14px; color: #555; line-height: 1.9; font-weight: 300; }
  .footer { text-align: center; padding: 48px; border-top: 1px solid #0f0f0f; color: #1e1e1e; font-size: 10px; letter-spacing: 3px; text-transform: uppercase; }
  %s
</style>
</head>
<body>
<div class="cover">
  <div class="report-id">REP-%s</div>
  <div class="cover-eyebrow">Аналитический отчёт · BlackSearch Intelligence</div>
  <h1 class="cover-title">%s</h1>
  <div class="cover-sub">Содержит %d раздел(ов)</div>
  <div class="cover-meta">
    <div class="meta-item"><span class="meta-label">Аналитик</span><span class="meta-value">%s</span></div>
    <div class="meta-item"><span class="meta-label">Дата создания</span><span class="meta-value">%s</span></div>
    <div class="meta-item"><span class="meta-label">Статус</span><span class="meta-value">Конфиденциально</span></div>
  </div>
</div>
<div class="content">%s</div>
<div class="footer">BlackSearch Intelligence · %s</div>
</body>
</html>`,
		subject, printCSS, id, subject,
		len(sections), analyst, date,
		sectionHTML.String(), date)
}

func buildTXTReport(subject, analyst, date string, sections []struct{ Title, Content string }) string {
	var sb strings.Builder
	line := strings.Repeat("═", 62)

	sb.WriteString(line + "\n")
	sb.WriteString("  АНАЛИТИЧЕСКИЙ ОТЧЁТ · BLACKSEARCH INTELLIGENCE\n")
	sb.WriteString(line + "\n\n")
	sb.WriteString(fmt.Sprintf("  Тема:       %s\n", subject))
	sb.WriteString(fmt.Sprintf("  Аналитик:   %s\n", analyst))
	sb.WriteString(fmt.Sprintf("  Дата:       %s\n", date))
	sb.WriteString(fmt.Sprintf("  Разделов:   %d\n\n", len(sections)))
	sb.WriteString(line + "\n\n")

	for i, s := range sections {
		sb.WriteString(fmt.Sprintf("  [%02d] %s\n", i+1, strings.ToUpper(s.Title)))
		sb.WriteString("  " + strings.Repeat("─", 50) + "\n")
		sb.WriteString("  " + strings.ReplaceAll(s.Content, "\n", "\n  ") + "\n\n")
	}

	sb.WriteString(line + "\n")
	sb.WriteString("  КОНФИДЕНЦИАЛЬНО · BlackSearch Intelligence\n")
	return sb.String()
}

func buildMDReport(subject, analyst, date string, sections []struct{ Title, Content string }) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", subject))
	sb.WriteString(fmt.Sprintf("> **Аналитик:** %s  \n> **Дата:** %s  \n> **Источник:** BlackSearch Intelligence  \n> **Гриф:** Конфиденциально\n\n", analyst, date))
	sb.WriteString("---\n\n")

	for i, s := range sections {
		sb.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, s.Title))
		sb.WriteString(s.Content + "\n\n")
		sb.WriteString("---\n\n")
	}

	sb.WriteString("*Создано BlackSearch Intelligence*\n")
	return sb.String()
}
