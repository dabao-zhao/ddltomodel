package command

import (
	"github.com/dabao-zhao/ddltomodel/gen"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// VarStringSrc describes the source file of sql.
	VarStringSrc string
	// VarStringDir describes the output directory of sql.
	VarStringDir string
	// VarStringDatabase describes the database.
	VarStringDatabase string
	// VarBoolStrict describes whether the strict mode is enabled.
	VarBoolStrict bool
)

type ddlArg struct {
	src, dir string
	database string
	strict   bool
}

// MysqlDDL generates model code from ddl
func MysqlDDL(_ *cobra.Command, _ []string) error {
	src := VarStringSrc
	dir := VarStringDir
	database := VarStringDatabase

	arg := ddlArg{
		src:      src,
		dir:      dir,
		database: database,
		strict:   VarBoolStrict,
	}
	return fromDDL(arg)
}

func fromDDL(arg ddlArg) error {
	src := strings.TrimSpace(arg.src)
	if len(src) == 0 {
		return errors.New("expected path or path globbing patterns, but nothing found")
	}

	files, err := filex.Match(src)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("not found any sql file")
	}

	generator, err := gen.NewDefaultGenerator(arg.dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		err = generator.StartFromDDL(f, arg.strict, arg.database)
		if err != nil {
			return err
		}
	}

	return nil
}
