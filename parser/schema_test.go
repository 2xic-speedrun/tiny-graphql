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
	assert.Equal(t, schema.fields["build"].(Fields).Type(), 1, "Wrong field type ")
	assert.Equal(t, schema.fields["build"].(Fields).Alias() == nil, true, "Wrong field alias")
}

func TestParseSingleFieldAlias(t *testing.T) {
	schema := Parse(`
		{
			alias:build
		}
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.fields["build"].(Fields).Type(), 1, "Wrong field type ")
	assert.Equal(t, *schema.fields["build"].(Fields).Alias(), "alias", "Wrong field alias ")
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
	assert.Equal(t, schema.Name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.fields["user"].(Fields).Type(), 2, "Wrong object name")
	assert.Equal(t, (schema.fields["user"].(*Object)).Type(), 2, "Wrong object name")
	assert.Equal(t, schema.fields["user"].(*Object).fields["name"].(Fields).Type(), 1, "Wrong object name")
}

/*
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
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
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
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].objects[0].Name, "nameObject", "Wrong object name")
	assert.Equal(t, schema.Objects[0].objects[0].Fields[0].Name, "nameField", "Wrong field name")
}

func TestParserNameOrQueryDoesNotNeedToBeSpecified(t *testing.T) {
	schema := Parse(`
	  {
		user(id: 4) {
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
}

func TestParserAlias(t *testing.T) {
	schema := Parse(`
	  {
		aliasObject: user(id: 4) {
		  alias: name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, nil != schema.Objects[0].Alias, true, "Wrong alias name")
	assert.Equal(t, *schema.Objects[0].Alias, "aliasObject", "Wrong alias name")
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
	assert.Equal(t, nil != schema.Objects[0].Fields[0].alias, true, "Wrong alias name")
	assert.Equal(t, *schema.Objects[0].Fields[0].alias, "alias", "Wrong field name")
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

	assert.Equal(t, schema.Objects[0].variables[0].key, "id", "Wrong variable name")
	assert.Equal(t, schema.Objects[0].variables[0].value, "$id", "Wrong variable name")
}

func TestParserShouldHandleArray(t *testing.T) {
	schema := Parse(`
	  {
		user(ids: [1,2,3]){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
	assert.Equal(t, len(schema.Objects[0].variables), 1, "Wrong variable length")
}

func TestParserShouldHandleDict(t *testing.T) {
	schema := Parse(`
	  {
		user(input: {"id": 4}){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
	assert.Equal(t, len(schema.Objects[0].variables), 1, "Wrong variable length")
}

func TestParserShouldHandleString(t *testing.T) {
	schema := Parse(`
	  {
		user(name: "mark"){
		  name
		}
	  }
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, schema.Objects[0].Name, "user", "Wrong object name")
	assert.Equal(t, schema.Objects[0].Fields[0].Name, "name", "Wrong field name")
	assert.Equal(t, len(schema.Objects[0].variables), 1, "Wrong variable length")
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
	assert.Equal(t, len(schema.Objects[0].fragments), 1, "Wrong object name")
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
	assert.Equal(t, schema.Objects[0].conditional.variant, "include", "Wrong conditional name")
	assert.Equal(t, len(schema.Objects[0].fragments), 1, "Wrong fragment name")
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
	assert.Equal(t, schema.Objects[0].conditional.variant, "skip", "Wrong conditional name")
	assert.Equal(t, len(schema.Objects[0].fragments), 1, "Wrong fragment name")
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
	`)

	assert.Equal(t, schema.Name, "root", "Wrong schema name")
	assert.Equal(t, len(schema.Objects[0].fragments[0].child.fragments), 1, "Wrong schema name")
	assert.Equal(t, schema.Objects[0].fragments[0].object, "Admin", "Wrong schema name")
	assert.Equal(t, schema.Objects[0].fragments[0].child.fragments[0].name, "SimpleUser", "Wrong schema name")
}
*/
