## GraphQL implementation

## Day 1.
Just setting things up. Wrote a basic parser that can parse a query schema. 

Tomorrow I will finish the parser, and clean up the code used for the parsing. 

## Day 2
Plan : More or less complete the parser. Missing features are for instance
- Validation of types (string, int, date, dict, Array, null, enum)
- Field alias
- Mutation
- Schema does not need to have query specified, or name.
- Input variables for schema ($variables)


### Todo left from today
Will implement the following (hopefully) tomorrow.
- Fragments
  - Type conditions
  - ... operator
- Object + Field order (currently this is not tracked)
- @include operator
- @skip operator

## Day 3
Plan : Complete the todo from yesterday, need to fix the way some tokens are constructed (i.e inside strings we should ignore terminators).

