package test

import (
	"cognitube.com/keywords-service/keydesc"
	"os"
	"testing"
)

func TestGPTClient_Describe(t *testing.T) {
	os.Setenv("APPLICATION_OPENAI_KEY", "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")
	os.Setenv("APPLICATION_GPT_URL", "https://api.openai.com/v1/chat/completions")
	g := keydesc.NewGPTClient()
	text := "This video is brought to you by Captivating History. Joan of Arc, or Jean Darc as she is known in France, remains an inspirational and controversial figure 600 years after her birth. St. or Heretic, touched by the grace of God or a tragically delusional young woman. Even while alive, these points of view were hotly debated. The France into which Jean was born in 1412 was divided. There was a relative peace in the long time war with England. However, internal unrest was caused by the feud between 2 royal factions who sat on the Regency council, the Armanyacs and the Burgundians. The council was headed by Queen Isabel, the wife of King Charles the 6th. He was also known as Charles the Man due to periods of psychosis which had begun in the early 1390s."
	result, err := g.Describe(text)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
