package test

import (
	"cognitube.com/keywords-service/azure"
	"testing"
)

func TestWebhook(t *testing.T) {
	azure.RegisterCallback("https://moccasin-known-doe.ngrok-free.app/api/v1/callback")
}
