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
	assert.Equal(t, schema.objects[0].fields[0], "name", "Wrong field name")
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
	assert.Equal(t, schema.objects[0].fields[0], "name", "Wrong field name")
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
	assert.Equal(t, schema.objects[0].objects[0].fields[0], "nameField", "Wrong field name")
}
