package media

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	assert.NoError(t, Validate("image", 5*1024*1024, "image/jpeg"))
	assert.Error(t, Validate("image", 20*1024*1024, "image/jpeg"))
	assert.Error(t, Validate("image", 1*1024*1024, "video/mp4"))
	assert.NoError(t, Validate("gif", 5*1024*1024, "image/gif"))
	assert.NoError(t, Validate("video", 20*1024*1024, "video/mp4"))
}
func TestGenerateKey(t *testing.T) {
	k := GenerateKey(uuid.New(), "meme.jpg")
	assert.Contains(t, k, ".jpg")
}
func TestPresignedURL(t *testing.T) {
	url := PresignedURL("http://localhost:9000", "memekita-media", "uploads/a/b.jpg")
	assert.Contains(t, url, "memekita-media")
}
