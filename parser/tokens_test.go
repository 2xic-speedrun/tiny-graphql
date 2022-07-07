package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringParser(t *testing.T) {
	tokens := GetTokens(`
		"long string { } [ 29 ["
	`)
	fmt.Println(tokens)
	assert.Equal(t, len(tokens), 1, "Wrong token length")
}
