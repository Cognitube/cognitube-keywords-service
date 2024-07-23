package azure

import (
	"cognitube.com/keywords-service/env"
	"context"
	"github.com/carlmjohnson/requests"
	"log"
)

func RegisterCallback(url string) {
	payload := struct {
		DisplayName string `json:"displayName"`
		Properties  struct {
			Secret string `json:"secret"`
		}
		WebUrl string `json:"webUrl"`
		Events struct {
			TranscriptionCompletion bool `json:"transcriptionCompletion"`
		}
		Description string `json:"description"`
	}{
		DisplayName: "TranscriptionCompletionWebHook",
		Properties: struct {
			Secret string `json:"secret"`
		}{Secret: env.GetInstance().AzureKey},
		WebUrl: url,
		Events: struct {
			TranscriptionCompletion bool `json:"transcriptionCompletion"`
		}{TranscriptionCompletion: true},
		Description: "Automatically registered by AI Service",
	}

	requests.URL("https://eastus.api.cognitive.microsoft.com/speechtotext/v3.1/webhooks").
		Method("POST").
		ContentType("application/json").
		Header("Ocp-Apim-Subscription-Key", env.GetInstance().AzureKey).
		BodyJSON(payload).
		Fetch(context.Background())

	log.Println("Callback Registered: ", url)
}
