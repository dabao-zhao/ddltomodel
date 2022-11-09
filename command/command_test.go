package command

import (
	"github.com/dabao-zhao/ddltomodel/gen"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	_ "embed"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	//go:embed fortest/user.sql
	sql string
)

func TestFromDDl(t *testing.T) {
	err := gen.Clean()
	assert.Nil(t, err)

	// 文件不存在
	err = fromDDL(ddlArg{
		src:      "./user.sql",
		dir:      filex.MustTempDir(),
		database: "for-test",
		strict:   false,
	})
	assert.Equal(t, errors.New("not found any sql file"), err)

	// 文件夹不存在
	unknownDir := filepath.Join(filex.MustTempDir(), "test", "user.sql")
	err = fromDDL(ddlArg{
		src:      unknownDir,
		dir:      filex.MustTempDir(),
		database: "for_test",
	})
	assert.True(t, func() bool {
		switch err.(type) {
		case *os.PathError:
			return true
		default:
			return false
		}
	}())

	// 空文件夹
	err = fromDDL(ddlArg{
		dir:      filex.MustTempDir(),
		database: "for_test",
	})
	if err != nil {
		assert.Equal(t, "expected path or path globbing patterns, but nothing found", err.Error())
	}

	tempDir := filepath.Join(filex.MustTempDir(), "test")
	err = filex.MkdirIfNotExist(tempDir)
	if err != nil {
		return
	}

	user1Sql := filepath.Join(tempDir, "user1.sql")
	user2Sql := filepath.Join(tempDir, "user2.sql")

	err = os.WriteFile(user1Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	err = os.WriteFile(user2Sql, []byte(sql), os.ModePerm)
	if err != nil {
		return
	}

	_, err = os.Stat(user1Sql)
	assert.Nil(t, err)

	_, err = os.Stat(user2Sql)
	assert.Nil(t, err)

	filename := filepath.Join(tempDir, "user_model.go")
	fromDDL := func(db string) {
		err = fromDDL(ddlArg{
			src:      filepath.Join(tempDir, "user*.sql"),
			dir:      tempDir,
			database: db,
		})
		assert.Nil(t, err)

		_, err = os.Stat(filename)
		assert.Nil(t, err)
	}

	fromDDL("for_test")
	_ = os.Remove(filename)
	fromDDL("for-test")
	_ = os.Remove(filename)
	fromDDL("1fortest")
}
