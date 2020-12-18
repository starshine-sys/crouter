package crouter

import (
	"fmt"
	"strings"
)

// TrimPrefixesSpace trims all given prefixes as well as whitespace from the given string
func TrimPrefixesSpace(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		s = strings.TrimPrefix(s, prefix)
	}
	s = strings.TrimSpace(s)
	return s
}

// HasAnyPrefix checks if the string has *any* of the given prefixes
func HasAnyPrefix(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// HasAnySuffix checks if the string has *any* of the given suffixes
func HasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// SprintfAll takes a slice of strings and uses them as input for Sprintf, returning a slice of strings
func SprintfAll(template string, in []string) []string {
	out := make([]string, 0)
	for _, i := range in {
		out = append(out, fmt.Sprintf(template, i))
	}

	return out
}
