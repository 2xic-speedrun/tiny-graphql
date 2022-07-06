package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParserName(t *testing.T) {
	schema := Parse(`
	  query GetUserName {
		user(id: 4) {
		  name
		}
	  }	  
	`)

	assert.Equal(t, schema.name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.objects[0].name, "user", "Wrong object name")
	assert.Equal(t, schema.objects[0].fields[0].name, "name", "Wrong field name")
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
	assert.Equal(t, schema.name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.objects[0].name, "user", "Wrong object name")
	assert.Equal(t, schema.objects[0].fields[0].name, "name", "Wrong field name")
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
	assert.Equal(t, schema.name, "GetUserName", "Wrong schema name")
	assert.Equal(t, schema.objects[0].name, "user", "Wrong object name")
	assert.Equal(t, schema.objects[0].objects[0].name, "nameObject", "Wrong object name")
	assert.Equal(t, schema.objects[0].objects[0].fields[0].name, "nameField", "Wrong field name")
}

func TestParserNameOrQueryDoesNotNeedToBeSpecified(t *testing.T) {
	schema := Parse(`
	  {
		user(id: 4) {
		  name
		}
	  }	  
	`)

	assert.Equal(t, schema.name, "root", "Wrong schema name")
	assert.Equal(t, schema.objects[0].name, "user", "Wrong object name")
	assert.Equal(t, schema.objects[0].fields[0].name, "name", "Wrong field name")
}

func TestParserAlias(t *testing.T) {
	schema := Parse(`
	  {
		aliasObject: user(id: 4) {
		  alias: name
		}
	  }	  
	`)

	assert.Equal(t, schema.name, "root", "Wrong schema name")
	assert.Equal(t, nil != schema.objects[0].alias, true, "Wrong alias name")
	assert.Equal(t, *schema.objects[0].alias, "aliasObject", "Wrong alias name")
	assert.Equal(t, schema.objects[0].name, "user", "Wrong object name")
	assert.Equal(t, schema.objects[0].fields[0].name, "name", "Wrong field name")
	assert.Equal(t, nil != schema.objects[0].fields[0].alias, true, "Wrong alias name")
	assert.Equal(t, *schema.objects[0].fields[0].alias, "alias", "Wrong field name")
}

func TestParserInputVariables(t *testing.T) {
	schema := Parse(`
	  query Test($id:int) {
		aliasObject: user(id: $id) {
		  alias: name
		}
	  }	  
	`)

	assert.Equal(t, schema.name, "Test", "Wrong schema name")
	assert.Equal(t, len(schema.variables), 1, "Wrong variable length")
	assert.Equal(t, schema.variables[0].key, "$id", "Wrong variable name")
	assert.Equal(t, schema.variables[0].value, "int", "Wrong variable name")

	assert.Equal(t, schema.objects[0].variables[0].key, "id", "Wrong variable name")
	assert.Equal(t, schema.objects[0].variables[0].value, "$id", "Wrong variable name")
}
