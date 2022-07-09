package resolver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleRequest(t *testing.T) {
	request_schema := `
		{
			build
		}
	`
	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
	}
	schema.Add_field(
		"build",
		func() string {
			return "0x42"
		},
	)
	raw_response := Request(request_schema, *schema)
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
	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
	}
	schema.Add_field(
		"build",
		func() string {
			return "0x42"
		},
	)
	raw_response := Request(request_schema, *schema)
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
	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
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
	raw_response := Request(request_schema, *schema)
	var response map[string]map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"]["build"]["id"], "0x42", "Wrong resolved value name")
}

func TestQueryObjectInObject(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				git {
					hash
				}
			}
		}
	`
	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
	}
	build_object := schema.Add_Object(
		"build",
	)
	git_object := schema.Add_Object(
		"git",
	)
	git_object.Add_field(
		"hash",
		func() string {
			return "0x42"
		},
	)

	build_object.Add_object(
		git_object,
	)
	raw_response := Request(request_schema, *schema)
	var response map[string]map[string]map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"]["build"]["git"]["hash"], "0x42", "Wrong resolved value name")
}

func TestMultipleObject(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				git {
					hash
				}
			}
			name
		}
	`
	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
	}
	schema.Add_field(
		"name",
		func() string {
			return "tiny-graphql"
		},
	)
	build_object := schema.Add_Object(
		"build",
	)
	git_object := schema.Add_Object(
		"git",
	)
	git_object.Add_field(
		"hash",
		func() string {
			return "0x42"
		},
	)

	build_object.Add_object(
		git_object,
	)

	raw_response := Request(request_schema, *schema)
	var response map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["git"].(map[string]interface{})["hash"], "0x42", "Wrong resolved value name")
	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["name"], "tiny-graphql", "Wrong resolved value name")
}
