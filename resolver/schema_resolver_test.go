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

func TestShouldCorrectlyHandleNullObject(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				deployment {
					timestamp
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
	git_object := schema.Add_Object_fetch(
		"deployment",
		func(_variables map[string]interface{}) interface{} {
			// no deployment has happened, so we return nil.
			return nil
		},
	)
	git_object.Add_field(
		"timestamp",
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

	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["deployment"], nil, "Wrong resolved value name")
}

func TestShouldBeAbleToBuildUponObjectContext(t *testing.T) {
	request_schema := `
		query BuildInfo {
			build {
				deployment {
					timestamp
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
	type Deployment struct {
		timestamp string
	}
	git_object := schema.Add_Object_fetch(
		"deployment",
		func(_variables map[string]interface{}) interface{} {
			return Deployment{
				timestamp: "42",
			}
		},
	)
	git_object.Add_field_resolver(
		"timestamp",
		func(context Context) interface{} {
			return context.working_object.value.(Deployment).timestamp
		},
	)

	build_object.Add_object(
		git_object,
	)
	raw_response := Request(request_schema, *schema)

	var response map[string]interface{}
	json.Unmarshal(raw_response, &response)

	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["deployment"].(map[string]interface{})["timestamp"], "42", "Wrong resolved value timestamp")
}

func TestShouldBeAbleToDealWithVariables(t *testing.T) {
	request_schema := `
	query BuildInfo($id: ID) {
		build (id: $id) {
			timestamp
		}
	}
	`
	type Deployment struct {
		timestamp string
		version   string
	}

	schema := &ResolverSchema{
		Resolvers: make(map[string]Resolvers),
	}
	build_object := schema.Add_Object_fetch(
		"build",
		func(m map[string]interface{}) interface{} {
			return &Deployment{
				timestamp: "42",
				version:   m["id"].(string),
			}
		},
	)
	build_object.Add_field_resolver(
		"timestamp",
		func(context Context) interface{} {
			if context.working_object.value.(*Deployment).version == "1" {
				return "1"
			}
			return context.working_object.value.(*Deployment).timestamp
		},
	)

	var response map[string]interface{}
	raw_response := Request_with_variables(request_schema, *schema, map[string]interface{}{
		"$id": "1",
	})
	json.Unmarshal(raw_response, &response)
	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["timestamp"], "1", "Wrong resolved value timestamp")

	raw_response_2 := Request_with_variables(request_schema, *schema, map[string]interface{}{
		"$id": "2",
	})
	json.Unmarshal(raw_response_2, &response)
	assert.Equal(t, response["BuildInfo"].(map[string]interface{})["build"].(map[string]interface{})["timestamp"], "42", "Wrong resolved value timestamp")
}
