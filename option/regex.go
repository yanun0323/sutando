package option

type Regex int

const (
	CaseInsensitive  Regex = 1
	MatchMultiLine   Regex = 1 << 1
	IgnoreWhitespace Regex = 1 << 2
	DotMatchNewLine  Regex = 1 << 3
)
