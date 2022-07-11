package resolver

import "fmt"

type Resolvers interface {
	Name() string
	Resolve(name *string, context *Context)
	Type() int
	Child() map[string]Resolvers
}

type Object struct {
	name     string
	resolves map[string]Resolvers
	fetch    func(map[string]interface{}) interface{}
	value    interface{}
}

type Field struct {
	name            string
	resolve         func() interface{}
	resolve_context func(Context) interface{}
	resolves        map[string]Resolvers
}

func (field *Object) Resolve(name *string, context *Context) {
	if name == nil {
		panic("name should not be nil")
	}
	if field.resolves[*name] == nil {
		panic(fmt.Sprintf("invalid field %s", *name))
	}
	// ops this could be another object...
	go field.resolves[*name].Resolve(nil, context)
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

func (object *Object) Add_field_resolver(name string, resolve func(context Context) interface{}) *Field {
	field := &Field{
		name:            name,
		resolve_context: resolve,
		resolves:        nil,
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

func (field *Field) Resolve(name *string, context *Context) {
	var results_value interface{}

	if field.resolve != nil {
		results_value = field.resolve()
	} else if field.resolve_context != nil {
		results_value = field.resolve_context(*context)
	} else {
		panic("Should not happened")
	}
	*context.changes <- ContextFieldReference{
		value:     results_value,
		reference: context.reference,
		name:      *name,
	}
}

type ContextFieldReference struct {
	reference *map[string]interface{}
	value     interface{}
	name      string
	object    bool
}

type Context struct {
	reference      *map[string]interface{}
	working_object *Object
	changes        *chan ContextFieldReference
}
