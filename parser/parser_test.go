package parser

import (
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

	if schema.name != "GetUserName" {
		panic("Wrong schema name")
	}
	if schema.objects[0].name != "user" {
		panic("Wrong with schema parser " + schema.objects[0].name)
	}
	if schema.objects[0].fields[0] != "name" {
		panic("Wrong with schema parser " + schema.objects[0].fields[0])
	}
}
