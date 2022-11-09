package gen

import (
	"strings"

	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
)

func genUpdate(table Table) (string, string, error) {
	expressionValues := make([]string, 0)
	pkg := "data."

	for _, field := range table.Fields {
		camel := stringx.SafeString(field.Name.ToCamel())
		if camel == "CreateTime" || camel == "UpdateTime" || camel == "CreateAt" || camel == "UpdateAt" {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			continue
		}

		expressionValues = append(expressionValues, pkg+camel)
	}

	expressionValues = append(
		expressionValues, pkg+table.PrimaryKey.Name.ToCamel(),
	)

	camelTableName := table.Name.ToCamel()
	text, err := filex.LoadTemplate(updateMethodTemplateFile, template.UpdateMethod)
	if err != nil {
		return "", "", err
	}

	methodBuffer, err := output.With("updateMethod").Parse(text).Execute(
		map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"lowerStartCamelObject": stringx.From(camelTableName).Untitle(),
			"upperStartCamelPrimaryKey": stringx.EscapeGolangKeyword(
				stringx.From(table.PrimaryKey.Name.ToCamel()).Title(),
			),
			"lowerStartCamelPrimaryKey": stringx.EscapeGolangKeyword(
				stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			),
			"originalPrimaryKey": wrapWithRawString(
				table.PrimaryKey.Name.Source(),
			),
			"expressionValues": strings.Join(
				expressionValues, ", ",
			),
			"data": table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	// update interface method
	text, err = filex.LoadTemplate(updateInterfaceTemplateFile, template.UpdateMethodInterface)
	if err != nil {
		return "", "", err
	}

	interfaceBuffer, err := output.With("updateInterface").Parse(text).Execute(
		map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"data":                  table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	return methodBuffer.String(), interfaceBuffer.String(), nil
}
