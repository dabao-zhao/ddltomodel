package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
)

func genDelete(table Table) (string, string, error) {
	camel := table.Name.ToCamel()
	text, err := filex.LoadTemplate(deleteMethodTemplateFile, template.DeleteMethod)
	if err != nil {
		return "", "", err
	}

	methodBuffer, err := output.With("delete").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"lowerStartCamelPrimaryKey": stringx.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.DataType,
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source()),
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = filex.LoadTemplate(deleteInterfaceTemplateFile, template.DeleteMethodInterface)
	if err != nil {
		return "", "", err
	}

	interfaceBuffer, err := output.With("deleteMethod").
		Parse(text).
		Execute(map[string]interface{}{
			"lowerStartCamelPrimaryKey": stringx.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.DataType,
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	return methodBuffer.String(), interfaceBuffer.String(), nil
}
