package util

import (
	"regexp"
)

var re_leadclose_whtsp = regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
var re_inside_whtsp = regexp.MustCompile(`[\s\p{Zs}]{2,}|[\r\n]+`)

func StripDuplicateWhiteSpaces(s string) string {

	if s != "" {
		s = re_leadclose_whtsp.ReplaceAllString(s, "")
		s = re_inside_whtsp.ReplaceAllString(s, " ")
	}

	return s
}
