package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	pw := "password123"
	hash := hashPassword(pw)
	assert.True(t, verifyPassword(hash, pw))
	assert.False(t, verifyPassword(hash, "salah"))
}
func TestGenerateUsername(t *testing.T) {
	u := generateUsername("Budi Santoso", "budi@test.id")
	assert.Contains(t, u, "Budi")
	assert.Greater(t, len(u), 5)
}
func TestVerifyPasswordEmpty(t *testing.T) {
	assert.False(t, verifyPassword("", "any"))
}
