package keywords_service

import (
	"cognitube.com/keywords-service/server/service/keywords/keydesc"
	"os"
	"testing"
)

func TestGPTClient_Describe(t *testing.T) {
	os.Setenv("APPLICATION_OPENAI_KEY", "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")
	os.Setenv("APPLICATION_GPT_URL", "https://api.openai.com/v1/chat/completions")
	g := keydesc.NewGPTClient()
	text := "{\n  \"duration\": 54.27,\n  \"text\": \"This video is brought to you by Captivating History. Joan of Arc, or Jean Darc as she is known in France, remains an inspirational and controversial figure 600 years after her birth. St. or Heretic, touched by the grace of God or a tragically delusional young woman. Even while alive, these points of view were hotly debated. The France into which Jean was born in 1412 was divided. There was a relative peace in the long time war with England. However, internal unrest was caused by the feud between 2 royal factions who sat on the Regency council, the Armanyacs and the Burgundians. The council was headed by Queen Isabel, the wife of King Charles the 6th. He was also known as Charles the Man due to periods of psychosis which had begun in the early 1390s.\",\n  \"segment\": [\n    {\n      \"start\": 0.44,\n      \"end\": 3.52,\n      \"text\": \"This video is brought to you by Captivating History.\"\n    },\n    {\n      \"start\": 4.88,\n      \"end\": 13.64,\n      \"text\": \"Joan of Arc, or Jean Darc as she is known in France, remains an inspirational and controversial figure 600 years after her birth.\"\n    },\n    {\n      \"start\": 14.32,\n      \"end\": 14.8,\n      \"text\": \"St.\"\n    },\n    {\n      \"start\": 15,\n      \"end\": 20.4,\n      \"text\": \"or Heretic, touched by the grace of God or a tragically delusional young woman.\"\n    },\n    {\n      \"start\": 21.28,\n      \"end\": 24.64,\n      \"text\": \"Even while alive, these points of view were hotly debated.\"\n    },\n    {\n      \"start\": 26.08,\n      \"end\": 29.759999999999998,\n      \"text\": \"The France into which Jean was born in 1412 was divided.\"\n    },\n    {\n      \"start\": 30.36,\n      \"end\": 33.12,\n      \"text\": \"There was a relative peace in the long time war with England.\"\n    },\n    {\n      \"start\": 33.52,\n      \"end\": 42.32000000000001,\n      \"text\": \"However, internal unrest was caused by the feud between 2 royal factions who sat on the Regency council, the Armanyacs and the Burgundians.\"\n    },\n    {\n      \"start\": 43,\n      \"end\": 47.24,\n      \"text\": \"The council was headed by Queen Isabel, the wife of King Charles the 6th.\"\n    },\n    {\n      \"start\": 47.92,\n      \"end\": 53.88,\n      \"text\": \"He was also known as Charles the Man due to periods of psychosis which had begun in the early 1390s.\"\n    }\n  ]\n}"
	result, err := g.Describe(text)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}
