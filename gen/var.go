package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
)

func genVars(table Table) (string, error) {
	text, err := filex.LoadTemplate(varTemplateFile, template.VarTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("var").Parse(text).
		GoFmt(true).Execute(map[string]interface{}{})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
