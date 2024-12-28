package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestData struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice
}

const PROMPT_TEMPLATE = `## Instructions##
You are a professional diary scribe. I want you to summarize my daily digital activities into a readable and concise bullet point format grouped by activity.
To assist you, I have the log lines of everything I have done during the day. The logs are in the format of "[time in hour:minutes:seconds] <description of what I was doing that time>".
For every activity (e.g. YouTube, Coding, etc.), itemize all the different sub-activities that I did with a short description.
Then, give me a paragraph summarizing my day using poetic language (be creative).

## Logs ##
%s

## Summary ##`

func readLogs(logFile string) string {
	content, err := os.ReadFile(logFile)
	if err != nil {
		log.Fatalf("Could not read log file content: %v", err)
	}
	return string(content)
}

// TODO: breakdown the logs into batches
// get { batch -> prompt -> response } -> responses[]
// formulate an aggregate prompt to summarize responses[]
func getPrompt(logFile string) string {
	content := readLogs(logFile)
	content = strings.Join(strings.Split(content, "\n")[:100], "\n")
	return fmt.Sprintf(PROMPT_TEMPLATE, content)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./scribe <path_to_logs>")
		os.Exit(1)
	}
	logFile := os.Args[1]
	day := strings.Split(strings.Split(logFile, "logs_")[1], ".txt")[0]
	prompt := getPrompt(logFile)
	requestData := RequestData{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Model:       "llama-3.3-70b-versatile",
		Temperature: 0.7,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}
	groqApiKey := os.Getenv("GROQ_API_KEY")
	url := "https://api.groq.com/openai/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+groqApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Print(resp)
	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error unmarshalling body: %v", err)
	}

	filePath := fmt.Sprintf("diary_%s.txt", day)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	content := response.Choices[0].Message.Content
	_, err = writer.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	writer.Flush()

	fmt.Printf("Wrote diary to %s\n", filePath)
}
