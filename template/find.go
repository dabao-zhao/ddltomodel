package template

const (
	FindOneMethod = `
func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error) {
	var ret *{{.upperStartCamelObject}}
	err := m.conn.WithContext(ctx).Table(m.TableName()).Where("{{.lowerStartCamelPrimaryKey}} = ?", {{.lowerStartCamelPrimaryKey}}).First(&ret).Error
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}
`
	// FindOneInterface defines find row method.
	FindOneInterface = `FindOne(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)`
)
