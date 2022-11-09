package gen

import (
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
)

func genTag(table Table, in string) (string, error) {
	if in == "" {
		return in, nil
	}

	text, err := filex.LoadTemplate(tagTemplateFile, template.TagTpl)
	if err != nil {
		return "", err
	}

	buffer, err := output.With("tag").Parse(text).Execute(map[string]interface{}{
		"field": in,
		"data":  table,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
