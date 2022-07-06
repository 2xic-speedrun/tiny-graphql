package parser

import (
	"strings"
)

func GetTokens(schema string) []string {
	tokens := TokensParser{
		tokens:       []string{},
		currentToken: "",
	}
	terminators := []string{"(", ")", "{", "}", ":", "[", "]", ","}
	index := 0

	for index < len(schema) {
		char := string(schema[index])
		if char == "#" {
			tokens.addToken(tokens.currentToken)
			for index < len(schema) && string(schema[index]) != "\n" {
				index++
			}
		} else if contains(terminators, char) {
			tokens.addToken(tokens.currentToken)
			tokens.addToken(char)
			index++
		} else if len(strings.TrimSpace(char)) == 0 {
			tokens.addToken(tokens.currentToken)
			index++
		} else {
			tokens.currentToken += char
			index++
		}
	}
	return tokens.tokens
}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (tokenParser *TokensParser) addToken(token string) {
	if 0 < len(token) {
		tokenParser.tokens = append(tokenParser.tokens, token)
	}
	tokenParser.currentToken = ""
}

type TokensParser struct {
	tokens       []string
	currentToken string
}
