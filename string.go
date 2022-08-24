package postgresql

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

// RemoveDiacritics ...
func RemoveDiacritics(s string) string {
	if s != "" {
		s = strings.ToLower(s)
		s = replaceStringWithRegex(s, `đ`, "d")
		t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		result, _, _ := transform.String(t, s)
		result = replaceStringWithRegex(result, `[^a-zA-Z0-9\s]`, "")

		return result
	}
	return ""
}

// replaceStringWithRegex ...
func replaceStringWithRegex(src string, regex string, replaceText string) string {
	reg := regexp.MustCompile(regex)
	return reg.ReplaceAllString(src, replaceText)
}

// TransformKeywordToSearchString ...
func TransformKeywordToSearchString(keyword string) string {
	s := strings.Trim(keyword, " ")
	s = RemoveDiacritics(s)
	s = strings.ReplaceAll(s, " ", "&")
	return s + ":*" // For prefix search
}
