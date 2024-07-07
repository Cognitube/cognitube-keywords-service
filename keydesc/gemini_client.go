package keydesc

import (
	"cognitube.com/keywords-service/env"
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"regexp"
)

// Deprecated: GeminiClient is deprecated
type GeminiClient struct {
	apiKey string
}

func NewGeminiClient() *GeminiClient {
	return &GeminiClient{
		apiKey: env.GetInstance().OpenAIKey,
	}
}

func (c *GeminiClient) Describe(transcript string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(c.apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")

	model.SetTemperature(1)
	model.SetTopP(0.95)
	model.SetTopK(0)
	model.SetMaxOutputTokens(65536)

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text("You are a keyword extraction model. \n Please extract 20 - 30 keywords from the following user context, where each keyword should be low frequency nouns that is important to the user context. The user input format is a list of json file, where the start and the text attributes in segment are very useful in this problem. your output should looks like this: \n[{keyword:\\\"keyword1\\\", description:\\\"description of keyword1 in the given context\\\",startTime:0.0,endTime:3.83}, {keyword:\\\"keyword2\\\", explain:\\\"description of keyword2 in the given context\\\",startTime:3.83,endTime:5.17}, ...]. The description of the keywords should be specific, first explain the actual meaning of this keyword, and then related to the context explain what this keyword means in the video, and no less than 150 words. The startTime and endTime of the return value should round to two decimal places. Display only json. Do not format it in markdown.\n\n")},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(transcript))
	if err != nil {
		log.Fatal(err)
	}

	res := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				res += string(part.(genai.Text))
			}
		}
	}

	re := regexp.MustCompile(`(?s)\[.*?\]`)
	matches := re.FindStringSubmatch(res)
	if len(matches) == 0 {
		return "", fmt.Errorf("no match found")
	}

	return matches[0], nil
}
