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

func (schema *ResolverSchema) Resolve(object parser.Schema) map[string]interface{} {
	data := map[string]interface{}{}
	//var reference *map[string]interface{}
	schema.reference = &data

	if 0 < len(object.Name) {
		data[object.Name] = make(map[string]interface{})
		newReference := data[object.Name].(map[string]interface{})
		schema.reference = &newReference
	}

	schema.recursive_resolve(object.Fields)

	return data
}

func (schema *ResolverSchema) recursive_resolve(object map[string]interface{}) {
	for field_name, field := range object {
		_, ok := field.(*parser.Object)
		if ok {
			object_field := field.(*parser.Object)
			old_reference := schema.reference
			(*schema.reference)[field_name] = make(map[string]interface{})
			new_reference := (*schema.reference)[field_name].(map[string]interface{})
			schema.reference = &new_reference

			last_object_reference := schema.working_object

			object := *schema.resolve_field(object_field.Name())
			object_reference := (object.(*Object))
			schema.working_object = object_reference

			schema.recursive_resolve(object_field.Fields)

			schema.working_object = last_object_reference

			schema.reference = old_reference

		} else {
			field_value := schema.resolve_field(field.(*parser.Field).Name())
			if field_value == nil {
				panic("invalid field")
			}
			value := (*field_value).Resolve(nil)
			(*schema.reference)[field_name] = value
		}
	}
}

const (
	field_type  = 1
	object_type = 2
)

func (schema *ResolverSchema) resolve_field(name string) *Resolvers {
	if schema.working_object == nil {
		value := schema.Resolvers[name]
		return &value
	} else {
		value := schema.working_object.resolves[name]
		return &value
	}
}

type ResolverSchema struct {
	Resolvers      map[string]Resolvers
	reference      *map[string]interface{}
	working_object *Object
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

func (object *Object) Add_object(child_object *Object) *Object {
	if object.resolves == nil {
		object.resolves = make(map[string]Resolvers)
	}
	object.resolves[child_object.name] = child_object
	return child_object
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
