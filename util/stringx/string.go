package stringx

import (
	"bytes"
	"log"
	"strings"
	"unicode"
)

// SafeString converts the input string into a safe naming style in golang
func SafeString(in string) string {
	if len(in) == 0 {
		return in
	}

	data := strings.Map(func(r rune) rune {
		if isSafeRune(r) {
			return r
		}
		return '_'
	}, in)

	headRune := rune(data[0])
	if isNumber(headRune) {
		return "_" + data
	}
	return data
}

func isSafeRune(r rune) bool {
	return isLetter(r) || isNumber(r) || r == '_'
}

func isLetter(r rune) bool {
	return 'A' <= r && r <= 'z'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

var goKeyword = map[string]string{
	"var":         "variable",
	"const":       "constant",
	"package":     "pkg",
	"func":        "function",
	"return":      "rtn",
	"defer":       "dfr",
	"go":          "goo",
	"select":      "slt",
	"struct":      "structure",
	"interface":   "itf",
	"chan":        "channel",
	"type":        "tp",
	"map":         "mp",
	"range":       "rg",
	"break":       "brk",
	"case":        "caz",
	"continue":    "ctn",
	"for":         "fr",
	"fallthrough": "fth",
	"else":        "es",
	"if":          "ef",
	"switch":      "swt",
	"goto":        "gt",
	"default":     "dft",
}

// EscapeGolangKeyword escapes the golang keywords.
func EscapeGolangKeyword(s string) string {
	if !isGolangKeyword(s) {
		return s
	}

	r := goKeyword[s]
	log.Printf("[EscapeGolangKeyword]: go keyword is forbidden %q, converted into %q", s, r)
	return r
}

func isGolangKeyword(s string) bool {
	_, ok := goKeyword[s]
	return ok
}

// String  provides for converting the source text into other spell case,like lower,snake,camel
type String struct {
	source string
}

// Source returns the source string value
func (s String) Source() string {
	return s.source
}

// From converts the input text to String and returns it
func From(data string) String {
	return String{source: data}
}

// ToCamel converts the input text into camel case
func (s String) ToCamel() string {
	list := s.splitBy(func(r rune) bool {
		return r == '_'
	}, true)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Title())
	}
	return strings.Join(target, "")
}

// Title calls the strings.Title
func (s String) Title() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	return strings.Title(s.source)
}

// it will not ignore spaces
func (s String) splitBy(fn func(r rune) bool, remove bool) []string {
	if s.IsEmptyOrSpace() {
		return nil
	}
	var list []string
	buffer := new(bytes.Buffer)
	for _, r := range s.source {
		if fn(r) {
			if buffer.Len() != 0 {
				list = append(list, buffer.String())
				buffer.Reset()
			}
			if !remove {
				buffer.WriteRune(r)
			}
			continue
		}
		buffer.WriteRune(r)
	}
	if buffer.Len() != 0 {
		list = append(list, buffer.String())
	}
	return list
}

// IsEmptyOrSpace returns true if the length of the string value is 0 after call strings.TrimSpace, or else returns false
func (s String) IsEmptyOrSpace() bool {
	if len(s.source) == 0 {
		return true
	}
	if strings.TrimSpace(s.source) == "" {
		return true
	}
	return false
}

// Untitle return the original string if rune is not letter at index 0
func (s String) Untitle() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	r := rune(s.source[0])
	if !unicode.IsUpper(r) && !unicode.IsLower(r) {
		return s.source
	}
	return string(unicode.ToLower(r)) + s.source[1:]
}
