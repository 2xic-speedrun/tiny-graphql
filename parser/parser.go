package parser

import (
	"fmt"
)

func Parse(schema string) Schema {
	tokens := GetTokens(schema)

	fmt.Println(tokens)
	parser := Parser{
		Tokens: tokens,
		index:  0,
	}
	return parser.ParseSchema()
}

func (parser *Parser) ParseSchema() Schema {
	if parser.Tokens[parser.index] == "query" || parser.Tokens[parser.index] == "mutation" {
		name := parser.Tokens[parser.index+1]
		parser.index += 2
		// parser arguments
		variables := parser.ParseArguments()
		parser.index += 1
		objects_and_fields := parser.ParseObjectAndFields(getEmptyObject(name), nil)

		return Schema{
			name:      name,
			objects:   objects_and_fields.objects,
			variant:   parser.Tokens[parser.index],
			variables: variables,
		}
	} else if parser.Tokens[parser.index] == "{" {
		objects_and_fields := parser.ParseObjectAndFields(getEmptyObject("root"), nil)
		parser.index += 1
		return Schema{
			name:    "root",
			objects: objects_and_fields.objects,
			variant: "query",
		}
	} else {
		panic("Invalid schema")
	}
}

func (parser *Parser) ParseObjectAndFields(results ObjectAndFields, alias *string) ObjectAndFields {
	/*	results := ObjectAndFields{
			name:    name,
			objects: []ObjectAndFields{},
			fields:  []Field{},
			alias:   alias,
		}
	*/
	results.alias = alias
	for parser.index < len(parser.Tokens) {
		peekToken := parser.Peek(0)
		if peekToken != nil && *peekToken == "}" {
			break
		}
		alias := parser.ParseAlias()
		object := parser.ParseObject()
		fmt.Printf("%d - %s (%d)\n", parser.index, parser.Tokens[parser.index], len(parser.Tokens[parser.index]))

		if object == nil {
			field := parser.ParseField()
			if field == nil {
				panic("Something is wrong")
			} else {
				fmt.Printf("Found field %s\n", *field)
				results.fields = append(results.fields, Field{
					name:  *field,
					alias: alias,
				})
			}
		} else {
			//fmt.Printf("object %s, current token %s\n", *object, parser.Tokens[parser.index])
			results.objects = append(results.objects, parser.ParseObjectAndFields(*object, alias))
		}

	}

	return results
}

func (parser *Parser) ParseObject() *ObjectAndFields {
	peeked := parser.Peek(1)
	if peeked != nil && (*peeked == "(" || *peeked == "{") {
		name := parser.Tokens[parser.index]
		results := ObjectAndFields{
			name:      name,
			objects:   []ObjectAndFields{},
			fields:    []Field{},
			alias:     nil,
			variables: []Variable{},
		}
		parser.index += 1

		// parser arguments
		results.variables = parser.ParseArguments()
		parser.index += 1

		return &results
	} else {
		return nil
	}
}

func (parser *Parser) ParseArguments() []Variable {
	peekArguments := parser.Peek(0)
	variables := []Variable{}
	if peekArguments != nil && *peekArguments == "(" {
		parser.index += 1
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
				variables = append(variables, Variable{
					key:   *key,
					value: *value,
				})
				// push
				parser.index += 3
			} else {
				panic("Invalid arguments")
			}
		}
		parser.index += 1
	}
	return variables
}

func (parser *Parser) ParseAlias() *string {
	peekedToken := parser.Peek(1)
	if peekedToken != nil && *peekedToken == ":" {
		results := parser.Tokens[parser.index]
		parser.index += 2
		return &results
	}
	return nil
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

func getEmptyObject(name string) ObjectAndFields {
	return ObjectAndFields{
		name:    name,
		objects: []ObjectAndFields{},
		fields:  []Field{},
		alias:   nil,
	}
}

type Schema struct {
	name      string
	variant   string
	variables []Variable
	objects   []ObjectAndFields
}

type Object struct {
	name   string
	alias  string
	fields []Field
}

type Parser struct {
	Tokens []string
	index  int
}

type ObjectAndFields struct {
	name      string
	alias     *string
	variables []Variable
	objects   []ObjectAndFields
	fields    []Field
}

type Field struct {
	name  string
	alias *string
}

type Variable struct {
	key   string
	value string
}
