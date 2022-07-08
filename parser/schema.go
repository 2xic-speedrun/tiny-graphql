package parser

import (
	"fmt"
	"regexp"
	"strings"
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
			Name:      name,
			Objects:   objects_and_fields.objects,
			Fields:    objects_and_fields.Fields,
			variant:   parser.Tokens[parser.index],
			Variables: variables,
		}
	} else if parser.Tokens[parser.index] == "{" {
		parser.index += 1
		objects_and_fields := parser.ParseObjectAndFields(getEmptyObject("root"), nil)
		//	fmt.Println(objects_and_fields.Fields[0])
		return Schema{
			Name:    "root",
			Objects: objects_and_fields.objects,
			Fields:  objects_and_fields.Fields,
			variant: "query",
		}
	} else {
		panic("Invalid schema")
	}
}

func (parser *Parser) ParseObjectAndFields(results ObjectAndFields, alias *string) ObjectAndFields {
	results.Alias = alias
	if *parser.Peek(0) == "{" {
		panic("error in parser")
	}
	for parser.index < len(parser.Tokens) {
		peekToken := parser.Peek(0)
		if peekToken != nil && *peekToken == "}" {
			break
		}
		if *parser.Peek(0) == "{" {
			panic("error in parser")
		}
		alias := parser.ParseAlias()
		fragment := parser.ParseFragmentReference()
		object := parser.ParseObject()
		fmt.Printf("%d - %s (%d)\n", parser.index, parser.Tokens[parser.index], len(parser.Tokens[parser.index]))

		if fragment != nil {
			results.fragments = append(results.fragments, *fragment)
		}

		if object == nil {
			field := parser.ParseField()
			if field == nil {
				panic("Something is wrong")
			} else {
				fmt.Printf("Found field %s\n", *field)
				results.Fields = append(results.Fields, Field{
					Name:  *field,
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
		results := getEmptyObject(name)

		parser.index += 1
		results.variables = parser.ParseArguments()
		condition := parser.ParseConditional()
		if condition != nil {
			results.conditional = *condition
		}
		parser.index += 1

		return &results
	} else {
		return nil
	}
}

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
				fmt.Println("Found dict")
			}
		}

		if key != nil && terminator != nil && value != nil {
			variables = append(variables, Variable{
				key:   *key,
				value: *value,
			})
		} else {
			panic("Invalid arguments")
		}
	},
		parser.DictAndArrayTerminatorFunction)

	return variables
}

func (parser *Parser) ParseConditional() *Conditional {
	if parser.isNextToken("@") {
		if parser.isNextToken("skip") {
			return &Conditional{
				variant:   "skip",
				variables: parser.ParseArguments(),
			}
		} else if parser.isNextToken("include") {
			return &Conditional{
				variant:   "include",
				variables: parser.ParseArguments(),
			}
		}
	}
	return nil
}

func (parser *Parser) ParseArray() *string {
	results := ""
	parser.ParseScope("[", "]", func() {
		results += *parser.Read()
	},
		parser.DictAndArrayTerminatorFunction)
	if len(results) == 0 {
		return nil
	}
	return &results
}

func (parser *Parser) ParseDict() *string {
	results := ""
	parser.ParseScope("{", "}", func() {
		results += *parser.Read()
	},
		parser.DictAndArrayTerminatorFunction)
	if len(results) == 0 {
		return nil
	}
	return &results
}

func (parser *Parser) ParseScope(init string, terminator string, callback func(), terminatorFunction func(terminator string) bool) bool {
	peekArguments := parser.Peek(0)
	if peekArguments != nil && *peekArguments == init {
		parser.index += 1
		for true {
			if terminatorFunction(terminator) {
				break
			}
			callback()
		}
		parser.index += 1
		return true
	}
	return false
}

func (parser *Parser) DictAndArrayTerminatorFunction(terminator string) bool {
	finished := parser.Peek(0)
	if finished != nil && *finished == terminator {
		return true
	} else if finished != nil && *finished == "," {
		parser.index += 1
	}
	return false
}

func (parser *Parser) ParseLiteral() {
	parser.Peek(0)
	parser.index += 1
}

func (parser *Parser) ParseAlias() *string {
	// alias : field <- leaving with pointer here
	peekedToken := parser.Peek(1)
	if peekedToken != nil && *peekedToken == ":" {
		results := parser.Read()
		parser.index += 1
		return results
	}
	return nil
}

func (parser *Parser) ParseFragment() *Fragment {
	if parser.isNextToken("fragment") {
		fragmentName := parser.Read()
		if parser.isNextToken("on") {
			fmt.Println("peek ?", parser.Peek(0))
			return &Fragment{
				name: *fragmentName,
				on:   *parser.Read(),
				fields: parser.ParseObjectAndFields(
					ObjectAndFields{
						Name:      "fragment",
						Alias:     nil,
						objects:   []ObjectAndFields{},
						Fields:    []Field{},
						variables: []Variable{},
					},
					nil,
				),
			}
		}
	}

	return nil
}

func (parser *Parser) ParseFragmentReference() *FragmentReference {
	if parser.isNextTokenSequence([]string{".", ".", "."}) {
		if parser.isNextToken("on") {
			object := *parser.Read()
			parser.index += 1
			return &FragmentReference{
				object: object,
				child: parser.ParseObjectAndFields(
					getEmptyObject(""),
					nil,
				),
			}
		}
		return &FragmentReference{
			name: *parser.Read(),
		}
	}
	return nil
}

func (parser *Parser) ParseField() *string {
	return parser.Read()
}

func getEmptyObject(name string) ObjectAndFields {
	return ObjectAndFields{
		Name:    name,
		objects: []ObjectAndFields{},
		Fields:  []Field{},
		Alias:   nil,
	}
}

type Schema struct {
	Name      string
	variant   string
	Variables []Variable
	Objects   []ObjectAndFields
	Fields    []Field
}

type Object struct {
	name   string
	alias  string
	fields []Field
}

type ObjectAndFields struct {
	Name        string
	Alias       *string
	variables   []Variable
	objects     []ObjectAndFields
	Fields      []Field
	fragments   []FragmentReference
	conditional Conditional
}

type Field struct {
	Name  string
	alias *string
}

type Fragment struct {
	name   string
	on     string
	fields ObjectAndFields
}

type FragmentReference struct {
	object string
	name   string
	child  ObjectAndFields
}

type Variable struct {
	key   string
	value string
}

type Conditional struct {
	variant   string
	variables []Variable
}
