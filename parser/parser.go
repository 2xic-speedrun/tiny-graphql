package parser

import (
	"fmt"
)

func Parse(schema string) QuerySchema {
	tokens := GetTokens(schema)

	fmt.Println(tokens)
	parser := Parser{
		Tokens: tokens,
		index:  0,
	}
	return parser.ParseSchema()
}

func (parser *Parser) ParseSchema() QuerySchema {
	if parser.Tokens[parser.index] == "query" {
		name := parser.Tokens[parser.index+1]
		parser.index += 3

		fmt.Printf("Current token %s\n", parser.Tokens[parser.index])

		objects_and_fields := parser.ParseObjectAndFields("root")

		return QuerySchema{
			name:    name,
			objects: objects_and_fields.objects,
		}
	} else {
		panic("Invalid schema")
	}
}

func (parser *Parser) ParseObjectAndFields(name string) ObjectAndFields {
	results := ObjectAndFields{
		name:    name,
		objects: []ObjectAndFields{},
		fields:  []string{},
	}
	for parser.index < len(parser.Tokens) {
		peekToken := parser.Peek(0)
		if peekToken != nil && *peekToken == "}" {
			break
		}

		object := parser.ParseObject()
		fmt.Printf("%d - %s (%d)\n", parser.index, parser.Tokens[parser.index], len(parser.Tokens[parser.index]))

		if object == nil {
			field := parser.ParseField()
			if field == nil {
				panic("Something is wrong")
			} else {
				fmt.Printf("Found field %s\n", *field)
				results.fields = append(results.fields, *field)
			}
		} else {
			fmt.Printf("object %s, current token %s\n", *object, parser.Tokens[parser.index])
			results.objects = append(results.objects, parser.ParseObjectAndFields(*object))
		}

	}

	return results
}

func (parser *Parser) ParseObject() *string {
	peeked := parser.Peek(1)
	if peeked != nil && (*peeked == "(" || *peeked == "{") {
		results := parser.Tokens[parser.index]
		parser.index += 1

		// parser arguments
		peekArguments := parser.Peek(0)
		if peekArguments != nil && *peekArguments == "(" {
			parser.index += 1

			// push arguments here
			parser.ParseArguments()

			parser.index += 2
		} else {
			parser.index += 1
		}

		return &results
	} else {
		return nil
	}
}

func (parser *Parser) ParseArguments() {
	for true {
		finished := parser.Peek(0)
		fmt.Printf("finished ? %s\n", *finished)
		if finished != nil && *finished == ")" {
			break
		}

		key := parser.Peek(0)
		terminator := parser.Peek(1)
		value := parser.Peek(2)
		if key != nil && terminator != nil && value != nil {
			// push
			parser.index += 3
		} else {
			panic("Invalid arguments")
		}
	}
}

func (parser *Parser) ParseField() *string {
	results := parser.Tokens[parser.index]
	parser.index += 1
	return &results
}

func (parser *Parser) Peek(length int) *string {
	if (parser.index + length) < len(parser.Tokens) {
		return &parser.Tokens[parser.index+length]
	}
	return nil
}

type QuerySchema struct {
	name    string
	objects []ObjectAndFields
}

type Object struct {
	name   string
	fields []string
}

type Parser struct {
	Tokens []string
	index  int
}

type ObjectAndFields struct {
	name    string
	objects []ObjectAndFields
	fields  []string
}
