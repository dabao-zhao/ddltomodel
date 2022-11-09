package command

import (
	"errors"
	"github.com/dabao-zhao/ddltomodel/model"
	"log"
	"strings"

	"github.com/dabao-zhao/ddltomodel/gen"
	"github.com/dabao-zhao/ddltomodel/util/filex"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	// VarStringURL describes the dsn of the sql.
	VarStringURL string
	// VarStringSliceTable describes tables.
	VarStringSliceTable []string
)

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

// MySqlDataSource generates model code from ddl
func MySqlDataSource(_ *cobra.Command, _ []string) error {
	url := strings.TrimSpace(VarStringURL)
	dir := strings.TrimSpace(VarStringDir)

	tableValue := VarStringSliceTable
	patterns := parseTableList(tableValue)

	arg := dataSourceArg{
		url:      url,
		dir:      dir,
		tablePat: patterns,
		strict:   VarBoolStrict,
	}
	return fromMysqlDataSource(arg)
}

type ddlArg struct {
	src, dir string
	database string
	strict   bool
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

type dataSourceArg struct {
	url, dir string
	tablePat pattern
	strict   bool
}

func fromMysqlDataSource(arg dataSourceArg) error {
	if len(arg.url) == 0 {
		log.Printf("%v", "expected data source of mysql, but nothing found")
		return nil
	}
	if len(arg.tablePat) == 0 {
		log.Printf("%v", "expected table or table globbing patterns, but nothing found")
		return nil
	}

	dsn, err := mysql.ParseDSN(arg.url)
	if err != nil {
		return err
	}

	databaseSource := strings.TrimSuffix(arg.url, "/"+dsn.DBName) + "/information_schema"
	db, err := gorm.Open(gormMysql.Open(databaseSource), &gorm.Config{})
	if err != nil {
		return err
	}
	im := model.NewInformationSchemaModel(db)

	tables, err := im.GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	matchTables := make(map[string]*model.Table)
	for _, item := range tables {
		if !arg.tablePat.Match(item) {
			continue
		}

		columnData, err := im.FindColumns(dsn.DBName, item)
		if err != nil {
			return err
		}

		table, err := columnData.Convert()
		if err != nil {
			return err
		}

		matchTables[item] = table
	}
	if len(matchTables) == 0 {
		return errors.New("no tables matched")
	}

	generator, err := gen.NewDefaultGenerator(arg.dir)
	if err != nil {
		return err
	}

	return generator.StartFromInformationSchema(matchTables, arg.strict)
}
