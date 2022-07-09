package resolver

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
		func() interface{} {
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
		func() interface{} {
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
		func() interface{} {
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
		func() interface{} {
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
		func() interface{} {
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
		func() interface{} {
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

func TestShouldFetchFragment(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				...buildInfo
			}
		}

		fragment buildInfo on Build {
			id
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
		func() interface{} {
			return "0x42"
		},
	)
	raw_response := Request(request_schema, *schema)
	var response map[string]map[string]map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"]["build"]["id"], "0x42", "Wrong resolved value name")
}

func TestGoRoutine(t *testing.T) {
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
		func() interface{} {
			time.Sleep(100 * time.Millisecond)
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
		func() interface{} {
			time.Sleep(200 * time.Millisecond)
			return "0x42"
		},
	)

	build_object.Add_object(
		git_object,
	)
	start := time.Now().UnixNano() / int64(time.Millisecond)
	raw_response := Request(request_schema, *schema)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	time := (end - start)

	fmt.Println(time)

	var response map[string]interface{}
	json.Unmarshal(raw_response, &response)

	fmt.Println(response)
	fmt.Println(time)

	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["name"], "tiny-graphql", "Wrong resolved value name")
	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["git"].(map[string]interface{})["hash"], "0x42", "Wrong resolved value name")
	assert.Equal(t, time < 220, true, "error in go routine")
}
