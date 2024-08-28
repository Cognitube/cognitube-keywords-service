package keywords_service

import (
	"cognitube.com/keywords-service/dal/azure"
	"testing"
)

func TestWebhook(t *testing.T) {
	azure.RegisterCallback("https://moccasin-known-doe.ngrok-free.app/api/v1/callback")
}
