package parser

func (parser *Parser) Peek(length int) *string {
	if (parser.index + length) < len(parser.Tokens) {
		return &parser.Tokens[parser.index+length]
	}
	return nil
}

func (parser *Parser) Read() *string {
	results := parser.Peek(0)
	parser.index++
	return results
}

func (parser *Parser) isNextTokenThenSkip(expected string) bool {
	if parser.isPeekToken(expected, 0) {
		parser.index += 1
		return true
	}
	return false
}

func (parser *Parser) isPeekToken(expected string, peek int) bool {
	value := parser.Peek(peek)
	if value != nil && *value == expected {
		return true
	}
	return false
}

func (parser *Parser) isNextTokenSequence(sequence []string) bool {
	for index, item := range sequence {
		reference := parser.Peek(index)
		if reference == nil || *reference != item {
			return false
		}
	}
	parser.index += len(sequence)
	return true
}

func (parser *Parser) ParseScope(init string, terminator string, callback func(), terminatorFunction func(terminator string) bool) bool {
	peekArguments := parser.Peek(0)
	if peekArguments != nil && *peekArguments == init {
		parser.index += 1
		for true {
			if terminatorFunction(terminator) {
				break
			}
			callback()
		}
		parser.index += 1
		return true
	}
	return false
}

type Parser struct {
	Tokens []string
	index  int
}
