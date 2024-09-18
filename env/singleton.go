package env

import (
	"os"
	"sync"
)

type Variables struct {
	OpenAIKey                string
	WhisperUrl               string
	GptUrl                   string
	Debug                    bool
	KeywordExtractionPrompt  string
	AzureKey                 string
	AzureUrl                 string
	KafkaHost                string
	KafkaPort                string
	KafkaTopic               string
	KafkaKeywordCheckTopic   string
	BlobConnectString        string
	TranscriptContainerName  string
	KeywordsContainerName    string
	KeywordDescMaxRetry      int
	ApplicationCallbackUrl   string
	EventHubNamespace        string
	EventHubConnectionString string
	Username                 string
	RedisConnectionString    string
	SubtitleContainerName    string
	KafkaBootstrapServers    string
}

var instance *Variables
var once sync.Once

func loadValues() {
	instance = &Variables{
		OpenAIKey:                os.Getenv("APPLICATION_OPENAI_KEY"),
		WhisperUrl:               os.Getenv("APPLICATION_WHISPER_URL"),
		GptUrl:                   os.Getenv("APPLICATION_GPT_URL"),
		Debug:                    os.Getenv("APPLICATION_DEBUG") == "true",
		KeywordExtractionPrompt:  "You are a keyword extraction model. Please extract 20 - 30 keywords from the following user context, where each keyword should be low frequency nouns that is important to the user context. The user input format is a list of json file, where the start and the text attributes in segment are very useful in this problem. your output should looks like this: [{keyword:\"keyword1\", description:\"description of keyword1 in the given context\",startTime:0.0,endTime:3.83}, {keyword:\"keyword2\", explain:\"description of keyword2 in the given context\",startTime:3.83,endTime:5.17}, ...]. The description of the keywords should be specific, first explain the actual meaning of this keyword, and then related to the context explain what this keyword means in the video, and no less than 150 words. The startTime and endTime of the return value should round to two decimal places.",
		AzureKey:                 os.Getenv("APPLICATION_AZURE_KEY"),
		AzureUrl:                 os.Getenv("APPLICATION_AZURE_URL"),
		KafkaHost:                os.Getenv("APPLICATION_KAFKA_HOST"),
		KafkaPort:                os.Getenv("APPLICATION_KAFKA_PORT"),
		KafkaTopic:               os.Getenv("APPLICATION_KAFKA_TOPIC"),
		KafkaKeywordCheckTopic:   os.Getenv("APPLICATION_KAFKA_KEYWORD_CHECK_TOPIC"),
		BlobConnectString:        os.Getenv("AZURE_BLOB_CONNECTION_STRING"),
		TranscriptContainerName:  os.Getenv("TRANSCRIPT_CONTAINER_NAME"),
		KeywordsContainerName:    os.Getenv("KEYWORDS_CONTAINER_NAME"),
		KeywordDescMaxRetry:      3,
		ApplicationCallbackUrl:   os.Getenv("APPLICATION_CALLBACK_URL"),
		EventHubNamespace:        os.Getenv("KAFKA_EVENTHUB_NAMESPACE"),
		EventHubConnectionString: os.Getenv("AZURE_EVENTHUB_CONNECTIONSTRING"),
		Username:                 os.Getenv("KAFKA_EVENTHUB_USERNAME"),
		RedisConnectionString:    os.Getenv("REDIS_CONNECTION_STRING"),
		SubtitleContainerName:    os.Getenv("SUBTITLE_CONTAINER_NAME"),
		KafkaBootstrapServers:    os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}
}

func GetInstance() *Variables {
	if instance == nil {
		once.Do(loadValues)
	}
	return instance
}
