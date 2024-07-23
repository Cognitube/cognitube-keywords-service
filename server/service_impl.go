package server

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"cognitube.com/keywords-service/azure"
	"cognitube.com/keywords-service/env"
	"cognitube.com/keywords-service/keydesc"
	"cognitube.com/keywords-service/publish"
	"cognitube.com/keywords-service/result"
	"cognitube.com/keywords-service/transcription"
	"github.com/tidwall/gjson"
)

type CongnitubeKeywordsService struct {
	descriptor       keydesc.KeywordsDescriptor
	transcriber      transcription.Transcriber
	asyncTranscriber transcription.AsyncTranscriber
	resultPublisher  publish.TranscriptionPublisher
	blobClient       *azure.BlobClient
}

func NewCognitubeKeywordsService() ICognitubeKeywordsService {
	return &CongnitubeKeywordsService{
		descriptor:       keydesc.NewKeywordsDescriptor("gpt"),
		transcriber:      transcription.NewTranscriber("whisper"),
		asyncTranscriber: transcription.NewAsyncTranscriberClient("azure"),
		resultPublisher:  publish.NewKafkaTranscriptionPublisher(),
		blobClient:       azure.NewBlobClient(),
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

func (c *CongnitubeKeywordsService) GetKeyDescFromHttpAudioFile(multiFile multipart.File) (string, error) {
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

func (c *CongnitubeKeywordsService) CreateAsyncTranscription(fileUrl string, videoID string) (string, error) {
	id, err := c.asyncTranscriber.CreateTranscription(fileUrl, videoID)
	if err != nil {
		return "", err
	}
	PutTranscriptIDToVideoID(videoID, id)
	return id, nil
}

var videoIDToTranscriptID = make(map[string]string)

func PutTranscriptIDToVideoID(tid string, vio string) {
	videoIDToTranscriptID[vio] = tid
}

func (c *CongnitubeKeywordsService) OnTranscriptionCallback(payload []byte) {
	success := true
	errStr := ""

	url := gjson.GetBytes(payload, "self").String()
	if url == "" {
		success, errStr = false, "self URL not found in payload"
	}

	id := azure.GetJobIDFromSelfURL(url)
	if id == "" {
		success, errStr = false, "invalid self URL"
	}

	vid := PopVideoIDFromTranscriptID(id) // vid might be ""
	log.Println("Find Transcription Job ID: ", id, " for Video ID: ", vid)

	log.Println("Getting transcription result for job ID: ", id)
	transcript, err := c.asyncTranscriber.OnTranscriptionCallback(id) // transcript might be ""
	if err != nil {
		success, errStr = false, err.Error()
	}

	retry := -1
	desc := ""
	for keydesc.ValidateKeywordsResult(desc) != true {
		log.Println("Generating keywords description for job ID: ", id)
		log.Println("Trying " + strconv.Itoa(retry+2) + " times")
		desc, err = c.descriptor.Describe(transcript)
		if err != nil {
			success = false
			errStr = err.Error()
			break
		}
		retry++
		if retry >= env.GetInstance().KeywordDescMaxRetry {
			success = false
			errStr = "Failed to generate keywords description (reached max retry)"
			log.Fatalf("Failure reason for video %s: %s", vid, desc)
			break
		}
	}

	transUrl := ""
	if transcript != "" && success {
		log.Println("Publishing transcription result for job ID: ", id)
		transUrl, _ = c.blobClient.UploadTranscript(id+".txt", transcript)
	}

	keywordsUrl := ""
	if desc != "" && success {
		log.Println("Publishing keywords result for job ID: ", id)
		keywordsUrl, _ = c.blobClient.UploadKeywords(id+".json", desc)
	}

	log.Println("Publishing final result for job ID: ", id)
	c.resultPublisher.PublishTranscriptionResult(&result.Result{
		Success:       success,
		VideoID:       vid,
		KeywordsURL:   keywordsUrl,
		TranscriptURL: transUrl,
		Error:         errStr,
	})
}

func PopVideoIDFromTranscriptID(tid string) string {
	vid := videoIDToTranscriptID[tid]
	if vid == "" {
		log.Println("No video ID found for transcript ID: ", tid)
		return ""
	}
	delete(videoIDToTranscriptID, tid)
	return vid
}
