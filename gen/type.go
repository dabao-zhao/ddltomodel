package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
)

func genTypes(table Table, methods string) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(table, fields)
	if err != nil {
		return "", err
	}

	text, err := filex.LoadTemplate(typeTemplateFile, template.TypeTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"method":                methods,
			"upperStartCamelObject": table.Name.ToCamel(),
			"lowerStartCamelObject": stringx.From(table.Name.ToCamel()).Untitle(),
			"fields":                fieldsString,
			"data":                  table,
		})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
