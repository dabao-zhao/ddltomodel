package template

// ErrorTpl defines an error template
const ErrorTpl = `package {{.pkg}}

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound
`
