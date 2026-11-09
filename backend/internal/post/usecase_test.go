package post

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHotScore(t *testing.T) {
	// hot = (up - down) / (age+2)^1.5
	score := HotScore(100, 10, 2)
	assert.Greater(t, score, 0.0)
	young := HotScore(10, 0, 1)
	old := HotScore(10, 0, 24)
	assert.Greater(t, young, old)
	neg := HotScore(2, 10, 1)
	assert.Less(t, neg, 0.0)
	assert.InDelta(t, 10/2.8284271247461903, HotScore(10, 0, 0), 0.0001)
}
