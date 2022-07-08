package main

import (
	"encoding/json"
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
	schema := &resolver.ResolverSchema{
		Resolvers: make(map[string]resolver.Resolvers),
	}
	schema.Add_field(
		"build",
		func() string {
			return "0x42"
		},
	)
	raw_response := resolver.Request(request_schema, *schema)
	var response map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)
	assert.Equal(t, response["root"]["build"], "0x42", "Wrong resolved value name")
}

func TestQuerySimpleRequest(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build
		}
	`
	schema := &resolver.ResolverSchema{
		Resolvers: make(map[string]resolver.Resolvers),
	}
	schema.Add_field(
		"build",
		func() string {
			return "0x42"
		},
	)
	raw_response := resolver.Request(request_schema, *schema)
	var response map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"]["build"], "0x42", "Wrong resolved value name")
}

func TestQuerySimpleObjectRequest(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				id
			}
		}
	`
	schema := &resolver.ResolverSchema{
		Resolvers: make(map[string]resolver.Resolvers),
	}
	object := schema.Add_Object(
		"build",
	)
	object.Add_field(
		"id",
		func() string {
			return "0x42"
		},
	)
	raw_response := resolver.Request(request_schema, *schema)
	var response map[string]map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"]["build"]["id"], "0x42", "Wrong resolved value name")
}
