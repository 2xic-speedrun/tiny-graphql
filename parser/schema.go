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
			Fragments: map[string]FragmentReference{}}
		schema.ParseObjectAndFields()
		fragment := schema.parser.ParseFragment()
		if fragment != nil {
			schema.Fragments[fragment.name] = *fragment
		}
		return schema
	} else if parser.Tokens[parser.index] == "{" {
		parser.index += 1

		schema := Schema{
			Name:      "root",
			variant:   parser.Tokens[parser.index],
			Fields:    data,
			reference: &data,
			parser:    *parser,
			Fragments: map[string]FragmentReference{},
		}
		schema.ParseObjectAndFields()

		fragment := schema.parser.ParseFragment()
		if fragment != nil {
			schema.Fragments[fragment.name] = *fragment
		}

		return schema
	} else {
		panic("Invalid schema")
	}
}

func (schema *Schema) ParseObjectAndFields() {
	schema.parser.BaseParser(func(
		alias *string,
		object *Object,
		fragment_reference *FragmentReference,
	) {
		current_map := *schema.reference

		if object != nil {
			object.alias = alias
			object.Fragment_reference = fragment_reference
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
	})
}

func (parser *Parser) BaseParser(callback func(alias *string, object *Object, fragment_reference *FragmentReference)) {
	if *parser.Peek(0) == "{" {
		panic("error in parser")
	}
	lastIndex := parser.index
	for parser.index < len(parser.Tokens) {
		peekToken := parser.Peek(0)
		if peekToken != nil && *peekToken == "}" {
			break
		}
		if *parser.Peek(0) == "{" {
			panic("error in parser")
		}

		alias := parser.ParseAlias()
		object := parser.ParseObject()
		fragment_reference := parser.ParseFragmentReference()

		callback(alias, object, fragment_reference)

		if parser.index == lastIndex {
			panic("something is wrong")
		} else {
			lastIndex = parser.index
		}
	}
	parser.index += 1
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
		parser.index += 1

		variables := parser.ParseArguments()
		condition := parser.ParseConditional()

		parser.index += 1

		return &Object{
			name:        name,
			Variables:   variables,
			Conditional: condition,
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

func (parser *Parser) ParseConditional() *Conditional {
	if parser.isNextTokenThenSkip("@") {
		if parser.isNextTokenThenSkip("skip") {
			return &Conditional{
				variant:   "skip",
				variables: parser.ParseArguments(),
			}
		} else if parser.isNextTokenThenSkip("include") {
			return &Conditional{
				variant:   "include",
				variables: parser.ParseArguments(),
			}
		}
	}
	return nil
}

func (parser *Parser) ParseFragmentReference() *FragmentReference {
	if parser.isNextTokenSequence([]string{".", ".", "."}) {
		if parser.isNextTokenThenSkip("on") {
			object := *parser.Read()
			fragment_reference := &FragmentReference{
				name:   object,
				Fields: make(map[string]interface{}),
			}
			fragment_reference.reference = &fragment_reference.Fields

			parser.index += 1
			parser.ConstructFragmentReference((fragment_reference))
			return fragment_reference
		}

		return &FragmentReference{
			name: *parser.Read(),
		}
	}
	return nil
}

func (parser *Parser) ConstructFragmentReference(fragment_reference *FragmentReference) {
	parser.BaseParser(func(alias *string, object *Object, _fragment_reference *FragmentReference) {
		current_map := *fragment_reference.reference
		if _fragment_reference != nil {
			fmt.Println(_fragment_reference.name)
			parser.ConstructFragmentReference(
				_fragment_reference,
			)
			current_map[_fragment_reference.name] = _fragment_reference
		}

		if alias != nil {
			panic("can a fragment has a alias ? ")
		}
		if object != nil {
			panic("not implemented")
		} else {
			field := parser.ParseField()
			current_map[*field] = &Field{
				name:  *field,
				alias: alias,
			}
		}
	})
}

func (parser *Parser) ParseFragment() *FragmentReference {
	if parser.isNextTokenThenSkip("fragment") {
		fragment_name := parser.Read()
		if parser.isNextTokenThenSkip("on") {
			on_object := parser.Read()
			fragment_reference := &FragmentReference{
				object: *on_object,
				name:   *fragment_name,
				Fields: make(map[string]interface{}),
			}
			fragment_reference.reference = &fragment_reference.Fields
			parser.index += 1
			parser.ConstructFragmentReference(fragment_reference)
			return fragment_reference
		}
	}
	return nil
}

const (
	field_type  = 1
	object_type = 2
)

type Schema struct {
	Name      string
	variant   string
	Variables []Variable
	Fragments map[string]FragmentReference
	Fields    map[string]interface{}

	reference *map[string]interface{}
	parser    Parser
}

/*
type Fragments struct {
	Name   string
	Fields map[string]interface{}
}
*/
type Fields interface {
	Type() int
	Name() string
	Alias() *string
}

type Object struct {
	name               string
	alias              *string
	Fields             map[string]interface{}
	Variables          []Variable
	Conditional        *Conditional
	Fragment_reference *FragmentReference
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

type Conditional struct {
	variant   string
	variables []Variable
}

type FragmentReference struct {
	object    string
	name      string
	Fields    map[string]interface{}
	reference *map[string]interface{}
}
