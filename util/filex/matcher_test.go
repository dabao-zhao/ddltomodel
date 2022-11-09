package filex

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	dir, err := filepath.Abs("./fortest/")
	assert.Nil(t, err)

	files, err := Match("./fortest/*.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "studeat.sql"), filepath.Join(dir, "student.sql"), filepath.Join(dir, "xx.sql")}, files)

	files, err = Match("./fortest/??.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "xx.sql")}, files)

	files, err = Match("./fortest/*.sq*")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "studeat.sql"), filepath.Join(dir, "student.sql"), filepath.Join(dir, "xx.sql"), filepath.Join(dir, "xx.sqlx")}, files)

	files, err = Match("./fortest/student.sql")
	assert.Nil(t, err)
	assert.Equal(t, []string{filepath.Join(dir, "student.sql")}, files)
}
