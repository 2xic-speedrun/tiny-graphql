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
