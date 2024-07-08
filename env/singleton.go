package env

import (
	"os"
	"sync"
)

type Variables struct {
	OpenAIKey               string
	WhisperUrl              string
	GptUrl                  string
	Debug                   bool
	KeywordExtractionPrompt string
}

var instance *Variables
var once sync.Once

func loadValues() {
	instance = &Variables{
		OpenAIKey:               os.Getenv("APPLICATION_OPENAI_KEY"),
		WhisperUrl:              os.Getenv("APPLICATION_WHISPER_URL"),
		GptUrl:                  os.Getenv("APPLICATION_GPT_URL"),
		Debug:                   os.Getenv("APPLICATION_DEBUG") == "true",
		KeywordExtractionPrompt: "You are a keyword extraction model. Please extract 20 - 30 keywords from the following user context, where each keyword should be low frequency nouns that is important to the user context. The user input format is a list of json file, where the start and the text attributes in segment are very useful in this problem. your output should looks like this: [{keyword:\"keyword1\", description:\"description of keyword1 in the given context\",startTime:0.0,endTime:3.83}, {keyword:\"keyword2\", explain:\"description of keyword2 in the given context\",startTime:3.83,endTime:5.17}, ...]. The description of the keywords should be specific, first explain the actual meaning of this keyword, and then related to the context explain what this keyword means in the video, and no less than 150 words. The startTime and endTime of the return value should round to two decimal places.",
	}
}

func GetInstance() *Variables {
	if instance == nil {
		once.Do(loadValues)
	}
	return instance
}
