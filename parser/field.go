package parser

type Field struct {
	name  string
	alias *string
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Type() int {
	return Field_type
}

func (field *Field) Alias() *string {
	return field.alias
}
