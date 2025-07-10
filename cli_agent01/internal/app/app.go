package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Role1776/agent/internal/config"
	"github.com/Role1776/agent/internal/shell"
	"github.com/Role1776/agent/internal/ui"

	"net/http"

	"github.com/fatih/color"
)

const (
	systemPromptBase            = "System Context:\n"
	simpleChatPromptTemplate    = "You are an AI agent. Answer the user's question based on the information available. NEVER USE SMILEYS."
	commandGenPromptTemplate    = "You are an AI agent. Your task is to generate commands based on the user's query. Answer only with the command, without extra words, explanations, and markdown formatting. Only raw command. Never use smileys. Try to generate commands that do not produce very long logs"
	errorAnalysisPromptTemplate = "Analyze the error execution of command and explain simply what went wrong. Original query: '%s'. Error: '%s'"
	summaryPromptTemplate       = "Briefly explain the result of executing the command, based on the original user query. Original query: '%s'. Command output: '%s'"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type APIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type APIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message ResponseMessage `json:"message"`
}

type ResponseMessage struct {
	Content string `json:"content"`
}

type App struct {
	client       *http.Client
	history      []Message
	autoComplete bool
	reader       *bufio.Reader
	config       *config.Config
}

func NewApp(config *config.Config, client *http.Client) *App {
	return &App{
		client:       client,
		autoComplete: false,
		reader:       bufio.NewReader(os.Stdin),
		config:       config,
		history: []Message{
			{
				Role:    "system",
				Content: systemPromptBase + "\n" + shell.GetSystemInfo(),
			},
		},
	}
}

func (a *App) promptUser() (string, error) {
	fmt.Println("╭─" + strings.Repeat("─", 66))
	color.New(color.FgHiWhite).Print("│ > Enter your query: ")
	userInput, err := a.reader.ReadString('\n')
	color.New(color.FgHiWhite).Println("╰" + strings.Repeat("─", 65))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(userInput), nil
}

func (a *App) Run() {
	ui.PrintHeader()

	for {
		userInput, err := a.promptUser()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		switch userInput {
		case "/exit":
			log.Println(ui.AiColor("exit"))
			return
		case "/clear":
			a.CleanHistory()
			log.Println(ui.SuccessColor("History cleared"))
			continue
		case "/auto-true":
			a.autoComplete = true
			log.Println(ui.SuccessColor("Auto complete enabled"))
			continue
		case "/auto-false":
			a.autoComplete = false
			log.Println(ui.SuccessColor("Auto complete disabled"))
			continue
		default:
			if isCommandPrefixed(userInput) {
				a.handleSimpleChat(userInput)
			} else {
				a.handleAgentMode(userInput, a.autoComplete)
			}
		}

	}
}

func (a *App) prepareAPImessages(systemPrompt, userInput string) []Message {
	messages := make([]Message, 0, len(a.history)+2)
	messages = append(messages, Message{Role: "system", Content: systemPrompt})
	messages = append(messages, a.history...)
	messages = append(messages, Message{Role: "user", Content: userInput})
	return messages
}

func (a *App) handleSimpleChat(userInput string) {
	spiner := ui.StartSpiner("Thinking...")

	cleanInput := strings.TrimPrefix(userInput, "!")
	messagesToGenerate := a.prepareAPImessages(simpleChatPromptTemplate, cleanInput)

	response, err := a.generateContent(messagesToGenerate)
	if err != nil {
		spiner <- true
		log.Println(ui.ErrorColor(err.Error()))
		return
	}
	spiner <- true
  
        a.history = append(a.history, Message{Role: "user", Content: userInput})
	a.history = append(a.history, Message{Role: "assistant", Content: response})
	
	ui.SimpleResultBox(response)
}

