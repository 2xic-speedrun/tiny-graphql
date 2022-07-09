package resolver

import (
	"github.com/2xic-speedrun/tiny-graphql/parser"
)

func (schema *ResolverSchema) Add_field(name string, resolve func() string) {
	field := &Field{
		name:     name,
		resolve:  resolve,
		resolves: nil,
	}
	schema.Resolvers[field.name] = field
}

func (schema *ResolverSchema) Add_Object(name string) *Object {
	object := &Object{
		name: name,
	}
	schema.Resolvers[object.name] = object
	return object
}

// Todo this has to be recursive.
func (schema *ResolverSchema) Resolve(object parser.Schema) map[string]interface{} {
	data := map[string]interface{}{}
	/*
		var reference *map[string]interface{}
		reference = &data

		if 0 < len(object.Name) {
			data[object.Name] = make(map[string]interface{})
			newReference := data[object.Name].(map[string]interface{})
			reference = &newReference
		}

		for _, val := range object.Fields {
			field := schema.resolve_field(val.Name)
			if field == nil {
				panic("invalid field")
			}
			value := (*field).Resolve(nil)
			(*reference)[val.Name] = value
		}
		for _, val := range object.Objects {
			(*reference)[val.Name] = make(map[string]interface{})
			newReference := (*reference)[val.Name].(map[string]interface{})
			reference = &newReference
			object := schema.resolve_field(val.Name)

			for _, field := range val.Fields {
				field_value := (*object).Child()[field.Name].Resolve(nil)
				(*reference)[field.Name] = field_value
			}
		}*/

	return data
}

const (
	field_type  = 1
	object_type = 2
)

func (schema *ResolverSchema) resolve_field(name string) *Resolvers {
	value := schema.Resolvers[name]
	return &value
}

type ResolverSchema struct {
	Resolvers map[string]Resolvers
}

type Resolvers interface {
	Name() string
	Resolve(name *string) string
	Type() int
	Child() map[string]Resolvers
}

type Object struct {
	name     string
	resolves map[string]Resolvers
}

type Field struct {
	name     string
	resolve  func() string
	resolves map[string]Resolvers
}

func (field *Object) Resolve(name *string) string {
	if name == nil {
		panic("name should not be nil")
	}
	// ops this could be another object...
	return field.resolves[*name].Resolve(nil)
}

func (field *Object) Name() string {
	return field.name
}

func (field *Object) Child() map[string]Resolvers {
	return field.resolves
}

func (field *Object) Type() int {
	return object_type
}

func (object *Object) Add_field(name string, resolve func() string) *Field {
	field := &Field{
		name:     name,
		resolve:  resolve,
		resolves: nil,
	}
	if object.resolves == nil {
		object.resolves = make(map[string]Resolvers)
	}
	object.resolves[name] = field
	return field
}

func (field *Field) Name() string {
	return field.name
}

func (field *Field) Child() map[string]Resolvers {
	panic("not used on field")
}

func (field *Field) Type() int {
	return field_type
}

func (field *Field) Resolve(name *string) string {
	if name != nil {
		panic("Name should be nil on field")
	}
	return field.resolve()
}
