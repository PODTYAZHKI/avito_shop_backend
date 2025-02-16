package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	generator := NewGenerator([]byte("secret"))
	token, err := generator.Generate("testuser")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
