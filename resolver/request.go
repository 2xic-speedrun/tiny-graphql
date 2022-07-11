package resolver

import (
	"encoding/json"

	"github.com/2xic-speedrun/tiny-graphql/parser"
)

func Request(request_schema string, resolver_schema ResolverSchema) []byte {
	request_Schema := parser.Parse(request_schema)
	return process_schema(request_Schema, resolver_schema)
}

func Request_with_variables(request_schema string, resolver_schema ResolverSchema, variables map[string]interface{}) []byte {
	request_Schema := parser.Parse(request_schema)
	request_Schema.Inject_variables(variables)
	return process_schema(request_Schema, resolver_schema)
}

func process_schema(request_Schema parser.Schema, resolver_schema ResolverSchema) []byte {
	resolved_fields := resolver_schema.Resolve(request_Schema)
	data, json_error := json.Marshal(resolved_fields)

	if json_error != nil {
		panic(json_error)
	}
	return data
}
