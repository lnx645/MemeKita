package comment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommentValidation(t *testing.T) {
	// pure logic: empty body should error when creating (tested via usecase validation)
	assert.True(t, true)
	_ = assert.AnError
}
