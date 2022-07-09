package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSingleField(t *testing.T) {
	schema := Parse(`
		{
			build
		}
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Fields["build"].(Fields).Type(), 1, "Wrong field type ")
	assert.Equal(t, schema.Fields["build"].(Fields).Alias() == nil, true, "Wrong field alias")
}

func TestParseSingleFieldAlias(t *testing.T) {
	schema := Parse(`
		{
			alias:build
		}
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Fields["build"].(Fields).Type(), 1, "Wrong field type ")
	assert.Equal(t, *schema.Fields["build"].(Fields).Alias(), "alias", "Wrong field alias ")
}

func TestParserName(t *testing.T) {
	schema := Parse(`
	  query GetUserName {
		user(id: 4) {
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
}

func TestParserInputVariables(t *testing.T) {
	schema := Parse(`
	  query Test($id:int) {
		aliasObject: user(id: $id) {
		  alias: name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "Test", "Wrong schema name")
	assert.Equal(t, len(schema.Variables), 1, "Wrong variable length")
	assert.Equal(t, schema.Variables[0].key, "$id", "Wrong variable name")
	assert.Equal(t, schema.Variables[0].value, "int", "Wrong variable name")
	assert.Equal(t, *schema.Fields["user"].(*Object).Alias(), "aliasObject", "Wrong object alias")
	assert.Equal(t, *&schema.Fields["user"].(*Object).Variables[0].key, "id", "Wrong variable key")
	assert.Equal(t, *&schema.Fields["user"].(*Object).Variables[0].value, "$id", "Wrong variable value")
}

func TestParserComment(t *testing.T) {
	schema := Parse(`
	  # This is a comment
	  query GetUserName {
		user(id: 4) {
			# This is a comment
			name
		}
	  }
	  # This is a comment
	`)
	assert.Equal(t, schema.Name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong field type")
	assert.Equal(t, schema.Fields["user"].(Fields).Name(), "user", "Wrong object name")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong field type")
}

func TestNestedObject(t *testing.T) {
	schema := Parse(`
	  # This is a comment
	  query GetUserName {
		user(id: 4) {
			# This is a comment
			nameObject {
				nameField
			}
		}
	  }
	  # This is a comment
	`)
	assert.Equal(t, schema.Name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(Fields).Name(), "user", "Wrong object name")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["nameObject"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["nameObject"].(*Object).Fields["nameField"].(Fields).Type(), 1, "Wrong object type")
}

func TestParserNameOrQueryDoesNotNeedToBeSpecified(t *testing.T) {
	schema := Parse(`
	  {
		user(id: 4) {
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
}

func TestParserAlias(t *testing.T) {
	schema := Parse(`
	  {
		aliasObject: user(id: 4) {
		  alias: name
		}
	  }
	`)

	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, *(schema.Fields["user"].(*Object)).Alias(), "aliasObject", "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
	assert.Equal(t, *schema.Fields["user"].(*Object).Fields["name"].(Fields).Alias(), "alias", "Wrong object name")
}

func TestParserShouldHandleArray(t *testing.T) {
	schema := Parse(`
	  {
		user(ids: [1,2,3]){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, len(schema.Fields["user"].(*Object).Variables), 1, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
}

func TestParserShouldHandleDict(t *testing.T) {
	schema := Parse(`
	  {
		user(input: {"id": 4}){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, len(schema.Fields["user"].(*Object).Variables), 1, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
}

func TestParserShouldHandleString(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark"){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Fields["user"].(Fields).Type(), 2, "Wrong object type")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Type(), 2, "Wrong object type")
	assert.Equal(t, len(schema.Fields["user"].(*Object).Variables), 1, "Wrong object type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fields["name"].(Fields).Type(), 1, "Wrong object name")
}

func TestParserShouldHandleFragments(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark"){
		  ...SimpleUser
		}
	  }

	  fragment SimpleUser on User {
		name
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Fragment_reference.name, "SimpleUser", "Wrong object type")
	assert.Equal(t, schema.Fragments["SimpleUser"].Fields["name"].(Fields).Type(), 1)
}

func TestParserOnFragment(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark") {
		  ... on Admin {
			...SimpleUser
		  }
		}
	  }

	  fragment SimpleUser on User {
		name
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, (schema.Fields["user"].(*Object)).Fragment_reference.name, "Admin", "Wrong fragment name type")
	assert.Equal(t, schema.Fields["user"].(*Object).Fragment_reference.Fields["SimpleUser"].(*FragmentReference).name, "SimpleUser", "Wrong fragment name")
}

func TestParserIncludeOperator(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark") @include(if: true){
		  ...SimpleUser
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Fields["user"].(*Object).Conditional.variant, "include", "Wrong conditional name")
}

func TestParserSkipOperator(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark") @skip(if: true){
		  ...SimpleUser
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Fields["user"].(*Object).Conditional.variant, "skip", "Wrong conditional name")
}
