package gen

import (
	"fmt"

	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
)

func genNew(table Table) (string, error) {
	text, err := filex.LoadTemplate(modelNewTemplateFile, template.NewTpl)
	if err != nil {
		return "", err
	}

	t := fmt.Sprintf(`"%s"`, table.Name.Source())

	buffer, err := output.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"table":                 t,
			"upperStartCamelObject": table.Name.ToCamel(),
			"data":                  table,
		})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
