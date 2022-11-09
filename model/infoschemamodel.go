package model

import (
	"fmt"
	"github.com/dabao-zhao/ddltomodel/util/trim"
	"gorm.io/gorm"
	"sort"
)

const indexPri = "PRIMARY"

type (
	// InformationSchemaModel defines information schema model
	InformationSchemaModel struct {
		conn *gorm.DB
	}

	// Column defines column in table
	Column struct {
		*DbColumn
		Index *DbIndex
	}

	// DbColumn defines column info of columns
	DbColumn struct {
		Name            string      `gorm:"column:COLUMN_NAME"`
		DataType        string      `gorm:"column:DATA_TYPE"`
		ColumnType      string      `gorm:"column:COLUMN_TYPE"`
		Extra           string      `gorm:"column:EXTRA"`
		Comment         string      `gorm:"column:COLUMN_COMMENT"`
		ColumnDefault   interface{} `gorm:"column:COLUMN_DEFAULT"`
		IsNullAble      string      `gorm:"column:IS_NULLABLE"`
		OrdinalPosition int         `gorm:"column:ORDINAL_POSITION"`
	}

	// DbIndex defines index of columns in information_schema.statistic
	DbIndex struct {
		IndexName  string `gorm:"column:INDEX_NAME"`
		NonUnique  int    `gorm:"column:NON_UNIQUE"`
		SeqInIndex int    `gorm:"column:SEQ_IN_INDEX"`
	}

	// ColumnData describes the columns of table
	ColumnData struct {
		Db      string
		Table   string
		Columns []*Column
	}

	// Table describes mysql table which contains database name, table name, columns, keys
	Table struct {
		Db      string
		Table   string
		Columns []*Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}

	// IndexType describes an alias of string
	IndexType string

	// Index describes a column index
	Index struct {
		IndexType IndexType
		Columns   []*Column
	}
)

// NewInformationSchemaModel creates an instance for InformationSchemaModel
func NewInformationSchemaModel(conn *gorm.DB) *InformationSchemaModel {
	return &InformationSchemaModel{conn: conn}
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *InformationSchemaModel) GetAllTables(database string) ([]string, error) {
	var tables []string

	err := m.conn.Table("TABLES").
		Select("TABLE_NAME").
		Where("TABLE_SCHEMA = ?", database).
		Find(&tables).Error
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// FindColumns return columns in specified database and table
func (m *InformationSchemaModel) FindColumns(db, table string) (*ColumnData, error) {
	var reply []*DbColumn

	err := m.conn.Table("COLUMNS c").
		Select("c.COLUMN_NAME,c.DATA_TYPE,c.COLUMN_TYPE,EXTRA,c.COLUMN_COMMENT,c.COLUMN_DEFAULT,c.IS_NULLABLE,c.ORDINAL_POSITION").
		Where("c.TABLE_SCHEMA = ?", db).
		Where("c.TABLE_NAME = ?", table).
		Find(&reply).Error

	if err != nil {
		return nil, err
	}

	var list []*Column
	for _, item := range reply {
		index, err := m.FindIndex(db, table, item.Name)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
			continue
		}

		if len(index) > 0 {
			for _, i := range index {
				list = append(list, &Column{
					DbColumn: item,
					Index:    i,
				})
			}
		} else {
			list = append(list, &Column{
				DbColumn: item,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].OrdinalPosition < list[j].OrdinalPosition
	})

	var columnData ColumnData
	columnData.Db = db
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

// FindIndex finds index with given db, table and column.
func (m *InformationSchemaModel) FindIndex(db, table, column string) ([]*DbIndex, error) {
	var reply []*DbIndex

	err := m.conn.Table("STATISTICS s").
		Select("s.INDEX_NAME,s.NON_UNIQUE,s.SEQ_IN_INDEX").
		Where("c.TABLE_SCHEMA = ?", db).
		Where("c.TABLE_NAME = ?", table).
		Where("c.COLUMN_NAME = ?", column).
		Find(&reply).Error

	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Convert converts column data into Table
func (c *ColumnData) Convert() (*Table, error) {
	var table Table
	table.Table = c.Table
	table.Db = c.Db
	table.Columns = c.Columns
	table.UniqueIndex = map[string][]*Column{}
	table.NormalIndex = map[string][]*Column{}

	m := make(map[string][]*Column)
	for _, each := range c.Columns {
		each.Comment = trim.NewLine(each.Comment)
		if each.Index != nil {
			m[each.Index.IndexName] = append(m[each.Index.IndexName], each)
		}
	}

	primaryColumns := m[indexPri]
	if len(primaryColumns) == 0 {
		return nil, fmt.Errorf("db:%s, table:%s, missing primary key", c.Db, c.Table)
	}

	if len(primaryColumns) > 1 {
		return nil, fmt.Errorf("db:%s, table:%s, joint primary key is not supported", c.Db, c.Table)
	}

	table.PrimaryKey = primaryColumns[0]
	for indexName, columns := range m {
		if indexName == indexPri {
			continue
		}

		for _, one := range columns {
			if one.Index != nil {
				if one.Index.NonUnique == 0 {
					table.UniqueIndex[indexName] = columns
				} else {
					table.NormalIndex[indexName] = columns
				}
			}
		}
	}

	return &table, nil
}
