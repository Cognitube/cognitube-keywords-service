package keywords

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	azure2 "cognitube.com/keywords-service/dal/azure"
	"cognitube.com/keywords-service/dal/mq"
	myredis "cognitube.com/keywords-service/dal/redis"
	keydesc2 "cognitube.com/keywords-service/server/service/keywords/keydesc"
	"cognitube.com/keywords-service/server/service/keywords/result"
	transcription2 "cognitube.com/keywords-service/server/service/keywords/transcription"
	"github.com/go-redis/redis/v8"

	"cognitube.com/keywords-service/env"
	"github.com/tidwall/gjson"
)

// ICognitubeKeywordsService receives the direct input from the HTTP request and returns the output to the HTTP response.
// Considering HTTP response time. Do not perform time-consuming operations outside the scope of the interface
// Move the time-consuming operations to the implementation of the interface, and if the return value is not necessary for the HTTP response, use goroutines
type ICognitubeKeywordsService interface {
	GetKeyDescFromHttpAudioFile(file multipart.File) (string, error)
	CreateAsyncTranscription(fileUrl string, videoID string) (string, error)
	OnTranscriptionCallback(payload []byte)
	ProcessTranscriptionResult(id string)
}

type CognitubeKeywordsService struct {
	descriptor       keydesc2.KeywordsDescriptor
	transcriber      transcription2.Transcriber
	asyncTranscriber transcription2.AsyncTranscriber
	resultPublisher  mq.TranscriptionPublisher
	blobClient       *azure2.BlobClient
}

func NewCognitubeKeywordsService() ICognitubeKeywordsService {
	return &CognitubeKeywordsService{
		descriptor:       keydesc2.NewKeywordsDescriptor("gpt"),
		transcriber:      transcription2.NewTranscriber("whisper"),
		asyncTranscriber: transcription2.NewAsyncTranscriberClient("azure"),
		resultPublisher:  mq.NewKafkaTranscriptionPublisher(),
		blobClient:       azure2.NewBlobClient(),
	}
}

func convertMultipartFileToOsFile(file multipart.File) (*os.File, error) {
	// Convert multipart.File to os.File
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "uploaded-*.mp3")
	if err != nil {
		return nil, err
	}

	// Copy the content of the multipart file to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close() // Attempt to close the file in case of error
		return nil, err
	}

	// It's important to seek to the beginning of the file before returning it
	if _, err := tempFile.Seek(0, 0); err != nil {
		tempFile.Close() // Attempt to close the file in case of error
		return nil, err
	}

	return tempFile, nil
}

func (c *CognitubeKeywordsService) GetKeyDescFromHttpAudioFile(multiFile multipart.File) (string, error) {
	// Convert multipart.File to os.File
	defer multiFile.Close()
	file, err := convertMultipartFileToOsFile(multiFile)
	if err != nil {
		return "", err
	}

	// Transcribe the audio file to text
	text, err := c.transcriber.Transcript(file)
	if err != nil {
		return "", err
	}

	// Describe the text to keywords
	keywords, err := c.descriptor.Describe(text)
	if err != nil {
		return "", err
	}
	return keywords, nil
}

func (c *CognitubeKeywordsService) CreateAsyncTranscription(fileUrl string, videoID string) (string, error) {
	id, err := c.asyncTranscriber.CreateTranscription(fileUrl, videoID)
	if err != nil {
		return "", err
	}

	PutTranscriptIDToVideoID(id, videoID)
	c.addTranscriptionResultCheckMessage(id, videoID)
	return id, nil
}

func (c *CognitubeKeywordsService) addTranscriptionResultCheckMessage(transcriptionID, videoID string) {
	newKafkaMessage := result.KeywordCheckMessage{
		ID:        transcriptionID,
		VideoID:   videoID,
		RetryTime: time.Now().Add(2 * time.Minute).Unix(),
	}

	msg, _ := json.Marshal(newKafkaMessage)
	// Add the message to the message queue
	err := c.resultPublisher.Publish(env.GetInstance().KafkaKeywordCheckTopic, msg)
	if err != nil {
		log.Println("Failed to add transcription result check message to Kafka")
	}
}

