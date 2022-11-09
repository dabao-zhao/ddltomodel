package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
)

func genImports(table Table, timeImport bool) (string, error) {
	text, err := filex.LoadTemplate(importTemplateFile, template.ImportTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("import").Parse(text).Execute(map[string]interface{}{
		"time": timeImport,
		"data": table,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
