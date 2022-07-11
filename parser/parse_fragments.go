package parser

import "fmt"

func (parser *Parser) ParseFragmentReference() *FragmentReference {
	if parser.isNextTokenSequence([]string{".", ".", "."}) {
		if parser.isNextTokenThenSkip("on") {
			object := *parser.Read()
			fragment_reference := &FragmentReference{
				Name:   object,
				Fields: make(map[string]interface{}),
			}
			fragment_reference.reference = &fragment_reference.Fields

			parser.index += 1
			parser.ConstructFragmentReference((fragment_reference))
			return fragment_reference
		}

		return &FragmentReference{
			Name: *parser.Read(),
		}
	}
	return nil
}

func (parser *Parser) ConstructFragmentReference(fragment_reference *FragmentReference) {
	parser.BaseParser(func(alias *string, object *Object, _fragment_reference *FragmentReference) {
		current_map := *fragment_reference.reference
		if _fragment_reference != nil {
			parser.ConstructFragmentReference(
				_fragment_reference,
			)
			current_map[_fragment_reference.Name] = _fragment_reference
		}

		if alias != nil {
			panic("can a fragment have a alias ? ")
		}
		if object != nil {
			// TODO: this can just reuse the base parser.
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
				Name:   *fragment_name,
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

type FragmentReference struct {
	object    string
	Name      string
	Fields    map[string]interface{}
	reference *map[string]interface{}
}
