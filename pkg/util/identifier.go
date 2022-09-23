package util

import "regexp"

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ReIdentifier = `[a-zA-Z][a-zA-Z0-9_\-]+`
)

var (
	reValidName = regexp.MustCompile(`^` + ReIdentifier + `$`)
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func IsIdentifier(s string) bool {
	return reValidName.MatchString(s)
}
