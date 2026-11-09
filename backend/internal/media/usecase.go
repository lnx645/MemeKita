package media

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
)

const maxImage int64 = 10*1024*1024
const maxGif int64 = 15*1024*1024
const maxVideo int64 = 50*1024*1024

func Validate(mediaType string, size int64, contentType string) error {
	limits := map[string]int64{"image":maxImage,"gif":maxGif,"video":maxVideo}
	allowed := map[string][]string{
		"image": {"image/jpeg","image/png","image/webp"},
		"gif": {"image/gif"},
		"video": {"video/mp4","video/webm","video/quicktime"},
	}
	max, ok := limits[mediaType]
	if !ok { return fmt.Errorf("tipe tidak dikenal") }
	if size > max { return fmt.Errorf("File terlalu besar. Maksimal %dMB", max/1024/1024) }
	for _, a := range allowed[mediaType] {
		if contentType == a { return nil }
	}
	return fmt.Errorf("Tipe file tidak didukung: %s", contentType)
}

func GenerateKey(userID uuid.UUID, fileName string) string {
	ext := ""
	if idx := strings.LastIndex(fileName, "."); idx >= 0 { ext = fileName[idx:] }
	return fmt.Sprintf("uploads/%s/%s%s", userID.String(), uuid.NewString(), ext)
}

func PresignedURL(endpoint, bucket, key string) string {
	return fmt.Sprintf("%s/%s/%s?presigned=true&expires=900", strings.TrimRight(endpoint, "/"), bucket, key)
}
