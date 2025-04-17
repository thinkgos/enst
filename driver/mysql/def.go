package mysql

import (
	"fmt"
	"regexp"

	"ariga.io/atlas/sql/schema"
	"github.com/thinkgos/enst"
)

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []struct {
	Key     string
	NewType func() enst.GoType
}{
	{`^(bool)`, enst.BoolType},                                 // bool
	{`^(tinyint)\b[(]1[)] unsigned`, enst.BoolType},            // bool
	{`^(tinyint)\b[(]1[)]`, enst.BoolType},                     // bool
	{`^(tinyint)\b([(]\d+[)])? unsigned`, enst.Uint8Type},      // uint8
	{`^(tinyint)\b([(]\d+[)])?`, enst.Int8Type},                // int8
	{`^(smallint)\b([(]\d+[)])? unsigned`, enst.Uint16Type},    // uint16
	{`^(smallint)\b([(]\d+[)])?`, enst.Int16Type},              // int16
	{`^(mediumint)\b([(]\d+[)])? unsigned`, enst.Uint32Type},   // uint32
	{`^(mediumint)\b([(]\d+[)])?`, enst.Int32Type},             // int32
	{`^(int)\b([(]\d+[)])? unsigned`, enst.Uint32Type},         // uint32
	{`^(int)\b([(]\d+[)])?`, enst.Int32Type},                   // int32
	{`^(integer)\b([(]\d+[)])? unsigned`, enst.Uint32Type},     // uint32
	{`^(integer)\b([(]\d+[)])?`, enst.Int32Type},               // int32
	{`^(bigint)\b([(]\d+[)])? unsigned`, enst.Uint64Type},      // uint64
	{`^(bigint)\b([(]\d+[)])?`, enst.Int64Type},                // int64
	{`^(float)\b([(]\d+,\d+[)])? unsigned`, enst.Float32Type},  // float32
	{`^(float)\b([(]\d+,\d+[)])?`, enst.Float32Type},           // float32
	{`^(double)\b([(]\d+,\d+[)])? unsigned`, enst.Float64Type}, // float64
	{`^(double)\b([(]\d+,\d+[)])?`, enst.Float64Type},          // float64
	{`^(char)\b[(]\d+[)]`, enst.StringType},                    // string
	{`^(varchar)\b[(]\d+[)]`, enst.StringType},                 // string
	{`^(datetime)\b([(]\d+[)])?`, enst.TimeType},               // time.Time
	{`^(date)\b([(]\d+[)])?`, enst.TimeType},                   // datatypes.Date
	{`^(timestamp)\b([(]\d+[)])?`, enst.TimeType},              // time.Time
	{`^(time)\b([(]\d+[)])?`, enst.TimeType},                   // time.Time
	{`^(year)\b([(]\d+[)])?`, enst.TimeType},                   // time.Time
	{`^(text)\b([(]\d+[)])?`, enst.StringType},                 // string
	{`^(tinytext)\b([(]\d+[)])?`, enst.StringType},             // string
	{`^(mediumtext)\b([(]\d+[)])?`, enst.StringType},           // string
	{`^(longtext)\b([(]\d+[)])?`, enst.StringType},             // string
	{`^(blob)\b([(]\d+[)])?`, enst.BytesType},                  // []byte
	{`^(tinyblob)\b([(]\d+[)])?`, enst.BytesType},              // []byte
	{`^(mediumblob)\b([(]\d+[)])?`, enst.BytesType},            // []byte
	{`^(longblob)\b([(]\d+[)])?`, enst.BytesType},              // []byte
	{`^(bit)\b[(]\d+[)]`, enst.BytesType},                      // []uint8
	{`^(json)\b`, enst.JSONRawMessageType},                     // datatypes.JSON
	{`^(enum)\b[(](.)+[)]`, enst.StringType},                   // string
	{`^(set)\b[(](.)+[)]`, enst.StringType},                    // string
	{`^(decimal)\b[(]\d+,\d+[)]`, enst.DecimalType},            // string
	{`^(binary)\b[(]\d+[)]`, enst.BytesType},                   // []byte
	{`^(varbinary)\b[(]\d+[)]`, enst.BytesType},                // []byte
	{`^(geometry)`, enst.StringType},                           // string
}

func intoGoType(columnType string) enst.GoType {
	for _, v := range typeDictMatchList {
		ok, _ := regexp.MatchString(v.Key, columnType)
		if ok {
			return v.NewType()
		}
	}
	panic(fmt.Sprintf("type (%v) not match in any way, need to add on (https://github.com/thinkgos/ormat/blob/main/driver/mysql/def.go)", columnType))
}

type TableDef struct {
	tb *schema.Table
}

func NewTableDef(tb *schema.Table) enst.TableDef {
	return &TableDef{tb: tb}
}

func (d *TableDef) Table() *schema.Table { return d.tb }

func (d *TableDef) PrimaryKey() enst.IndexDef {
	if d.tb.PrimaryKey != nil {
		return NewIndexDef(d.tb.PrimaryKey)
	}
	return nil
}

func (d *TableDef) Definition() string {
	return intoTableSql(d.tb)
}

type ColumnDef struct {
	col *schema.Column
}

func NewColumnDef(col *schema.Column) enst.ColumnDef {
	return &ColumnDef{col: col}
}

func (d *ColumnDef) Column() *schema.Column { return d.col }

func (d *ColumnDef) Definition() string { return intoColumnSql(d.col) }

func (d *ColumnDef) GormTag(tb *schema.Table) string { return intoGormTag(tb, d.col) }

type IndexDef struct {
	index *schema.Index
}

func NewIndexDef(index *schema.Index) enst.IndexDef { return &IndexDef{index: index} }

func (d *IndexDef) Index() *schema.Index { return d.index }

func (d *IndexDef) Definition() string { return intoIndexSql(d.index) }

type ForeignKeyDef struct {
	fk *schema.ForeignKey
}

func NewForeignKey(fk *schema.ForeignKey) enst.ForeignKeyDef {
	return &ForeignKeyDef{fk: fk}
}

func (d *ForeignKeyDef) ForeignKey() *schema.ForeignKey { return d.fk }

func (d *ForeignKeyDef) Definition() string { return intoForeignKeySql(d.fk) }
