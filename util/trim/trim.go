package trim

import "strings"

// StringSlice returns a copy slice without empty string item
func StringSlice(list []string) []string {
	var out []string
	for _, item := range list {
		if len(item) == 0 {
			continue
		}
		out = append(out, item)
	}
	return out
}

// NewLine trims \r and \n chars.
func NewLine(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
