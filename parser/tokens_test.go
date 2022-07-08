package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringParser(t *testing.T) {
	tokens := GetTokens(`
		"long string { } [ 29 ["
	`)
	assert.Equal(t, len(tokens), 1, "Wrong token length")
}
