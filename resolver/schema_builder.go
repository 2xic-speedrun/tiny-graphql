package resolver

import (
	"fmt"

	"github.com/2xic-speedrun/tiny-graphql/parser"
)

/*
	We need to build the schema with resolvers.
	- Resolvers resolve objects
	- FieldResolver are in resolvers ?

	Idea
		We create an object interface
			-> Object interface has a registered fields and callback function.
			-> callback functions can be "anything"
			-> based on the callback type we know the schema type, I guess we could also specify it.
				-> Custom Types are a thing in graphql
*/

func (schema *ResolverSchema) Add_field(name string, Resolve func() string) {
	field := &Field{
		name:  name,
		value: Resolve(),
	}
	schema.resolvers = append(schema.resolvers, field)
}

func (schema *ResolverSchema) Resolve(object parser.Schema) map[string]interface{} {
	data := map[string]interface{}{}
	if 0 < len(object.Name) {
		//	data[object.Name] = map[string]interface{}{}
		//		workingRef = data[object.Name]
	}

	fmt.Println("raw ", object.Objects)
	for _, val := range object.Fields {
		field := schema.resolve_field(val.Name)
		fmt.Println("hello :)")
		if field == nil {
			panic("invalid field")
		}
		fmt.Println()
		fmt.Println(field)
		value := (*field).Resolve()
		data[val.Name] = value
	}

	return data
}

func (schema *ResolverSchema) resolve_field(name string) *Resolvers {
	for _, val := range schema.resolvers {
		if val.Name() == name {
			return &val
		}
	}
	return nil
}

type ResolverSchema struct {
	resolvers []Resolvers
}

type Resolvers interface {
	Name() string
	Resolve() string
}

type Object struct {
	name    string
	resolve func(arguments string)
}

type Field struct {
	name  string
	value string
}

func (field *Field) Resolve() string {
	return field.value
}

func (field *Field) Name() string {
	return field.name
}
