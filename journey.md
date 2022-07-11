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

## Day 4
Most pieces are now in place, and I added some basic mapping of the resolvers to the request schema.

Tomorrow I will refactor the old array based schema generator where fields, objects, etc were separated and instead do it in a way similar to how it was solved in the schema builder.

Main reason for doing this is to make the code a lot cleaner. You can traverse both schemas in the same way, and it makes it a lot easier to know if you have defined some objects twice (ambiguous).

Goal for tomorrow
- Refactor the parser/schema.go
- Add for recursive objects and resolve them 
```
query test {
  build {
    git {
      hash # return a hash.
    }
  }
}
```
- Add test for key / value order of response
- Add support object arguments. 
- 

## Day 5
Starting on the refactoring of the schema parser now. 

One hour later (ish): reimplemented most of the logic now to use the new map interface. Things are looking good, and working a lot better. 

Still have to move over the fragment and conditional logic, but that should not be too hard. 

So I was able to solve the two first key points from yesterday.

Things to do tomorrow

- Add test for key / value order of response (schema and response should have same order)
- Add support object arguments in resolver
- Update the fragment + conditional logic
- 

^ ended up coding a bit more on this, mainly added support for go routines to allow faster fetching.

## Day 6
Starting to want to find a end date, the point of the project was to have a dive a bit deeper into golang, and I think it's done it's job.

Features I want to have in place before finishing off
- Object fetch -> fetch and object, and have field resolvers on the object. Need to support nil type.
  - Currently all fields on a object could be seen as a "field resolver", but not fully since you can't run arguments on the fields. 
  - we do support objects in objects, but without context, so the next big thing is adding function context.
- Should support nullable fields
  - Don't build upon nullable fields i.e
- Field validation of types, and variables (?)
- Some meaningful error messages ? 
- 

----
Okay, I solved the two biggest points I wanted to have solved mainly nullable objects, and fetching objects. 

I could solve the two next problems, but then it's start to become somewhat of an implementation, and I have to finish it. This was meant to be fun, and not a big serious project.

So I think I will draw the line here for now. Might come back to this project since I still have not played with generics which would be useful here.
