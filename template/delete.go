package template

const (
	// DeleteMethod defines a delete template
	DeleteMethod = `
func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	err := m.conn.Table(m.TableName()).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).Delete().Error
	return err
}
`

	// DeleteMethodInterface defines a delete template for interface method
	DeleteMethodInterface = `Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error`
)
