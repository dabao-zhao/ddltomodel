package template

// TypeTpl defines a template for types in model.
const TypeTpl = `
type (
	{{.lowerStartCamelObject}}Model interface{
		{{.method}}
	}

	default{{.upperStartCamelObject}}Model struct {
		conn  *gorm.DB
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
`
