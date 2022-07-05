package parser

import (
	"fmt"
	"strings"
)

func GetTokens(schema string) []string {
	var tokens []string
	terminators := []string{"(", ")", "{", "}", ":"}
	currentToken := ""
	index := 0
	for index < len(schema) {
		char := string(schema[index])
		if contains(terminators, char) {
			if 0 < len(currentToken) {
				tokens = append(tokens, currentToken)
				fmt.Println(currentToken)
			}
			tokens = append(tokens, char)
			currentToken = ""
			index++
		} else if len(strings.TrimSpace(char)) == 0 {
			if 0 < len(currentToken) {
				tokens = append(tokens, currentToken)
				fmt.Println(currentToken)
			}
			currentToken = ""
			index++
		} else {
			currentToken += char
			index++
		}
	}
	return tokens
}

func contains(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
