package test

import (
	"cognitube.com/keywords-service/keydesc"
	"os"
	"testing"
)

//	func TestGeminiClient_Describe(t *testing.T) {
//		os.Setenv("APPLICATION_DEBUG", "true")
//		os.Setenv("APPLICATION_OPENAI_KEY", "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")
//		os.Setenv("APPLICATION_KEYWORD_SERVICE_URL", "https://cognitube-keywords-service.azurewebsites.net/api")
//
//		client := keydesc.NewGeminiClient()
//		keywords, err := client.Describe("This video is brought to you by Captivating History. Joan of Arc, or Jeanne d'Arc as she is known in France, remains an inspirational and controversial figure 600 years after her birth. Saint or heretic, touched by the grace of God, or a tragically delusional young woman. Even while alive, these points of view were hotly debated. The France into which Joan was born in 1412 was divided. There was a relative peace in the long-time war with England. However, internal unrest was caused by the feud between two royal factions who sat on the Regency Council, the Armagnacs and the Burgundians. The council was headed by Queen Isabeau, the wife of King Charles VI. He was also known as Charles the Mad, due to periods of psychosis which had begun in the early 1390s.")
//		if err != nil {
//			t.Errorf("Error getting keywords: %v", err)
//		}
//		t.Logf("Keywords: %v", keywords)
//	}
func TestGPTClient_Describe(t *testing.T) {
	os.Setenv("APPLICATION_DEBUG", "true")
	os.Setenv("APPLICATION_OPENAI_KEY", "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")
	os.Setenv("APPLICATION_GPT_URL", "https://api.openai.com/v1/chat/completions")
	os.Setenv("APPLICATION_KEYWORD_SERVICE_URL", "https://cognitube-keywords-service.azurewebsites.net/api")

	client := keydesc.NewGPTClient()
	keywords, err := client.Describe("This video is brought to you by Captivating History. Joan of Arc, or Jeanne d'Arc as she is known in France, remains an inspirational and controversial figure 600 years after her birth. Saint or heretic, touched by the grace of God, or a tragically delusional young woman. Even while alive, these points of view were hotly debated. The France into which Joan was born in 1412 was divided. There was a relative peace in the long-time war with England. However, internal unrest was caused by the feud between two royal factions who sat on the Regency Council, the Armagnacs and the Burgundians. The council was headed by Queen Isabeau, the wife of King Charles VI. He was also known as Charles the Mad, due to periods of psychosis which had begun in the early 1390s.")
	if err != nil {
		t.Errorf("Error getting keywords: %v", err)
	}
	t.Logf("Keywords: %v", keywords)
}
