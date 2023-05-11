package sqlutil

import (
	"fmt"
	"strings"
)

func ToMaxText(s string, maxLength int) (string, bool) {
	if maxLength == 0 {
		panic(fmt.Sprintf("tpm-batis::ToMaxText - maxLength cannot be 0 (%d)", maxLength))
	}

	truncated := false
	absMaxLength := maxLength
	if absMaxLength < 0 {
		absMaxLength = -maxLength
	}
	if len(s) > absMaxLength {
		truncated = true
		if maxLength > 0 {
			s = s[0:maxLength]
		} else {
			s = s[len(s)+maxLength:]
		}
	}

	return s, truncated
}

func ToSqlString(s string) string {
	var sb strings.Builder
	sb.WriteRune('\'')
	for _, c := range s {
		sb.WriteRune(c)
		if c == '\'' {
			sb.WriteRune(c)
		}
	}
	sb.WriteRune('\'')
	return sb.String()
}
