package option

type Regex int

const (
	// CaseInsensitive makes case insensitivity to match upper and lower cases.
	CaseInsensitive Regex = 1

	// MatchMultiLine includes anchors
	// (i.e. ^ for the start, $ for the end),
	// match at the beginning or end of each line for
	// strings with multiline values. Without this option,
	// these anchors match at beginning or end of the string.
	//
	// If the pattern contains no anchors or
	// if the string value has no newline characters
	// (e.g. \n), the m option has no effect.
	MatchMultiLine Regex = 1 << 1

	// IgnoreWhitespace ignores all white space characters
	// in the $regex pattern unless escaped
	// or included in a character class.
	//
	// Additionally, it ignores characters in-between
	// and including an un-escaped hash/pound (#) character
	// and the next new line, so that you may include comments
	// in complicated patterns. This only applies to
	// data characters; white space characters may never appear
	// within special character sequences in a pattern.
	//
	// The x option does not affect the handling
	// of the VT character (i.e. code 11).
	IgnoreWhitespace Regex = 1 << 2

	// Allows the dot character (i.e. .)
	// to match all characters
	// including newline characters.
	DotMatchNewLine Regex = 1 << 3
)
