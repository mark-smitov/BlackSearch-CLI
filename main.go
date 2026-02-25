package main

import (
	"blacksearch/modules"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

var (
	colorGray   = lipgloss.Color("#6272a4")
	colorBlue   = lipgloss.Color("#7abeda")
	colorWhite  = lipgloss.Color("#f8f8f2")
	colorGreen  = lipgloss.Color("#c5dac5")
	colorPurple = lipgloss.Color("#c2d4da")
	colorDim    = lipgloss.Color("#6b7cee")
	colorRed    = lipgloss.Color("#ff5555")
)

func main() {
	_ = godotenv.Load()
	for {
		clearScreen()
		showInterface()

		choice := readInput()

		switch choice {
		case "1":
			modules.RunPhoneAnalysis()
		case "2":
			modules.RunIPAnalysis()
		case "3":
			modules.RunEmailAnalysis()
		case "4":
			modules.RunDomainAnalysis()
		case "5":
			modules.RunDNSAnalysis()
		case "6":
			modules.RunNickSearch()
		case "7":
			modules.RunFIOSearch()
		case "8":
			modules.RunINNSearch()
		case "9":
			modules.RunOGRNSearch()
		case "10":
			modules.RunCompanySearch()
		case "11":
			modules.RunLocalDBSearch()
		case "12":
			modules.RunVKSearch()
		case "13":
			modules.RunGitHubSearch()
		case "14":
			modules.RunGoogleDork()
		case "15":
			modules.RunBINSearch()
		case "16":
			modules.RunGeoSearch()
		case "17":
			modules.RunGraphServer()
		case "18":
			modules.RunDossier()
		case "19":
			modules.RunReport()
		case "0", "q", "exit":
			bye()
			os.Exit(0)
		default:
			modules.PrintWarn("Неверный номер команды. Нажмите Enter...")
		}

		if choice != "0" {
			fmt.Print(lipgloss.NewStyle().Foreground(colorDim).Render("\n  Enter — в меню"))
			fmt.Scanln()
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func readInput() string {
	promptStyle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	fmt.Print(promptStyle.Render("\n [➤] Команда: "))
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func showInterface() {
	bannerText := `
    ____  __            __  _____                      __       _    ____ __
   / __ )/ /___ ______/ /__/ ___/___  ____ ___________/ /_     | |  / / // /
  / __  / / __  / ___/ //_/\__ \/ _ \/ __  / ___/ ___/ __ \    | | / / // /_
 / /_/ / / /_/ / /__/ ,<  ___/ /  __/ /_/ / /  / /__/ / / /    | |/ /__  __/
/_____/_/\__,_/\___/_/|_|/____/\___/\__,_/_/   \___/_/ /_/     |___/  /_/  `

	justiceArt := `
	    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣤⣤⣤⣤⣘⣿⣀⣤⣤⣤⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢠⣤⣤⣤⣤⣤⣴⡾⠟⠋⠀⠈⠈⠉⢻⣿⡿⠉⠈⠀⠉⠉⠻⢷⣦⣤⣤⣤⣤⣤⡄⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠸⣿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⡅⠀⠀⠀⠀⠀⠀⠀⠈⠈⢹⣿⡀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⢷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠠⣿⡅⠀⠀⠀⠀⠀⠀⠀⠀⠀⡾⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡟⢼⡉⣧⠀⠀⠀⠀⠀⠀⠀⠀⢈⣿⡆⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⣿⢹⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⡿⠀⢽⡀⠸⣇⠀⠀⠀⠀⠀⠀⠀⠨⣿⡆⠀⠀⠀⠀⠀⠀⠀⢰⠀⠀⣿⠀⢿⡀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⠃⠀⢽⠂⠀⢻⡄⠀⠀⠀⠀⠀⠀⠨⣿⡆⠀⠀⠀⠀⠀⠀⠢⡟⠀⠀⣿⠀⠘⣷⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢰⡽⠀⠀⢽⡂⠀⠀⢿⡀⠀⠀⠀⠀⠀⢠⣿⠆⠀⠀⠀⠀⠀⢀⡾⠀⠀⠀⣿⠀⠀⠹⣇⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⡿⠀⠀⠀⢽⡄⠀⠀⠈⣷⠀⠀⠀⠀⠀⣴⣿⣧⠀⠀⠀⠀⠀⣼⠃⠀⠀⠀⣿⠀⠀⠀⢻⡄⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣾⠃⠀⠀⠀⢽⠄⠀⠀⠀⠸⣇⠀⠀⠀⢸⣿⣿⣿⡇⠀⠀⠀⣸⠇⠀⠀⠀⠀⣿⠀⠀⠀⠈⣷⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢶⣤⣤⣄⣄⣀⣀⣉⣀⣀⣀⣄⣤⣤⣶⠆⠀⠸⣿⣿⣿⡇⠀⠰⣶⣤⣤⣄⣄⣀⣀⣉⣀⣀⣀⣤⣤⣤⡶⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠹⣿⡟⠀⠀⠀⠈⻣⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠙⻿⠿⣿⣿⣿⡿⠿⠛⠀⠀⠀⠀⠀⠀⠨⣿⡆⠀⠀⠀⠀⠀⠈⠛⠿⢿⣿⣿⣿⠿⠟⠋⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣾⣿⣿⣿⣷⣄⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠉⠀⠀`

	fmt.Println(lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render(bannerText))

	menuBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorGray).
		Padding(1, 2)

	artBorderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorGray).
		Foreground(colorBlue).
		Padding(0, 1)

	item := func(num, text string) string {
		n := lipgloss.NewStyle().Foreground(colorPurple).Bold(true).Render(num)
		t := lipgloss.NewStyle().Foreground(colorWhite).Render(text)
		return n + "  " + t
	}

	col1 := strings.Join([]string{
		item(" 1.", "Анализ Номера"),
		item(" 2.", "Анализ IP"),
		item(" 3.", "Анализ почты"),
		item(" 4.", "Анализ Сайта"),
		item(" 5.", "Анализ DNS"),
		item(" 6.", "Поиск по Нику"),
		item(" 7.", "Поиск по ФИО"),
		item(" 8.", "Поиск по ИНН"),
		item(" 9.", "Поиск по ОГРН"),
		item("10.", "Поиск по Компании"),
	}, "\n")

	col2 := strings.Join([]string{
		item("11.", "Поиск по Локал БД"),
		item("12.", "Поиск по ВК"),
		item("13.", "Поиск по Github"),
		item("14.", "Google Dork"),
		item("15.", "Поиск по карте"),
		item("16.", "Гео.Поиск"),
		item("17.", "Граф связей"),
		item("18.", "Создать досье"),
		item("19.", "Создать отчёт"),
		lipgloss.NewStyle().Foreground(colorRed).Bold(true).Render(" 0.  Выход"),
	}, "\n")

	menuContent := lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)
	menuBox := menuBorderStyle.Render(menuContent)
	artBox := artBorderStyle.Render(justiceArt)
	finalLayout := lipgloss.JoinHorizontal(lipgloss.Top, menuBox, " ", artBox)

	fmt.Println(finalLayout)
}

func bye() {
	fmt.Println(lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("\n  BlackSearch  ·  Bye!"))
}
