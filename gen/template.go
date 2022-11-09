package gen

import "github.com/dabao-zhao/ddltomodel/util/filex"

const (
	importTemplateFile           = "import.tpl"
	varTemplateFile              = "var.tpl"
	insertTemplateMethodFile     = "insert_method.tpl"
	insertTemplateInterfaceFile  = "insert_interface.tpl"
	findOneMethodTemplateFile    = "find_one_method.tpl"
	findOneInterfaceTemplateFile = "find_one_interface.tpl"
	updateMethodTemplateFile     = "update_method.tpl"
	updateInterfaceTemplateFile  = "update_interface.tpl"
	deleteMethodTemplateFile     = "delete_method.tpl"
	deleteInterfaceTemplateFile  = "delete_interface.tpl"
	tagTemplateFile              = "tag.tpl"
	fieldTemplateFile            = "field.tpl"
	typeTemplateFile             = "type.tpl"
	modelNewTemplateFile         = "model_new.tpl"
	tableNameTemplateFile        = "table_name.tpl"
	modelGenTemplateFile         = "model_gen.tpl"
	modelCustomTemplateFile      = "model_custom.tpl"
	errTemplateFile              = "err.tpl"
)

// Clean deletes all template files
func Clean() error {
	return filex.Clean()
}
