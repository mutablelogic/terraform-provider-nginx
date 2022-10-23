package util

import "strings"

const (
	domainSep = "."
)

// Fqn returns a fully-qualified value, which includes a trailing dot
func Fqn(value ...string) string {
	var result []string
	if len(value) == 0 {
		return ""
	}
	for _, v := range value {
		result = append(result, strings.Trim(v, domainSep))
	}
	return strings.Join(result, domainSep) + domainSep
}

// Unfqn remove final domain separator and domain
func Unfqn(value, domain string) string {
	return strings.Trim(strings.Trim(strings.Trim(value, domainSep), strings.Trim(domain, domainSep)), domainSep)
}
