package gen

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dabao-zhao/ddltomodel/model"
	"github.com/dabao-zhao/ddltomodel/output"
	"github.com/dabao-zhao/ddltomodel/parser"
	"github.com/dabao-zhao/ddltomodel/template"
	"github.com/dabao-zhao/ddltomodel/util/filex"
	"github.com/dabao-zhao/ddltomodel/util/stringx"
	"github.com/dabao-zhao/ddltomodel/util/trim"
)

const pwd = "."

type (
	defaultGenerator struct {
		dir string
		pkg string
	}

	code struct {
		importsCode string
		varsCode    string
		typesCode   string
		newCode     string
		insertCode  string
		findCode    []string
		updateCode  string
		deleteCode  string
		cacheExtra  string
		tableName   string
	}

	codeTuple struct {
		modelCode       string
		modelCustomCode string
	}
)

// NewDefaultGenerator creates an instance for defaultGenerator
func NewDefaultGenerator(dir string) (*defaultGenerator, error) {
	if dir == "" {
		dir = pwd
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dir = dirAbs
	pkg := stringx.SafeString(filepath.Base(dirAbs))
	err = filex.MkdirIfNotExist(dir)
	if err != nil {
		return nil, err
	}

	generator := &defaultGenerator{
		dir: dir,
		pkg: pkg,
	}

	return generator, nil
}

func (g *defaultGenerator) StartFromDDL(filename string, strict bool, database string) error {
	modelList, err := g.genFromDDL(filename, strict, database)
	if err != nil {
		return err
	}

	return g.createFile(modelList)
}

func (g *defaultGenerator) genFromDDL(filename string, strict bool, database string) (map[string]*codeTuple, error) {
	m := make(map[string]*codeTuple)
	tables, err := parser.Parse(filename, database, strict)
	if err != nil {
		return nil, err
	}

	for _, e := range tables {
		code, err := g.genModel(*e)
		if err != nil {
			return nil, err
		}
		customCode, err := g.genModelCustom(*e)
		if err != nil {
			return nil, err
		}

		m[e.Name.Source()] = &codeTuple{
			modelCode:       code,
			modelCustomCode: customCode,
		}
	}

	return m, nil
}

func (g *defaultGenerator) StartFromInformationSchema(tables map[string]*model.Table, strict bool) error {
	m := make(map[string]*codeTuple)
	for _, each := range tables {
		table, err := parser.ConvertDataType(each, strict)
		if err != nil {
			return err
		}

		code, err := g.genModel(*table)
		if err != nil {
			return err
		}
		customCode, err := g.genModelCustom(*table)
		if err != nil {
			return err
		}

		m[table.Name.Source()] = &codeTuple{
			modelCode:       code,
			modelCustomCode: customCode,
		}
	}

	return g.createFile(m)
}

func (g *defaultGenerator) createFile(modelList map[string]*codeTuple) error {
	dirAbs, err := filepath.Abs(g.dir)
	if err != nil {
		return err
	}

	g.dir = dirAbs
	g.pkg = stringx.SafeString(filepath.Base(dirAbs))
	err = filex.MkdirIfNotExist(dirAbs)
	if err != nil {
		return err
	}

	for tableName, codes := range modelList {
		tn := stringx.From(tableName)
		modelFilename := fmt.Sprintf("%s_model", tn.Source())

		name := stringx.SafeString(modelFilename) + "_gen.go"
		filename := filepath.Join(dirAbs, name)
		err = os.WriteFile(filename, []byte(codes.modelCode), os.ModePerm)
		if err != nil {
			return err
		}

		name = stringx.SafeString(modelFilename) + ".go"
		filename = filepath.Join(dirAbs, name)
		if filex.FileExists(filename) {
			log.Printf("%s already exists, ignored.", name)
			continue
		}
		err = os.WriteFile(filename, []byte(codes.modelCustomCode), os.ModePerm)
		if err != nil {
			return err
		}
	}

	// generate error file
	varFilename := "vars"

	filename := filepath.Join(dirAbs, varFilename+".go")
	text, err := filex.LoadTemplate(errTemplateFile, template.ErrorTpl)
	if err != nil {
		return err
	}

	err = output.With("vars").Parse(text).SaveTo(map[string]interface{}{
		"pkg": g.pkg,
	}, filename, false)
	if err != nil {
		return err
	}

	log.Println("Done.")
	return nil
}

// Table defines mysql table
type Table struct {
	parser.Table
}

func (g *defaultGenerator) genModel(in parser.Table) (string, error) {
	if len(in.PrimaryKey.Name.Source()) == 0 {
		return "", fmt.Errorf("table %s: missing primary key", in.Name.Source())
	}

	var table Table
	table.Table = in

	importsCode, err := genImports(table, in.ContainsTime())
	if err != nil {
		return "", err
	}

	varsCode, err := genVars(table)
	if err != nil {
		return "", err
	}

	insertCode, insertCodeMethod, err := genInsert(table)
	if err != nil {
		return "", err
	}

	findCode := make([]string, 0)
	findOneCode, findOneCodeMethod, err := genFindOne(table)
	if err != nil {
		return "", err
	}
	findCode = append(findCode, findOneCode)

	updateCode, updateCodeMethod, err := genUpdate(table)
	if err != nil {
		return "", err
	}

	deleteCode, deleteCodeMethod, err := genDelete(table)
	if err != nil {
		return "", err
	}

	var list []string
	list = append(
		list,
		insertCodeMethod,
		findOneCodeMethod,
		updateCodeMethod,
		deleteCodeMethod,
	)

	typesCode, err := genTypes(table, strings.Join(trim.StringSlice(list), filex.NL))
	if err != nil {
		return "", err
	}

	newCode, err := genNew(table)
	if err != nil {
		return "", err
	}

	tableName, err := genTableName(table)
	if err != nil {
		return "", err
	}

	code := &code{
		importsCode: importsCode,
		varsCode:    varsCode,
		typesCode:   typesCode,
		newCode:     newCode,
		insertCode:  insertCode,
		findCode:    findCode,
		updateCode:  updateCode,
		deleteCode:  deleteCode,
		tableName:   tableName,
	}

	buffer, err := g.executeModel(table, code)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (g *defaultGenerator) executeModel(table Table, code *code) (*bytes.Buffer, error) {
	text, err := filex.LoadTemplate(modelGenTemplateFile, template.ModelGen)
	if err != nil {
		return nil, err
	}
	t := output.With("model").
		Parse(text).
		GoFmt(true)
	buffer, err := t.Execute(map[string]interface{}{
		"pkg":         g.pkg,
		"imports":     code.importsCode,
		"vars":        code.varsCode,
		"types":       code.typesCode,
		"new":         code.newCode,
		"insert":      code.insertCode,
		"find":        strings.Join(code.findCode, "\n"),
		"update":      code.updateCode,
		"delete":      code.deleteCode,
		"extraMethod": code.cacheExtra,
		"tableName":   code.tableName,
		"data":        table,
	})
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (g *defaultGenerator) genModelCustom(in parser.Table) (string, error) {
	text, err := filex.LoadTemplate(modelCustomTemplateFile, template.ModelCustom)
	if err != nil {
		return "", err
	}

	t := output.With("model-custom").
		Parse(text).
		GoFmt(true)
	buffer, err := t.Execute(map[string]interface{}{
		"pkg":                   g.pkg,
		"upperStartCamelObject": in.Name.ToCamel(),
		"lowerStartCamelObject": stringx.From(in.Name.ToCamel()).Untitle(),
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
