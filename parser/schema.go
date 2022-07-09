package parser

import (
	"fmt"
	"regexp"
	"strings"
)

func Parse(schema string) Schema {
	tokens := GetTokens(schema)

	parser := Parser{
		Tokens: tokens,
		index:  0,
	}
	return parser.parseSchema()
}

func (parser *Parser) parseSchema() Schema {
	data := map[string]interface{}{}
	if parser.Tokens[parser.index] == "query" || parser.Tokens[parser.index] == "mutation" {
		name := parser.Tokens[parser.index+1]
		parser.index += 2
		// parser arguments
		variables := parser.ParseArguments()
		parser.index += 1

		schema := Schema{
			Name:      name,
			variant:   parser.Tokens[parser.index],
			Fields:    data,
			reference: &data,
			parser:    *parser,
			Variables: variables,
		}
		schema.ParseObjectAndFields()
		return schema
	} else if parser.Tokens[parser.index] == "{" {
		parser.index += 1

		schema := Schema{
			Name:      "root",
			variant:   parser.Tokens[parser.index],
			Fields:    data,
			reference: &data,
			parser:    *parser,
		}
		schema.ParseObjectAndFields()
		return schema
	} else {
		panic("Invalid schema")
	}
}

func (schema *Schema) ParseObjectAndFields() {
	//	results.Alias = alias
	fmt.Println("HELLLLOOO")
	if *schema.parser.Peek(0) == "{" {
		panic("error in parser")
	}
	fmt.Println("HELLLLOOO")
	lastIndex := schema.parser.index
	for schema.parser.index < len(schema.parser.Tokens) {
		peekToken := schema.parser.Peek(0)
		if peekToken != nil && *peekToken == "}" {
			break
		}
		if *schema.parser.Peek(0) == "{" {
			panic("error in parser")
		}

		alias := schema.parser.ParseAlias()
		object := schema.parser.ParseObject()
		current_map := *schema.reference
		fmt.Println(object)

		if object != nil {
			object.alias = alias
			old_reference := schema.reference
			current_map[object.Name()] = object

			current_map[object.Name()].(*Object).Fields = make(map[string]interface{})
			object_field_reference := current_map[object.Name()].(*Object).Fields
			schema.reference = &object_field_reference
			schema.ParseObjectAndFields()

			schema.reference = old_reference
		} else {
			field := schema.parser.ParseField()
			current_map[*field] = &Field{
				name:  *field,
				alias: alias,
			}
		}

		if schema.parser.index == lastIndex {
			panic("something is wrong")
		} else {
			lastIndex = schema.parser.index
		}

		/*
			fragment := parser.ParseFragmentReference()
			object := parser.ParseObject()
			fmt.Printf("%d - %s (%d)\n", parser.index, parser.Tokens[parser.index], len(parser.Tokens[parser.index]))

			if fragment != nil {
				results.fragments = append(results.fragments, *fragment)
			}

			if object == nil {
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
		*/
	}
}

func (parser *Parser) ParseField() *string {
	return parser.Read()
}

func (parser *Parser) ParseAlias() *string {
	// alias : field <- leaving with pointer here
	if parser.isPeekToken(":", 1) {
		results := parser.Read()
		parser.index += 1
		return results
	}
	return nil
}

func (parser *Parser) ParseObject() *Object {
	//	peeked := parser.Peek(1)
	if parser.isPeekToken("(", 1) || parser.isPeekToken("{", 1) {
		name := parser.Tokens[parser.index]
		//	results := getEmptyObject(name)

		parser.index += 1
		//results.variables =
		variables := parser.ParseArguments()
		/*
			condition := parser.ParseConditional()
			if condition != nil {
				results.conditional = *condition
			}
		*/
		parser.index += 1

		return &Object{
			name:      name,
			Variables: variables,
		}
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

func (parser *Parser) DictAndArrayTerminatorFunction(terminator string) bool {
	finished := parser.Peek(0)
	if finished != nil && *finished == terminator {
		return true
	} else if finished != nil && *finished == "," {
		parser.index += 1
	}
	return false
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

/*

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


func (parser *Parser) ParseLiteral() {
	parser.Peek(0)
	parser.index += 1
}

func (parser *Parser) ParseFragment() *Fragment {
	if parser.isNextToken("fragment") {
		fragmentName := parser.Read()
		if parser.isNextToken("on") {
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

func getEmptyObject(name string) ObjectAndFields {
	return ObjectAndFields{
		Name:    name,
		objects: []ObjectAndFields{},
		Fields:  []Field{},
		Alias:   nil,
	}
}
*/

const (
	field_type  = 1
	object_type = 2
)

type Schema struct {
	Name      string
	variant   string
	Fields    map[string]interface{}
	Variables []Variable

	reference *map[string]interface{}
	parser    Parser
}

type Fields interface {
	Type() int
	Name() string
	Alias() *string
}

type Object struct {
	name      string
	alias     *string
	Fields    map[string]interface{}
	Variables []Variable
}

func (object *Object) Name() string {
	return object.name
}

func (object *Object) Type() int {
	return object_type
}

func (object *Object) Alias() *string {
	return object.alias
}

type Field struct {
	name  string
	alias *string
	/*	Name  string
		alias *string*/
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Type() int {
	return field_type
}

func (field *Field) Alias() *string {
	return field.alias
}

type Variable struct {
	key   string
	value string
}

/*
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

type Conditional struct {
	variant   string
	variables []Variable
}
*/
