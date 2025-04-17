package mysql

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/thinkgos/enst"
	"github.com/thinkgos/enst/driver"
	"github.com/thinkgos/enst/internal/insql"
	"github.com/xwb1989/sqlparser"
)

var _ driver.Driver = (*SQL)(nil)

type SQL struct{}

// InspectSchema implements driver.Driver.
func (sq *SQL) InspectSchema(ctx context.Context, arg *driver.InspectOption) (*enst.Schema, error) {
	table, err := sq.inspectSchema(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &enst.Schema{
		Name:     "",
		Entities: []*enst.EntityDescriptor{intoSchema(table)},
	}, nil
}

func (sq *SQL) inspectSchema(_ context.Context, arg *driver.InspectOption) (*schema.Table, error) {
	statement, err := sqlparser.Parse(arg.Data)
	if err != nil {
		return nil, err
	}
	switch stmt := statement.(type) {
	case *sqlparser.DDL:
		if stmt.Action != sqlparser.CreateStr {
			return nil, errors.New("不是创建表语句")
		}
		if stmt.TableSpec == nil {
			return nil, errors.New("未解析到任何字段")
		}
		return parseSqlTable(stmt)
	default:
		return nil, errors.New("不是DDL语句")
	}
}

func parseSqlTable(stmt *sqlparser.DDL) (*schema.Table, error) {
	table := schema.NewTable(stmt.NewName.Name.String())

	//* table
	// ENGINE=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci comment='我是注释'
	tbOptions := strings.Split(stmt.TableSpec.Options, " ")
	for _, option := range tbOptions {
		keyValue := strings.Split(option, "=")
		if len(keyValue) >= 2 {
			switch keyValue[0] {
			case "ENGINE":
				table.AddAttrs(&mysql.Engine{V: keyValue[1], Default: false})
			case "charset":
				table.AddAttrs(&schema.Charset{V: keyValue[1]})
			case "collate":
				table.AddAttrs(&schema.Collation{V: keyValue[1]})
			case "comment":
				table.AddAttrs(&schema.Comment{Text: strings.ReplaceAll(keyValue[1], "'", "")})
			}
		}
	}

	//* columns
	columns := make([]*schema.Column, 0, len(stmt.TableSpec.Columns))
	for _, col := range stmt.TableSpec.Columns {
		column, err := parseSqlColumnDefinition(col)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	table.AddColumns(columns...)

	// * indexes
	indexes := make([]*schema.Index, 0, len(stmt.TableSpec.Indexes))
	for _, idx := range stmt.TableSpec.Indexes {
		index, err := parseSqlIndexDefinition(table, idx, columns)
		if err != nil {
			return nil, err
		}
		indexes = append(indexes, index)
	}
	table.AddIndexes(indexes...)
	return table, nil
}

func parseSqlColumnDefinition(col *sqlparser.ColumnDefinition) (*schema.Column, error) {
	coldef := schema.NewColumn(col.Name.String())
	colType := col.Type
	if colType.Default != nil {
		coldef.Default = &schema.Literal{V: string(colType.Default.Val)}
	}
	if colType.OnUpdate != nil {
		coldef.AddAttrs(&mysql.OnUpdate{
			Attr: nil,
			A:    string(colType.OnUpdate.Val),
		})
	}
	if colType.Autoincrement {
		coldef.AddAttrs(&mysql.AutoIncrement{})
	}

	if colType.Comment != nil {
		coldef.AddAttrs(&schema.Comment{Text: string(colType.Comment.Val)})
	}
	parseInt := func(l *sqlparser.SQLVal) int64 {
		if l == nil {
			return 0
		}
		length, err := strconv.ParseInt(string(l.Val), 0, 0)
		if err != nil {
			return 0
		}
		return length
	}

	isUnsigned := bool(colType.Unsigned)
	switch colType.Type {
	case mysql.TypeBool,
		mysql.TypeBoolean,
		mysql.TypeTinyInt,
		mysql.TypeSmallInt,
		mysql.TypeMediumInt,
		mysql.TypeInt,
		mysql.TypeBigInt,
		"integer":
		// `^(tinyint)\b[(]1[)] unsigned`
		// `^(tinyint)\b[(]1[)]`
		// `^(tinyint)\b([(]\d+[)])? unsigned`
		// `^(tinyint)\b([(]\d+[)])?`
		// `^(smallint)\b([(]\d+[)])? unsigned`
		// `^(smallint)\b([(]\d+[)])?`
		// `^(mediumint)\b([(]\d+[)])? unsigned`
		// `^(mediumint)\b([(]\d+[)])?`
		// `^(int)\b([(]\d+[)])? unsigned`
		// `^(int)\b([(]\d+[)])?`
		// `^(bigint)\b([(]\d+[)])? unsigned`
		// `^(bigint)\b([(]\d+[)])?`
		// `^(integer)\b([(]\d+[)])? unsigned`
		// `^(integer)\b([(]\d+[)])?`
		rawColumnType := func() string {
			length := int64(0)
			if colType.Length != nil {
				length = parseInt(colType.Length)
			}
			if length > 0 {
				if isUnsigned {
					return fmt.Sprintf("%s(%d) unsigned", colType.Type, length)
				} else {
					return fmt.Sprintf("%s(%d)", colType.Type, length)
				}
			}
			if isUnsigned {
				return fmt.Sprintf("%s unsigned", colType.Type)
			} else {
				return colType.Type
			}
		}

		coldef.Type = &schema.ColumnType{
			Type: &schema.IntegerType{
				T:        colType.Type,
				Unsigned: isUnsigned,
			},
			Raw:  rawColumnType(),
			Null: !bool(colType.NotNull),
		}

	case mysql.TypeDecimal,
		mysql.TypeNumeric,
		mysql.TypeFloat,
		mysql.TypeDouble,
		mysql.TypeReal:
		// `^(float)\b([(]\d+,\d+[)])? unsigned`
		// `^(float)\b([(]\d+,\d+[)])?`
		// `^(double)\b([(]\d+,\d+[)])? unsigned`
		// `^(double)\b([(]\d+,\d+[)])?`
		// `^(decimal)\b[(]\d+,\d+[)]`
		length, scale := parseInt(colType.Length), parseInt(colType.Scale)
		raw := colType.Type
		if length > 0 {
			raw = fmt.Sprintf("%s(%d,%d)", colType.Type, length, scale)
		}

		if colType.Type == mysql.TypeDecimal ||
			colType.Type == mysql.TypeNumeric {
			coldef.Type = &schema.ColumnType{
				Type: &schema.DecimalType{
					T:         colType.Type,
					Precision: int(length),
					Scale:     int(scale),
					Unsigned:  isUnsigned,
				},
				Raw:  raw,
				Null: !bool(colType.NotNull),
			}
		} else {
			coldef.Type = &schema.ColumnType{
				Type: &schema.FloatType{
					T:         colType.Type,
					Unsigned:  isUnsigned,
					Precision: int(length),
				},
				Raw:  raw,
				Null: !bool(colType.NotNull),
			}
		}

	case mysql.TypeVarchar,
		mysql.TypeChar,
		mysql.TypeText,
		mysql.TypeTinyText,
		mysql.TypeMediumText,
		mysql.TypeLongText:
		// `^(char)\b[(]\d+[)]`
		// `^(varchar)\b[(]\d+[)]`
		// `^(text)\b([(]\d+[)])?`
		// `^(tinytext)\b([(]\d+[)])?`
		// `^(mediumtext)\b([(]\d+[)])?`
		// `^(longtext)\b([(]\d+[)])?`
		raw := colType.Type
		length := parseInt(colType.Length)
		if length > 0 {
			raw = fmt.Sprintf("%s(%d)", colType.Type, length)
		}

		coldef.Type = &schema.ColumnType{
			Type: &schema.StringType{
				T:    colType.Type,
				Size: int(length),
			},
			Raw:  raw,
			Null: !bool(colType.NotNull),
		}
	case
		mysql.TypeVarBinary,
		mysql.TypeBinary,
		mysql.TypeBlob,
		mysql.TypeTinyBlob,
		mysql.TypeMediumBlob,
		mysql.TypeLongBlob,
		mysql.TypeBit:
		// `^(blob)\b([(]\d+[)])?`
		// `^(tinyblob)\b([(]\d+[)])?`
		// `^(mediumblob)\b([(]\d+[)])?`
		// `^(longblob)\b([(]\d+[)])?`
		// `^(binary)\b[(]\d+[)]`
		// `^(varbinary)\b[(]\d+[)]`
		// `^(bit)\b[(]\d+[)]`
		var size *int

		raw := colType.Type
		if colType.Length != nil {
			length := parseInt(colType.Length)
			size = insql.P(int(length))
			raw = fmt.Sprintf("%s(%d)", colType.Type, length)
		}

		coldef.Type = &schema.ColumnType{
			Type: &schema.BinaryType{
				T:    colType.Type,
				Size: size,
			},
			Raw:  raw,
			Null: !bool(colType.NotNull),
		}

	case mysql.TypeTimestamp,
		mysql.TypeDate,
		mysql.TypeTime,
		mysql.TypeDateTime,
		mysql.TypeYear:
		// `^(date)\b([(]\d+[)])?`
		// `^(datetime)\b([(]\d+[)])?`
		// `^(timestamp)\b([(]\d+[)])?`
		// `^(time)\b([(]\d+[)])?`
		// `^(year)\b([(]\d+[)])?`
		var precision, scale *int

		if colType.Length != nil {
			precision = insql.P(int(parseInt(colType.Length)))
		}
		if colType.Scale != nil {
			scale = insql.P(int(parseInt(colType.Scale)))
		}
		coldef.Type = &schema.ColumnType{
			Type: &schema.TimeType{
				T:         colType.Type,
				Precision: precision,
				Scale:     scale,
			},
			Raw:  colType.Type,
			Null: !bool(colType.NotNull),
		}
	case mysql.TypeEnum:
		// `^(enum)\b[(](.)+[)]`
		coldef.Type = &schema.ColumnType{
			Type: &schema.EnumType{
				T:      colType.Type,
				Values: colType.EnumValues,
			},
			Raw:  "enum(" + strings.Join(colType.EnumValues, ",") + ")",
			Null: !bool(colType.NotNull),
		}
	case mysql.TypeJSON:
		// `^(json)\b`
		coldef.Type = &schema.ColumnType{
			Type: &schema.JSONType{
				T: colType.Type,
			},
			Raw:  colType.Type,
			Null: !bool(colType.NotNull),
		}
	case mysql.TypeGeometry,
		mysql.TypePoint,
		mysql.TypeMultiPoint,
		mysql.TypeLineString,
		mysql.TypeMultiLineString,
		mysql.TypePolygon,
		mysql.TypeMultiPolygon,
		mysql.TypeGeoCollection,
		mysql.TypeGeometryCollection:
		// `geometry`
		coldef.Type = &schema.ColumnType{
			Type: &schema.SpatialType{
				T: colType.Type,
			},
			Raw:  colType.Type,
			Null: !bool(colType.NotNull),
		}
	default:
		return nil, fmt.Errorf("unsupported column type(%s)", colType.Type)
	}
	return coldef, nil
}

func parseSqlIndexDefinition(table *schema.Table, idx *sqlparser.IndexDefinition, columns []*schema.Column) (*schema.Index, error) {
	indexName := idx.Info.Name.String()
	isPk := false
	unique := false
	switch idx.Info.Type {
	case "primary key":
		isPk = true
		unique = true
	case "unique key", "unique":
		unique = true
	case "key":
		unique = false
	}
	indexType := "BTREE"
	for _, option := range idx.Options {
		if option.Name == "using" {
			indexType = option.Using
		}
	}

	cols := make([]*schema.Column, 0, len(idx.Columns))
	for _, idxCol := range idx.Columns {
		columnName := idxCol.Column.String()
		col, ok := insql.FindColumn(columns, columnName)
		if ok {
			cols = append(cols, col)
		} else {
			return nil, fmt.Errorf("Key('%s') column '%s' doesn't exist in table '%s'", indexName, columnName, table.Name)
		}
	}
	index := schema.NewIndex(indexName).
		SetUnique(unique).
		AddAttrs(&mysql.IndexType{T: indexType}).
		AddColumns(cols...)
	if isPk {
		table.SetPrimaryKey(index)
	}
	return index, nil
}
