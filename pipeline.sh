GO_PATH=~/go
PATH=$PATH:/$GO_PATH/bin

go test github.com/2xic-speedrun/tiny-graphql/parser
go test github.com/2xic-speedrun/tiny-graphql/resolver
go test github.com/2xic-speedrun/tiny-graphql/poc-app
go build -o poc ./poc-app
