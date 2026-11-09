package reaction

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidReactions(t *testing.T) {
	for _, r := range []string{"ngakak", "relate", "anjay", "sedih", "marah"} {
		assert.True(t, valid[r])
	}
	assert.False(t, valid["like"])
}
