package template

const (
	InsertMethod = `
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}}, error) {
	err := m.conn.Table(m.TableName()).Create(&data).Error
	return data, err
}
`

	// InsertMethodInterface defines an interface method template for insert code in model
	InsertMethodInterface = `Insert(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}},error)`
)
