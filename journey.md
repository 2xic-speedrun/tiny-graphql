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

----

Have now more or less written a functional parser. Still some minor things to fix like the types of array and dict.

The main challenge now is how to represent this all. Should it all be a object like it is today ? 

No, I think the correct way is to do something like an interface.

Fragments has to be calculated on the server side, so we cannot compress this.

The variants we have are 
- Object -> can be conditional
- Field -> string -> can be conditional
- Fragments -> can be conditional
- 


```
type BaseElement interface {
  name string
}
```

Okay, maybe we should create a basic resolver object first to get a feeling for how the interface should be ? 

Yeah, let's try that first.

So I will do that tomorrow, I think :)

