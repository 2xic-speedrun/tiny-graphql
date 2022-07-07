package parser

import (
	"strings"
)

func GetTokens(schema string) []string {
	tokens := TokensParser{
		tokens:       []string{},
		currentToken: "",
	}
	terminators := []string{"(", ")", "{", "}", ":", "[", "]", ",", "@", "."}
	index := 0

	for index < len(schema) {
		char := string(schema[index])
		if char == "#" {
			tokens.addTokenAndClearCurrentToken(tokens.currentToken)
			for index < len(schema) && string(schema[index]) != "\n" {
				index++
			}
		} else if char == "\"" {
			index += 1
			tokens.addTokenAndClearCurrentToken(tokens.currentToken)
			for index < len(schema) {
				current := string(schema[index])
				if current != "\"" {
					tokens.currentToken += string(schema[index])
					index++
				} else {
					break
				}
			}
			tokens.addTokenAndClearCurrentToken(tokens.currentToken)
		} else if contains(terminators, char) {
			tokens.addTokenAndClearCurrentToken(tokens.currentToken)
			tokens.addTokenAndClearCurrentToken(char)
		} else if len(strings.TrimSpace(char)) == 0 {
			tokens.addTokenAndClearCurrentToken(tokens.currentToken)
		} else {
			tokens.currentToken += char
		}
		index++
	}
	tokens.addTokenAndClearCurrentToken(tokens.currentToken)

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

func (tokenParser *TokensParser) addTokenAndClearCurrentToken(token string) {
	if 0 < len(token) {
		tokenParser.tokens = append(tokenParser.tokens, token)
	}
	tokenParser.currentToken = ""
}

type TokensParser struct {
	tokens       []string
	currentToken string
}
