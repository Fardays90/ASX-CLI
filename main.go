package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"time"

	"regexp"

	"os/signal"
	"syscall"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var apiUrl = "https://api.neunelabs.com/v1/endpoint"

type Request struct {
	Query   string      `json:"query"`
	History []ChatEntry `json:"history"`
}

type Response struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}
type ChatEntry struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
}

var previousHistory []ChatEntry = []ChatEntry{{User: "", Assistant: ""}}

func formatResponse(response string) string {
	titleColorCode := "\033[1;92m"
	resetCode := "\033[0m"

	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	formatted := boldRegex.ReplaceAllString(response, titleColorCode+"$1"+resetCode)
	lines := strings.Split(formatted, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "*") {
			lines[i] = strings.Replace(line, "*", "  \033[92m•\033[0m", 1)
		}
	}
	return strings.Join(lines, "\n")
}

func sendQuery(query string) (Response, error) {
	requestPayload := Request{
		Query:   query,
		History: previousHistory,
	}
	data, err := json.Marshal(requestPayload)
	if err != nil {
		return Response{Query: query, Response: err.Error()}, err
	}
	newReq, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request")
		return Response{Query: query, Response: err.Error()}, err
	}
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Authorization", "asx_UwEtbqIibFIDee1DTS4TjXEy0jfKG4EBt3xOKWoFxeuoLZljgk9iJbUsXehJ") //use it if you want lmao
	client := &http.Client{}
	response, err := client.Do(newReq)
	if err != nil {
		fmt.Println("Error sending the req")
		return Response{Query: query, Response: err.Error()}, err
	}
	defer response.Body.Close()
	var responseObj Response
	err = json.NewDecoder(response.Body).Decode(&responseObj)
	if err != nil {
		return Response{Query: query, Response: err.Error()}, err
	}
	chatHistory := ChatEntry{User: query, Assistant: responseObj.Response}
	previousHistory = append(previousHistory, chatHistory)

	return responseObj, nil
}

func main() {
	fmt.Println()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	welcomeColor := color.New(color.FgHiGreen).Add(color.Bold)
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	assistantColor := color.New(color.FgHiWhite).Add(color.Bold)
	assistantAnsColor := color.New(color.FgHiGreen)
	arrowColor := color.New(color.FgHiCyan).Add(color.Bold)
	messageColor := color.New(color.FgYellow).Add(color.Bold)
	newColor := color.New(color.FgHiWhite).Add(color.Bold)
	welcomeColor.Println("Welcome to ASX answer engine! Ask a question to get answers. Do CTRL+C to quit.")
	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		<-sigChan
		fmt.Print("\r\033[K")
		color.NoColor = false
		fmt.Println()
		assistantColor.Println("Thanks for using asx, check out Neune Labs")
		fmt.Println()
		os.Exit(0)
	}()
	var iteration = 0
	for {
		fmt.Println()
		arrowColor.Print("-> ")
		scanner.Scan()
		s.Start()
		response, err := sendQuery(scanner.Text())
		if err != nil {
			fmt.Println(err)
			break
		}
		s.Stop()
		fmt.Print("\r\033[K")
		color.NoColor = false
		fmt.Println()
		assistantColor.Print("ASX: ")
		assistantAnsColor.Println(formatResponse(response.Response))

		if iteration == 3 {
			fmt.Println()
			messageColor.Print("You've reached the limit, please make a ")
			newColor.Print("new")
			messageColor.Println(" chat!")
			fmt.Println()
			break
		}
		iteration++
	}
}
