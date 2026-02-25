package progress_bar

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	StyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("#696969")).Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Bold(true)
	StyleWarn    = lipgloss.NewStyle().Foreground(lipgloss.Color("#707070")).Bold(true)
	StyleInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd"))
	StyleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("rgb(119, 119, 121)"))
	StyleBold    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Bold(true)
	StyleKey     = lipgloss.NewStyle().Foreground(lipgloss.Color("#7c7c7c"))
	StyleVal     = lipgloss.NewStyle().Foreground(lipgloss.Color("#999999"))
	StyleAccent  = lipgloss.NewStyle().Foreground(lipgloss.Color("#444344")).Bold(true)
	StyleGreen   = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
)

var stdinReader = bufio.NewReader(os.Stdin)

type ProgressBar struct {
	label   string
	total   int
	current int
}

func NewProgressBar(label string, total int) *ProgressBar {
	return &ProgressBar{label: label, total: total}
}

func (p *ProgressBar) Set(current int) {
	p.current = current
	p.render()
}

func (p *ProgressBar) Inc() {
	p.current++
	if p.current > p.total {
		p.current = p.total
	}
	p.render()
}

func (p *ProgressBar) Done() {
	p.current = p.total
	p.render()
	fmt.Println()
}

func (p *ProgressBar) Render() {
	p.render()
}

func (p *ProgressBar) render() {
	pct := 0.0
	if p.total > 0 {
		pct = float64(p.current) / float64(p.total)
	}
	pctStr := fmt.Sprintf("%3.0f%%", pct*100)
	fmt.Printf("\r  %s  %s  %d/%d   ",
		StyleDim.Render(p.label),
		StyleBold.Render(pctStr),
		p.current, p.total,
	)
}

type Spinner struct {
	label  string
	done   chan struct{}
	doneCh chan struct{}
}

func NewSpinner(label string) *Spinner {
	return &Spinner{
		label:  label,
		done:   make(chan struct{}),
		doneCh: make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.done:
				fmt.Printf("\r%s\r", strings.Repeat(" ", 60))
				close(s.doneCh)
				return
			default:
				fmt.Printf("\r  %s %s", StyleDim.Render(frames[i%len(frames)]), StyleDim.Render(s.label))
				time.Sleep(80 * time.Millisecond)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop() {
	close(s.done)
	<-s.doneCh
}

func PrintSuccess(msg string) {
	fmt.Println(StyleSuccess.Render("  ✓ ") + StyleBold.Render(msg))
}

func PrintError(msg string) {
	fmt.Println(StyleError.Render("  ✗ ") + StyleBold.Render(msg))
}

func PrintWarn(msg string) {
	fmt.Println(StyleWarn.Render("  · ") + StyleDim.Render(msg))
}

func PrintInfo(msg string) {
	fmt.Println(StyleInfo.Render("  ➤ ") + msg)
}

func PrintKV(key, val string) {
	fmt.Printf("  %s %s\n", StyleKey.Render(key+":"), StyleVal.Render(val))
}

func PrintSection(title string) {
	fmt.Printf("\n  %s %s\n\n", StyleAccent.Render("·"), StyleDim.Render(title))
}

func PrintHeader(title string) {
	fmt.Printf("\n  %s\n", StyleBold.Render(title))
}

func PrintTree(items []TreeItem) {
	for i, item := range items {
		last := i == len(items)-1
		prefix := "├── "
		if last {
			prefix = "└── "
		}
		fmt.Printf("  %s%s: %s\n", StyleDim.Render(prefix), StyleKey.Render(item.Key), StyleVal.Render(item.Value))
		for j, child := range item.Children {
			childLast := j == len(item.Children)-1
			childPrefix := "│   ├── "
			if last {
				childPrefix = "    ├── "
			}
			if childLast {
				if last {
					childPrefix = "    └── "
				} else {
					childPrefix = "│   └── "
				}
			}
			fmt.Printf("  %s%s: %s\n",
				StyleDim.Render(childPrefix),
				StyleKey.Render(child.Key),
				StyleVal.Render(child.Value),
			)
		}
	}
}

type TreeItem struct {
	Key      string
	Value    string
	Children []TreeItem
}

func AskInput(prompt string) string {
	fmt.Print(lipgloss.NewStyle().Foreground(lipgloss.Color("#737279")).Bold(true).Render("\n  [➤] "+prompt+": "))
	input, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(input)
}

func Divider() {
	fmt.Println(StyleDim.Render("  " + strings.Repeat("·", 40)))
}

func PrintDim(s string) {
	fmt.Println(StyleDim.Render(s))
}
