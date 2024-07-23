package keydesc

import (
	"bytes"
	"cognitube.com/keywords-service/env"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GPTClient struct {
	apiKey string
	gptUrl string
}

func NewGPTClient() *GPTClient {
	return &GPTClient{
		apiKey: env.GetInstance().OpenAIKey,
		gptUrl: env.GetInstance().GptUrl,
	}
}

// TODO: Add retry logic
func (g *GPTClient) Describe(transcript string) (string, error) {
	systemMessage := map[string]string{
		"role":    "system",
		"content": env.GetInstance().KeywordExtractionPrompt,
	}
	userMessage := map[string]string{
		"role":    "user",
		"content": transcript,
	}
	messages := []map[string]string{systemMessage, userMessage}

	requestObj := map[string]interface{}{
		"model":    "gpt-4-turbo-preview",
		"messages": messages,
	}
	jsonData, err := json.Marshal(requestObj)
	if err != nil {
		return "", fmt.Errorf("error marshalling request data: %w", err)
	}

	req, err := http.NewRequest("POST", g.gptUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	var responseObj map[string]interface{}
	if err := json.Unmarshal(responseData, &responseObj); err != nil {
		return "", fmt.Errorf("error parsing response JSON: %w", err)
	}

	choices, ok := responseObj["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("choices field not found or empty")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid format for first choice")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message field not found or invalid in first choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content field not found or invalid in message")
	}

	return content, nil
}
