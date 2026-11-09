package vote

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolToInt(t *testing.T) {
	assert.Equal(t, 1, boolToInt(true))
	assert.Equal(t, 0, boolToInt(false))
}
