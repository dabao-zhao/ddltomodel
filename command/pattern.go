package command

import (
	"log"
	"path/filepath"
	"strings"
)

type pattern map[string]struct{}

func (p pattern) Match(s string) bool {
	for v := range p {
		match, err := filepath.Match(v, s)
		if err != nil {
			log.Printf("%+v", err)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func (p pattern) list() []string {
	var ret []string
	for v := range p {
		ret = append(ret, v)
	}
	return ret
}

func parseTableList(tableValue []string) pattern {
	tablePattern := make(pattern)
	for _, v := range tableValue {
		fields := strings.FieldsFunc(v, func(r rune) bool {
			return r == ','
		})
		for _, f := range fields {
			tablePattern[f] = struct{}{}
		}
	}
	return tablePattern
}
