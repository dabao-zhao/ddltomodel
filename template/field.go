package template

// FieldTpl defines a filed template for types
const FieldTpl = `{{.name}} {{.type}} {{.tag}} {{if .hasComment}}// {{.comment}}{{end}}`
