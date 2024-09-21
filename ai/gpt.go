package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	Messages    []Message `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

var content string = `
you are ai assistant
`

func Gpt(userMessage string) (string, error) {
	apiKey := os.Getenv("gpt_token") // Fetch API key from environment variables
	url := "https://api.openai.com/v1/chat/completions"

	// Append extra data if it exists

	// Prepare the request body (messages array)
	requestBody := ChatCompletionRequest{
		Model:       "gpt-4o-mini", // Use GPT-4 or another available model
		Temperature: 0.4,
		Messages: []Message{
			{Role: "system", Content: content},
			{Role: "user", Content: userMessage},
		},
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %v", err)
	}

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	var completionResponse ChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Extract the AI's response message
	if len(completionResponse.Choices) > 0 {
		aiResponse := completionResponse.Choices[0].Message.Content
		return aiResponse, nil
	}

	return "", fmt.Errorf("no AI response received")
}
