package parser

import (
	"regexp"
	"strings"
)

func (parser *Parser) ParseArguments() []Variable {
	variables := []Variable{}
	parser.ParseScope("(", ")", func() {
		key := parser.Peek(0)
		terminator := parser.Peek(1)
		parser.index += 2
		var value *string

		isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
		isNumeric := regexp.MustCompile(`^0|[1-9]\d*$`).MatchString
		if isAlpha(*parser.Peek(0)) || isNumeric(*parser.Peek(0)) {
			value = parser.Peek(0)
			parser.index++
		} else if strings.HasPrefix(*parser.Peek(0), "$") {
			value = parser.Peek(0)
			parser.index++
		} else if strings.HasPrefix(*parser.Peek(0), "\"") {
			value = parser.Peek(0)
			parser.index++
		} else {
			value = parser.ParseArray()
			if value == nil {
				value = parser.ParseDict()
			}
		}

		if key != nil && terminator != nil && value != nil {
			variables = append(variables, Variable{
				Key:   *key,
				Value: *value,
			})
		} else {
			panic("Invalid arguments")
		}
	},
		parser.DictAndArrayTerminatorFunction)

	return variables
}
