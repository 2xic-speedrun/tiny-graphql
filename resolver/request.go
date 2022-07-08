package resolver

import (
	"encoding/json"

	"github.com/2xic-speedrun/tiny-graphql/parser"
)

func Request(request_schema string, resolver_schema ResolverSchema) []byte {
	request_fields := parser.Parse(request_schema)
	resolved_fields := resolver_schema.Resolve(request_fields)
	data, error := json.Marshal(resolved_fields)

	if error != nil {
		panic(error)
	}
	return data
}
