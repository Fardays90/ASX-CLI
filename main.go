package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var apiUrl = "https://api.neunelabs.com/v1/endpoint"

type Request struct {
	Query   string `json:"query"`
	History string `json:"history"`
}

type Response struct {
	Query    string `json:"query"`
	Response string `json:"response"`
}

func sendQuery(query string) (Response, error) {
	requestPayload := Request{
		Query:   query,
		History: "",
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
	return responseObj, nil
}

func main() {
	fmt.Println("Welcome to ASX answer engine ask a question to get answers!")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("-> ")
		scanner.Scan()
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "exit" {
			fmt.Println("bye bye")
			break
		}
		response, err := sendQuery(scanner.Text())
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Print("ASX: ")
		fmt.Println(response.Response)
	}
}
