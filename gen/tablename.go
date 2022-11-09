package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
)

func genTableName(table Table) (string, error) {
	text, err := filex.LoadTemplate(tableNameTemplateFile, template.TableNameTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"tableName":             table.Name.Source(),
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
