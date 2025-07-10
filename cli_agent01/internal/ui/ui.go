package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	SuccessColor = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	BorderColor  = color.New(color.FgWhite).Add(color.Bold).SprintFunc()
	LabelColor   = color.New(color.FgHiWhite).Add(color.Bold).SprintFunc()
	AiColor      = color.New(color.FgHiCyan).Add(color.Bold).SprintFunc()
	ErrorColor   = color.New(color.FgRed).Add(color.Bold).SprintFunc()
)

func PrintHeader() {
	var asciiHeader = []string{
		"  ███████╗██╗     ██╗     █████╗  ██████╗  █████╗ ███████╗███╗   ██╗████████╗",
		"  ██╔════╝██║     ██║    ██╔══██╗██╔════╝ ██╔══██╗██╔════╝████╗  ██║╚══██╔══╝",
		"  ██║     ██║     ██║    ███████║██║  ███╗███████║█████╗  ██╔██╗ ██║   ██║   ",
		"  ██║     ██║     ██║    ██╔══██║██║   ██║██╔══██║██╔══╝  ██║╚██╗██║   ██║   ",
		"  ███████╗███████╗██║    ██║  ██║╚██████╔╝██║  ██║███████╗██║ ╚████║   ██║   ",
		"  ╚══════╝╚══════╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝   ",
		"                                                                            ",
		"                             (by oyminirole)                                ",
	}
	colors := []*color.Color{
		color.New(color.FgHiCyan),
		color.New(color.FgCyan),
		color.New(color.FgHiBlue),
		color.New(color.FgBlue),
		color.New(color.FgHiMagenta),
		color.New(color.FgMagenta),
	}

	for _, line := range asciiHeader {
		lineLength := len(line)
		step := float64(len(colors)) / float64(lineLength)

		for i, char := range line {
			colorIndex := int(float64(i) * step)
			if colorIndex >= len(colors) {
				colorIndex = len(colors) - 1
			}
			colors[colorIndex].Printf("%c", char)
		}
		fmt.Println()
	}

	color.New(color.FgHiCyan).Add(color.Bold).Println("\nInformation:")
	fmt.Println(" • Use a \"!\" before your query to get answers to questions, or just enter your query. Example: !How to create a folder?")
	fmt.Println(" • To agent mode, just enter your query without \"!\". Agent will execute commands and provide reports with results.")
	fmt.Println(" • To exit the program, press Ctrl+C or close the terminal window.")
	fmt.Println(strings.Repeat("─", 70))
}

func PrintErrorBox(errorOutput, aiAnalysis string) {
	width := 80

	fmt.Printf(" %s\n", ErrorColor("╭─[ Error ]"+strings.Repeat("─", width-16)))
	fmt.Printf(" %s %s\n", BorderColor("│"), LabelColor("Error log:"))
	for _, line := range strings.Split(strings.TrimSpace(errorOutput), "\n") {
		fmt.Printf(" %s   %s\n", BorderColor("│"), ErrorColor(line))
	}
	if aiAnalysis != "" {
		fmt.Printf(" %s %s\n", BorderColor("│"), BorderColor(strings.Repeat("·", width-4)))
		fmt.Printf(" %s %s\n", BorderColor("│"), AiColor("AI summary:"))
		for _, line := range strings.Split(strings.TrimSpace(aiAnalysis), "\n") {
			fmt.Printf(" %s   %s\n", BorderColor("│"), AiColor(line))
		}
	}
	fmt.Printf(" %s\n", ErrorColor("╰"+strings.Repeat("─", width-2)))
}

func PrintResultBox(commandOutput, aiSummary string) {
	width := 80

	fmt.Printf(" %s\n", SuccessColor("╭─[ Result ]"+strings.Repeat("─", width-16)))
	fmt.Printf(" %s %s\n", BorderColor("│"), LabelColor("Command output:"))
	for _, line := range strings.Split(strings.TrimSpace(commandOutput), "\n") {
		fmt.Printf(" %s   %s\n", BorderColor("│"), line)
	}

	fmt.Printf(" %s %s\n", BorderColor("│"), BorderColor(strings.Repeat("·", width-4)))
	fmt.Printf(" %s %s\n", BorderColor("│"), AiColor("AI summary:"))
	for _, line := range strings.Split(strings.TrimSpace(aiSummary), "\n") {
		fmt.Printf(" %s   %s\n", BorderColor("│"), AiColor(line))
	}

	fmt.Printf(" %s\n", SuccessColor("╰"+strings.Repeat("─", width-2)))
}

func SimpleResultBox(commandOutput string) {
	width := 80

	fmt.Printf(" %s\n", SuccessColor("╭─[ Result ]"+strings.Repeat("─", width-16)))
	fmt.Printf(" %s %s\n", BorderColor("│"), LabelColor("Command output:"))
	for _, line := range strings.Split(strings.TrimSpace(commandOutput), "\n") {
		fmt.Printf(" %s   %s\n", BorderColor("│"), line)
	}
	fmt.Printf(" %s\n", SuccessColor("╰"+strings.Repeat("─", width-2)))
}

func StartSpiner(text string) chan bool {
	stop := make(chan bool)
	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇"}
		i := 0
		for {
			select {
			case <-stop:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Print("\r\033[K")
				fmt.Print(text + " ")
				color.New(color.FgHiCyan).Add(color.Bold).Print(frames[i%len(frames)])
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()
	return stop
}
