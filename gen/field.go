package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/parser"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
	"strings"
)

func genFields(table Table, fields []*parser.Field) (string, error) {
	var list []string

	for _, field := range fields {
		result, err := genField(table, field)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(table Table, field *parser.Field) (string, error) {
	tag, err := genTag(table, field.NameOriginal)
	if err != nil {
		return "", err
	}

	text, err := filex.LoadTemplate(fieldTemplateFile, template.FieldTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"name":       stringx.SafeString(field.Name.ToCamel()),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
			"data":       table,
		})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
