package gen

import (
	"strings"

	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
)

func genFindOne(table Table) (string, string, error) {
	camel := table.Name.ToCamel()
	text, err := filex.LoadTemplate(findOneMethodTemplateFile, template.FindOneMethod)
	if err != nil {
		return "", "", err
	}

	methodBuffer, err := output.With("findOneMethod").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"lowerStartCamelObject":     stringx.From(camel).Untitle(),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source()),
			"lowerStartCamelPrimaryKey": stringx.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.DataType,
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	text, err = filex.LoadTemplate(findOneInterfaceTemplateFile, template.FindOneInterface)
	if err != nil {
		return "", "", err
	}

	interfaceBuffer, err := output.With("findOneInterface").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"lowerStartCamelPrimaryKey": stringx.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle()),
			"dataType":                  table.PrimaryKey.DataType,
			"data":                      table,
		})
	if err != nil {
		return "", "", err
	}

	return methodBuffer.String(), interfaceBuffer.String(), nil
}

func wrapWithRawString(v string) string {
	if v == "`" {
		return v
	}

	if !strings.HasPrefix(v, "`") {
		v = "`" + v
	}

	if !strings.HasSuffix(v, "`") {
		v = v + "`"
	} else if len(v) == 1 {
		v = v + "`"
	}

	return v
}
