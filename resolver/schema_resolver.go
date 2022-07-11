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

func (schema *ResolverSchema) Add_Object_fetch(name string, fetch func(map[string]interface{}) interface{}) *Object {
	object := &Object{
		name:  name,
		fetch: fetch,
	}
	schema.Resolvers[object.name] = object
	return object
}

func (schema *ResolverSchema) Resolve(object parser.Schema) map[string]interface{} {
	data := map[string]interface{}{}
	reference := &data

	if 0 < len(object.Name) {
		data[object.Name] = make(map[string]interface{})
		newReference := data[object.Name].(map[string]interface{})
		reference = &newReference
	}

	schema.fragments = object.Fragments
	changes := make(chan ContextFieldReference)
	schema.recursive_resolve(object.Fields,
		Context{
			reference:      reference,
			working_object: nil,
			changes:        &changes,
		},
	)

	return data
}

func (schema *ResolverSchema) recursive_resolve(object map[string]interface{}, context Context) {
	open_channels := 0
	for field_name, field := range object {
		_, ok := field.(*parser.Object)
		if ok {
			object_field := field.(*parser.Object)
			(*context.reference)[field_name] = make(map[string]interface{})
			new_reference := (*context.reference)[field_name].(map[string]interface{})

			if object_field.Fragment_reference != nil {
				// TODO : this should check if the reference is partial
				object_field.Fields = schema.fragments[object_field.Fragment_reference.Name].Fields
			}

			object := *schema.resolve_field(object_field.Name(), context.working_object)
			object_reference := (object.(*Object))

			if object_reference.fetch != nil {
				// TODO: This should just be a map from the start.
				variables_map := make(map[string]interface{})
				for _, variable := range object_field.Variables {
					variables_map[variable.Key] = variable.Value
				}
				object_reference.value = object_reference.fetch(variables_map)
				if object_reference.value == nil {
					(*context.reference)[field_name] = nil
					continue
				}
			}
			open_channels++
			go schema.recursive_resolve(object_field.Fields,
				Context{
					reference:      &new_reference,
					working_object: object_reference,
					changes:        context.changes,
				},
			)
		} else {
			field_value := schema.resolve_field(field.(*parser.Field).Name(), context.working_object)
			if (*field_value) == nil {
				panic(fmt.Sprintf("invalid field %s\n", field_name))
			}
			go (*field_value).Resolve(&field_name, &context)
			open_channels++
		}
	}
	for i := 0; i < open_channels; i++ {
		change := <-*context.changes
		if change.object {
			continue
		}
		(*change.reference)[change.name] = change.value
	}
	if context.working_object != nil {
		*context.changes <- ContextFieldReference{
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
	fragments map[string]parser.FragmentReference
}
