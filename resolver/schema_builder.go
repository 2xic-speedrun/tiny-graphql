package resolver

/*
	We need to build the schema with resolvers.
	- Resolvers resolve objects
	- FieldResolver are in resolvers ?

	Idea
		We create an object interface
			-> Object interface has a registered fields and callback function.
			-> callback functions can be "anything"
			-> based on the callback type we know the schema type, I guess we could also specify it.
				-> Custom Types are a thing in graphql
*/
