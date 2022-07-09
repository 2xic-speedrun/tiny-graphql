package resolver

import (
	"fmt"

	"github.com/2xic-speedrun/tiny-graphql/parser"
)

func (schema *ResolverSchema) Add_field(name string, resolve func() interface{}) {
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
	reference := &data

	if 0 < len(object.Name) {
		data[object.Name] = make(map[string]interface{})
		newReference := data[object.Name].(map[string]interface{})
		reference = &newReference
	}

	schema.fragments = object.Fragments
	//	c := make(chan bool)
	changes := make(chan GoSpeed)
	schema.recursive_resolve(object.Fields, reference, nil, &changes)
	//	fmt.Println(<-c)

	return data
}

func (schema *ResolverSchema) recursive_resolve(object map[string]interface{}, reference *map[string]interface{}, working_object *Object, changes *chan GoSpeed) {
	open_channels := 0
	for field_name, field := range object {
		_, ok := field.(*parser.Object)
		if ok {
			object_field := field.(*parser.Object)
			(*reference)[field_name] = make(map[string]interface{})
			new_reference := (*reference)[field_name].(map[string]interface{})

			if object_field.Fragment_reference != nil {
				// TODO : this should check if the reference is partial
				object_field.Fields = schema.fragments[object_field.Fragment_reference.Name].Fields
			}

			object := *schema.resolve_field(object_field.Name(), working_object)
			object_reference := (object.(*Object))

			open_channels++
			go schema.recursive_resolve(object_field.Fields, &new_reference, object_reference, changes)
		} else {
			field_value := schema.resolve_field(field.(*parser.Field).Name(), working_object)
			if field_value == nil {
				panic("invalid field")
			}
			go (*field_value).Resolve(&field_name, *changes, reference)
			open_channels++
		}
	}

	fmt.Println(open_channels)
	for i := 0; i < open_channels; i++ {
		change := <-*changes
		if change.object {
			continue
		}
		(*change.reference)[change.name] = change.value
	}
	if working_object != nil {
		*changes <- GoSpeed{
			object: true,
		}
	}
}

const (
	field_type  = 1
	object_type = 2
)

func (schema *ResolverSchema) resolve_field(name string, working_object *Object) *Resolvers {
	if working_object == nil {
		value := schema.Resolvers[name]
		return &value
	} else {
		value := working_object.resolves[name]
		return &value
	}
}

type ResolverSchema struct {
	Resolvers map[string]Resolvers
	/*	reference      *map[string]interface{}
		working_object *Object*/
	fragments map[string]parser.FragmentReference
}

type Resolvers interface {
	Name() string
	Resolve(name *string, results chan GoSpeed, reference *map[string]interface{})
	Type() int
	Child() map[string]Resolvers
}

type Object struct {
	name     string
	resolves map[string]Resolvers
}

type Field struct {
	name     string
	resolve  func() interface{}
	resolves map[string]Resolvers
}

func (field *Object) Resolve(name *string, results chan GoSpeed, reference *map[string]interface{}) {
	if name == nil {
		panic("name should not be nil")
	}
	fmt.Println("hey hey", name)
	// ops this could be another object...
	go field.resolves[*name].Resolve(nil, results, reference)
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

func (object *Object) Add_field(name string, resolve func() interface{}) *Field {
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

func (field *Field) Resolve(name *string, results chan GoSpeed, reference *map[string]interface{}) {
	/*	if name != nil {
		panic("Name should be nil on field")
	}*/
	results_value := field.resolve()
	results <- GoSpeed{
		value:     results_value,
		reference: reference,
		name:      *name,
	}
}

type GoSpeed struct {
	reference *map[string]interface{}
	value     interface{}
	name      string
	object    bool
}
