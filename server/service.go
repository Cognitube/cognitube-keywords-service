package server

import (
	"mime/multipart"
)

type ICognitubeKeywordsService interface {
	GetKeyDescFromHttpAudioFile(file multipart.File) (string, error)
}
