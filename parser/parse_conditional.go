package parser

func (parser *Parser) ParseConditional() *Conditional {
	if parser.isNextTokenThenSkip("@") {
		if parser.isNextTokenThenSkip("skip") {
			return &Conditional{
				variant:   "skip",
				variables: parser.ParseArguments(),
			}
		} else if parser.isNextTokenThenSkip("include") {
			return &Conditional{
				variant:   "include",
				variables: parser.ParseArguments(),
			}
		}
	}
	return nil
}
