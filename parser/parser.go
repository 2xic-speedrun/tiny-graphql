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

func (parser *Parser) isNextToken(expected string) bool {
	value := parser.Peek(0)
	if value != nil && *value == expected {
		parser.index += 1
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

type Parser struct {
	Tokens []string
	index  int
}