func (c *CognitubeKeywordsService) ProcessTranscriptionResult(id string) {
	go c.processTranscriptionResult(id, true, "")
}

func (c *CognitubeKeywordsService) processTranscriptionResult(id string, isSuccess bool, errString string) {
	success := isSuccess
	errStr := errString

	vid := PopVideoIDFromTranscriptID(id) // vid might be ""
	log.Println("Find Transcription Job ID: ", id, " for Video ID: ", vid)

	log.Println("Getting transcription result for job ID: ", id)
	transcript, err := c.asyncTranscriber.OnTranscriptionCallback(id) // transcript might be ""
	if err != nil {
		success, errStr = false, err.Error()
	}

	standard := transcript.ToStandard()
	subtitle := standard.ToJson()

	retry := -1
	desc := ""
	for keydesc2.ValidateKeywordsResult(desc) != true {
		log.Println("Generating keywords description for job ID: ", id)
		log.Println("Trying " + strconv.Itoa(retry+2) + " times")
		desc, err = c.descriptor.Describe(subtitle)
		if err != nil {
			success = false
			errStr = err.Error()
			break
		}
		retry++
		if retry >= env.GetInstance().KeywordDescMaxRetry {
			success = false
			errStr = "Failed to generate keywords description (reached max retry)"
			log.Printf("Failure reason for video %s: %s", vid, desc)
			break
		}
	}

	transUrl := ""
	if transcript.GetText() != "" && success {
		log.Println("Publishing transcription result for job ID: ", id)
		transUrl, err = c.blobClient.UploadTranscript(id+".txt", transcript.GetText())
	}

	keywordsUrl := ""
	if desc != "" && success {
		log.Println("Publishing keywords result for job ID: ", id)
		keywordsUrl, _ = c.blobClient.UploadKeywords(id+".json", desc)
	}

	subtitleUrl := ""
	if subtitle != "" && success {
		log.Println("Publishing subtitle for job ID: ", id)
		subtitleUrl, err = c.blobClient.UploadSubtitle(id+".json", subtitle)
	}

	log.Println("Publishing final result for job ID: ", id)
	c.resultPublisher.PublishTranscriptionResult(&result.Result{
		Success:       success,
		VideoID:       vid,
		KeywordsURL:   keywordsUrl,
		TranscriptURL: transUrl,
		Error:         errStr,
		SubtitleURL:   subtitleUrl,
	})
}

func (c *CognitubeKeywordsService) OnTranscriptionCallback(payload []byte) {
	success := true
	errStr := ""

	url := gjson.GetBytes(payload, "self").String()
	if url == "" {
		success, errStr = false, "self URL not found in payload"
	}

	id := azure2.GetJobIDFromSelfURL(url)
	if id == "" {
		success, errStr = false, "invalid self URL"
	}

	c.processTranscriptionResult(id, success, errStr)
}

func GetRedisOptions() *redis.Options {
	if env.GetInstance().Debug {
		option := &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}
		return option
	}
	option, err := redis.ParseURL(env.GetInstance().RedisConnectionString)
	if err != nil {
		log.Println("Failed to parse redis connection string")
	}
	return option
}

func PutTranscriptIDToVideoID(tid string, vid string) {
	log.Println("Putting transcript ID to video ID mapping, ", tid, " -> ", vid)
	option := GetRedisOptions()
	client := redis.NewClient(option)
	defer client.Close()
	err := client.Set(context.Background(), myredis.GetTranscriptRedisKey(tid), vid, 5*time.Hour).Err()
	if err != nil {
		log.Println("Failed to set transcript ID to video ID mapping")
	}
}

func PopVideoIDFromTranscriptID(tid string) string {
	log.Println("Popping video ID from transcript ID: ", tid)
	option := GetRedisOptions()
	client := redis.NewClient(option)
	defer client.Close()

	vid, err := client.Get(context.Background(), myredis.GetTranscriptRedisKey(tid)).Result()
	if err != nil {
		log.Println("Failed to get video ID for transcript ID: ", tid)
	}

	// Delete the key-value pair after retrieving it
	err = client.Del(context.Background(), myredis.GetTranscriptRedisKey(tid)).Err()
	if err != nil {
		log.Println("Failed to delete video ID for transcript ID: ", tid)
	}
	return vid
}