func (a *App) handleAgentMode(userInput string, autoComplete bool) {
	attemptHistory := a.prepareAPImessages(commandGenPromptTemplate, userInput)
	var lastError string
	for i := 0; i < a.config.Retries; i++ {
		spiner := ui.StartSpiner("Command generation...")
		command, err := a.generateContent(attemptHistory)
		spiner <- true

		if err != nil {
			ui.PrintErrorBox(fmt.Sprintf("Command generation error:\n%v", err), "")
			return
		}

		cleanCommand := cleanCommand(command)
		if !autoComplete {
			if !a.askForConfirmation(cleanCommand) {
				log.Println(ui.AiColor("Command cancelled by user."))
				return
			}
		}
		attemptHistory = append(attemptHistory, Message{Role: "assistant", Content: command})
		spinerExecution := ui.StartSpiner("Command execution...")
		commandOutput, err := shell.ExecuteCommand(cleanCommand)
		spinerExecution <- true

		if err == nil {
			spinnerSummary := ui.StartSpiner("Generating summary...")

			summaryPrompt := []Message{
				{
					Role:    "system",
					Content: "You are an AI agent who short explain the result of executing the command, based on the original user query. NEVER USE SMILEYS.",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf(summaryPromptTemplate, userInput, commandOutput),
				},
			}
			summary, err := a.generateContent(summaryPrompt)
			spinnerSummary <- true
			if err != nil {
				ui.PrintErrorBox(fmt.Sprintf("Summary generation error:\n%v", err), "")
				return
			}
			ui.PrintResultBox(commandOutput, summary)
			a.history = append(a.history,
				Message{Role: "user", Content: userInput},
				Message{Role: "assistant", Content: fmt.Sprintf("Command: `%s`\nSummary: %s", command, summary)},
			)
			return
		}
		lastError = commandOutput
		errorFeedback := fmt.Sprintf("This command did not work. Output was:\n%s\nTry another command.", commandOutput)
		attemptHistory = append(attemptHistory, Message{Role: "user", Content: errorFeedback})

		time.Sleep(time.Second)
	}

	spinnerErrorAnalysis := ui.StartSpiner("Generating error analysis...")
	analysis, err := a.generateContent([]Message{
		{
			Role:    "system",
			Content: "You are an AI agent who short explain the result of executing the command, based on the original user query. NEVER USE SMILEYS.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf(errorAnalysisPromptTemplate, userInput, lastError),
		},
	})
	spinnerErrorAnalysis <- true
	if err != nil {
		ui.PrintErrorBox(fmt.Sprintf("Error analysis generation error:\n%v", err), "")
		return
	}
	ui.PrintErrorBox(fmt.Sprintf("Failed to execute task after %d attempts.", a.config.Retries), analysis)
}

func (a *App) askForConfirmation(command string) bool {
	fmt.Print("Confirm command? [y/n]: ", ui.AiColor(command), "\n> ")
	confirmInput, _ := a.reader.ReadString('\n')
	confirmInput = strings.ToLower(strings.TrimSpace(confirmInput))
	fmt.Print("\033[2A\033[J")
	return confirmInput == "y" || confirmInput == "yes"
}

func (a *App) generateContent(messages []Message) (string, error) {
	reqBody := APIRequest{
		Model:    a.config.Model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error encoding JSON: %w", err)
	}

	req, err := http.NewRequest("POST", a.config.ApiUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.config.ApiToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned error (status %d): %s", resp.StatusCode, string(respBytes))
	}
	var result APIResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) > 0 && result.Choices[0].Message.Content != "" {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("API did not return content in the response")
}

func isCommandPrefixed(input string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(input)), "!")
}

func cleanCommand(cmdStr string) string {
	cmdStr = strings.TrimPrefix(cmdStr, "```powershell")
	cmdStr = strings.TrimPrefix(cmdStr, "```bash")
	cmdStr = strings.TrimPrefix(cmdStr, "```")
	cmdStr = strings.TrimSuffix(cmdStr, "```")
	cmdStr = strings.TrimPrefix(cmdStr, "Command: ")

	if strings.HasPrefix(strings.ToLower(cmdStr), "powershell -command ") {
		firstQuote := strings.Index(cmdStr, "\"")
		lastQuote := strings.LastIndex(cmdStr, "\"")
		if firstQuote != -1 && lastQuote > firstQuote {
			cmdStr = cmdStr[firstQuote+1 : lastQuote]
		}
	}
	return strings.TrimSpace(cmdStr)
}

func (a *App) CleanHistory() {
	a.history = []Message{
		{
			Role:    "system",
			Content: systemPromptBase + "\n" + shell.GetSystemInfo(),
		},
	}
}
