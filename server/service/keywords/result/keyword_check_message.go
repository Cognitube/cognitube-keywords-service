package result

type KeywordCheckMessage struct {
	ID        string `json:"id"`
	VideoID   string `json:"videoId"`
	RetryTime int64  `json:"retryTime"`
}
