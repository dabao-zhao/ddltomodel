package template

const ImportTpl = `import (
	"context"
	"errors"
	{{if .time}}"time"{{end}}

	"gorm.io/gorm"
)
`
