module tiny-graphql

go 1.18

require (
	github.com/2xic-speedrun/tiny-graphql/parser v1.0.0
	github.com/2xic-speedrun/tiny-graphql/resolver v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/2xic-speedrun/tiny-graphql/parser => ./parser
replace github.com/2xic-speedrun/tiny-graphql/resolver => ./resolver
