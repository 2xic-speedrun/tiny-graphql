package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/2xic-speedrun/tiny-graphql/resolver"

	"github.com/stretchr/testify/assert"
)

func TestSimpleRequest(t *testing.T) {
	request_schema := `
		{
			build
		}
	`
	schema := &resolver.ResolverSchema{}
	schema.Add_field(
		"build",
		func() string {
			return "0x42"
		},
	)
	raw_response := resolver.Request(request_schema, *schema)
	fmt.Println(raw_response)

	var response map[string]interface{}
	json.Unmarshal(raw_response, &response)
	fmt.Println(response)

	assert.Equal(t, response["build"], "0x42", "Wrong resolved value name")
}

type SimpleResponseBuild struct {
	build string
}
