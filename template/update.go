package template

const (
	// UpdateMethod defines a template for generating update codes
	UpdateMethod = `
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}}, error) {
	err := m.conn.Table(m.TableName()).Where("{{.lowerStartCamelPrimaryKey}} = ?", data.{{.upperStartCamelPrimaryKey}}).Updates(&data).Error
	return data, err
}
`

	// UpdateMethodInterface defines an interface method template for generating update codes
	UpdateMethodInterface = `Update(ctx context.Context, data *{{.upperStartCamelObject}}) (*{{.upperStartCamelObject}}, error)`
)
