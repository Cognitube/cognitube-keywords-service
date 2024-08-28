package redis

const (
	RedisKeyTranscriptPrefix = "ai-service:transcription-key:"
)

func GetTranscriptRedisKey(tid string) string {
	return RedisKeyTranscriptPrefix + tid
}
