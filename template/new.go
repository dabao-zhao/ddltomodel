package template

// NewTpl defines the template for creating model instance.
const NewTpl = `
func new{{.upperStartCamelObject}}Model(conn *gorm.DB) *default{{.upperStartCamelObject}}Model {
	return &default{{.upperStartCamelObject}}Model{
		conn:conn,
		table:{{.table}},
	}
}
`
