package parser

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
	return Object_type
}

func (object *Object) Alias() *string {
	return object.alias
}
