## Graphql speedrun

This was supposed to be a project for me to play some more with go with a codebase that I thought would allow be to visit more aspects of go in a short amount of time.

This is a speedrun after all :).

The plan was to implement as much as possible of the GraphQL spec ( https://spec.graphql.org/October2021/ ) in a week (ish). The point was never to complete everything, but more having it as a common thread to write the go code around.

Results
- Super basic parser was written
    - Will parse basic GraphQL schemas, but does not have full support for fragments.
    - Does parse types, but have not tested with dict / array passing, I think I only made it parse, but not correctly store it in the parsed schema.
- Super basic schema processor was written
  - Parses a schema
    - Uses user defined resolvers to resolve the objects / fields from the schema. 
