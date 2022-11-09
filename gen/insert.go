package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
	"strings"
)

func genInsert(table Table) (string, string, error) {
	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	var count int
	for _, field := range table.Fields {
		camel := stringx.SafeString(field.Name.ToCamel())
		if camel == "CreateTime" || camel == "UpdateTime" || camel == "CreateAt" || camel == "UpdateAt" {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			if table.PrimaryKey.AutoIncrement {
				continue
			}
		}

		count += 1
		expressions = append(expressions, "?")
		expressionValues = append(expressionValues, "data."+camel)
	}

	camel := table.Name.ToCamel()
	text, err := filex.LoadTemplate(insertTemplateMethodFile, template.InsertMethod)
	if err != nil {
		return "", "", err
	}

	methodBuffer, err := output.With("insertMethod").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).Untitle(),
			"expression":            strings.Join(expressions, ", "),
			"expressionValues":      strings.Join(expressionValues, ", "),
			"data":                  table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = filex.LoadTemplate(insertTemplateInterfaceFile, template.InsertMethodInterface)
	if err != nil {
		return "", "", err
	}

	interfaceBuffer, err := output.With("insertInterface").Parse(text).Execute(map[string]interface{}{
		"upperStartCamelObject": camel,
		"data":                  table,
	})
	if err != nil {
		return "", "", err
	}

	return methodBuffer.String(), interfaceBuffer.String(), nil
}
