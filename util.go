package schema

import (
	"strings"
)

func splitNonEmptyAndTrim(s, sep string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	secs := strings.Split(s, sep)
	for i := range secs {
		secs[i] = strings.TrimSpace(secs[i])
	}
	return secs
}

func hasString(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}
