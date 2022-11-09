package template

// ErrorTpl defines an error template
const ErrorTpl = `package {{.pkg}}

import "gorm.io/gorm"

var ErrNotFound = gorm.ErrRecordNotFound
`
