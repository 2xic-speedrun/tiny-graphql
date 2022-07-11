package parser

func (parser *Parser) ParseArray() *string {
	results := ""
	parser.ParseScope("[", "]", func() {
		results += *parser.Read()
	},
		parser.DictAndArrayTerminatorFunction)
	if len(results) == 0 {
		return nil
	}
	return &results
}

func (parser *Parser) ParseDict() *string {
	results := ""
	parser.ParseScope("{", "}", func() {
		results += *parser.Read()
	},
		parser.DictAndArrayTerminatorFunction)
	if len(results) == 0 {
		return nil
	}
	return &results
}

const (
	NO_TYPE = -1
	ID      = 0
	INT     = 1
	STRING  = 2
	FLOAT   = 3
)

func (parser *Parser) ParseValidType() int {
	if parser.isNextTokenThenSkip("id") {
		return ID
	} else if parser.isNextTokenThenSkip("int") {
		return INT
	} else if parser.isNextTokenThenSkip("string") {
		return STRING
	} else if parser.isNextTokenThenSkip("float") {
		return FLOAT
	} else {
		return NO_TYPE
	}
}
