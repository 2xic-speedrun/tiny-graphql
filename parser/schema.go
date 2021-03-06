package parser

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
	if parser.isNextTokenThenSkip("query") || parser.isNextTokenThenSkip("mutation") {
		name := parser.Tokens[parser.index]
		parser.index += 1

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
			schema.Fragments[fragment.Name] = *fragment
		}
		return schema
	} else if parser.isNextTokenThenSkip("{") {
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
			schema.Fragments[fragment.Name] = *fragment
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

			// TODO: variables should be a map.
			for _, object_key := range object.Variables {
				for index_variable_key, variable_key := range schema.Variables {
					if variable_key.Key == object_key.Value {
						schema.Variables[index_variable_key].usage = append(schema.Variables[index_variable_key].usage, object)
						break
					}
				}
			}

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

func (schema *Schema) Inject_variables(variables map[string]interface{}) {
	for index, entry := range schema.Variables {
		if variables[entry.Key] == nil {
			panic("did not find presented variable")
		} else {
			// TODO: THis should just be a map from the start.
			schema.Variables[index].Value = variables[entry.Key].(string)
			for _, value := range schema.Variables[index].usage {
				for index_usage_variables, usage_variables := range value.Variables {
					if usage_variables.Value == entry.Key {
						value.Variables[index_usage_variables].Value = variables[entry.Key].(string)
						break
					}
				}
			}
		}
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
	if parser.isPeekToken("(", 1) || parser.isPeekToken("{", 1) {
		name := parser.Tokens[parser.index]
		parser.index += 1

		variables := parser.ParseArguments()
		condition := parser.ParseConditional()

		parser.index += 1

		object := &Object{
			name:        name,
			Variables:   variables,
			Conditional: condition,
		}

		return object
	} else {
		return nil
	}
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

const (
	Field_type  = 1
	Object_type = 2
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

type Fields interface {
	Type() int
	Name() string
	Alias() *string
}

type Variable struct {
	Key   string
	Value string
	usage []*Object
}

type Conditional struct {
	variant   string
	variables []Variable
}
