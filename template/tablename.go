package template

// TableNameTpl defines a template that generate the tableName method.
const TableNameTpl = `
func (m *default{{.upperStartCamelObject}}Model) TableName() string {
	return m.table
}
`
