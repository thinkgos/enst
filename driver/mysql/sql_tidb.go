package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/thinkgos/enst"
	"github.com/thinkgos/enst/driver"
	"github.com/thinkgos/enst/internal/insql"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	parsermysql "github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/parser/types"

	_ "github.com/pingcap/tidb/parser/test_driver"
)

type SQLTidb struct{}

// InspectSchema implements driver.Driver.
func (st *SQLTidb) InspectSchema(_ context.Context, arg *driver.InspectOption) (*enst.Schema, error) {
	pr := parser.New()
	stmts, _, err := pr.ParseSQL(arg.Data)
	if err != nil {
		return nil, err
	}

	entities := make([]*enst.EntityDescriptor, 0, len(stmts))
	for _, stmt := range stmts {
		createStmt, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			continue
		}
		table, err := parserCreateTableStmtTable(createStmt)
		if err != nil {
			return nil, err
		}
		entities = append(entities, intoSchema(table))
	}
	return &enst.Schema{
		Name:     "",
		Entities: entities,
	}, nil
}

func parserCreateTableStmtTable(stmt *ast.CreateTableStmt) (*schema.Table, error) {
	table := schema.NewTable(stmt.Table.Name.L)

	//* table
	// ENGINE=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci comment='我是注释'
	for _, option := range stmt.Options {
		switch option.Tp { // nolint: exhaustive
		case ast.TableOptionEngine:
			table.AddAttrs(&mysql.Engine{V: option.StrValue, Default: false})
		case ast.TableOptionCharset:
			table.AddAttrs(&schema.Charset{V: option.StrValue})
		case ast.TableOptionCollate:
			table.AddAttrs(&schema.Collation{V: option.StrValue})
		case ast.TableOptionComment:
			table.AddAttrs(&schema.Comment{Text: option.StrValue})
		}
	}
	//* columns
	columns := make([]*schema.Column, 0, len(stmt.Cols))
	for _, col := range stmt.Cols {
		column, err := parserCreateTableStmtColumn(col)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	table.AddColumns(columns...)
	// * indexes
	indexes := make([]*schema.Index, 0, len(stmt.Table.IndexHints))
	for _, idx := range stmt.Constraints {
		index, err := parseCreateTableStmtIndex(table, idx, columns)
		if err != nil {
			return nil, err
		}
		if index == nil {
			continue
		}
		indexes = append(indexes, index)
	}
	table.AddIndexes(indexes...)
	return table, nil
}

func parserCreateTableStmtColumn(col *ast.ColumnDef) (*schema.Column, error) {
	coldef := schema.NewColumn(col.Name.Name.L)
	nullable := true
	for _, opt := range col.Options {
		switch opt.Tp { // nolint: exhaustive
		case ast.ColumnOptionNotNull:
			nullable = false
		case ast.ColumnOptionAutoIncrement:
			coldef.AddAttrs(&mysql.AutoIncrement{})
		case ast.ColumnOptionDefaultValue:
			val := formatExprNode(opt.Expr)
			if !strings.EqualFold(val, "null") {
				coldef.Default = &schema.Literal{V: strings.Trim(formatExprNode(opt.Expr), `"`)}
			}

		case ast.ColumnOptionComment:
			coldef.AddAttrs(&schema.Comment{Text: strings.Trim(formatExprNode(opt.Expr), `"`)})
		case ast.ColumnOptionOnUpdate:
			coldef.AddAttrs(&mysql.OnUpdate{
				Attr: nil,
				A:    strings.Trim(formatExprNode(opt.Expr), `"`),
			})
		}
	}

	ft := col.Tp
	ftFlag := ft.GetFlag()
	ftType := ft.GetType()
	ftTypeValue := types.TypeToStr(ftType, ft.GetCharset())
	ftTypeRaw := col.Tp.InfoSchemaStr()
	isUnsigned := parsermysql.HasUnsignedFlag(ftFlag)
	length := ft.GetFlen()
	decimal := ft.GetDecimal()
	// log.Printf("%v ---> %v--> %v -- %v -- %v --- %v -- %v\n", ftType, ft, ftTypeValue, ftTypeRaw, isUnsigned, length, decimal)

	switch ftType {
	case parsermysql.TypeTiny, // [bool, boolean, tinyint]
		parsermysql.TypeShort,
		parsermysql.TypeInt24,
		parsermysql.TypeLong,
		parsermysql.TypeLonglong:
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

		coldef.Type = &schema.ColumnType{
			Type: &schema.IntegerType{
				T:        ftTypeValue,
				Unsigned: isUnsigned,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}

	case parsermysql.TypeFloat, // [float]
		parsermysql.TypeDouble,     // [real, double]
		parsermysql.TypeNewDecimal: // [numeric, decimal]
		// `^(float)\b([(]\d+,\d+[)])? unsigned`
		// `^(float)\b([(]\d+,\d+[)])?`
		// `^(double)\b([(]\d+,\d+[)])? unsigned`
		// `^(double)\b([(]\d+,\d+[)])?`
		// `^(decimal)\b[(]\d+,\d+[)]`
		// `^(decimal)\b[(]\d+,\d+[)]`
		if ftType == parsermysql.TypeNewDecimal {
			coldef.Type = &schema.ColumnType{
				Type: &schema.DecimalType{
					T:         ftTypeValue,
					Precision: length,
					Scale:     decimal,
					Unsigned:  isUnsigned,
				},
				Raw:  ftTypeRaw,
				Null: nullable,
			}
		} else {
			coldef.Type = &schema.ColumnType{
				Type: &schema.FloatType{
					T:         ftTypeValue,
					Unsigned:  isUnsigned,
					Precision: length,
				},
				Raw:  ftTypeRaw,
				Null: nullable,
			}
		}

	case parsermysql.TypeVarchar, // [varchar, varbinary]
		parsermysql.TypeVarString,
		parsermysql.TypeString: //[char, binary]
		// `^(char)\b[(]\d+[)]`
		// `^(varchar)\b[(]\d+[)]`
		// `^(binary)\b[(]\d+[)]`
		// `^(varbinary)\b[(]\d+[)]`

		coldef.Type = &schema.ColumnType{
			Type: &schema.StringType{
				T:    ftTypeValue,
				Size: length,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}
	case parsermysql.TypeBlob, // [text]
		parsermysql.TypeTinyBlob,   // [tinytext, blob, tinyblob]
		parsermysql.TypeMediumBlob, // [mediumtext, mediumblob]
		parsermysql.TypeLongBlob,   // [longtext, longblob]
		parsermysql.TypeBit:        // [bit]
		// `^(blob)\b([(]\d+[)])?`
		// `^(tinyblob)\b([(]\d+[)])?`
		// `^(mediumblob)\b([(]\d+[)])?`
		// `^(longblob)\b([(]\d+[)])?`
		// `^(text)\b([(]\d+[)])?`
		// `^(tinytext)\b([(]\d+[)])?`
		// `^(mediumtext)\b([(]\d+[)])?`
		// `^(longtext)\b([(]\d+[)])?`
		// `^(bit)\b[(]\d+[)]`
		var size *int
		if length != -1 {
			size = &length
		}

		coldef.Type = &schema.ColumnType{
			Type: &schema.BinaryType{
				T:    ftTypeValue,
				Size: size,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}

	case parsermysql.TypeTimestamp, // [timestamp]
		parsermysql.TypeDate,     // [date]
		parsermysql.TypeDuration, // [time]
		parsermysql.TypeDatetime, // [datetime]
		parsermysql.TypeYear,     // [year]
		parsermysql.TypeNewDate:
		// `^(date)\b([(]\d+[)])?`
		// `^(datetime)\b([(]\d+[)])?`
		// `^(timestamp)\b([(]\d+[)])?`
		// `^(time)\b([(]\d+[)])?`
		// `^(year)\b([(]\d+[)])?`
		var precision, scale *int

		if length != -1 {
			precision = &length
		}
		if decimal != -1 {
			scale = &decimal
		}
		coldef.Type = &schema.ColumnType{
			Type: &schema.TimeType{
				T:         ftTypeValue,
				Precision: precision,
				Scale:     scale,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}
	case parsermysql.TypeEnum,
		parsermysql.TypeSet:
		// `^(enum)\b[(](.)+[)]`
		// `^(set)\b[(](.)+[)]`
		coldef.Type = &schema.ColumnType{
			Type: &schema.EnumType{
				T:      ftTypeValue,
				Values: col.Tp.GetElems(),
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}
	case parsermysql.TypeJSON:
		// `^(json)\b`
		coldef.Type = &schema.ColumnType{
			Type: &schema.JSONType{
				T: ftTypeValue,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}
	case parsermysql.TypeGeometry: // FIXME: not support
		// `geometry`
		coldef.Type = &schema.ColumnType{
			Type: &schema.SpatialType{
				T: ftTypeValue,
			},
			Raw:  ftTypeRaw,
			Null: nullable,
		}
	default:
		return nil, fmt.Errorf("unsupported column type(%s)", ftTypeRaw)
	}
	return coldef, nil
}

func parseCreateTableStmtIndex(table *schema.Table, idx *ast.Constraint, columns []*schema.Column) (*schema.Index, error) {
	indexName := idx.Name
	isPk := false
	unique := false
	switch idx.Tp { // nolint: exhaustive
	case ast.ConstraintPrimaryKey:
		indexName = "PRIMARY"
		isPk = true
		unique = true
	case ast.ConstraintUniq,
		ast.ConstraintUniqKey,
		ast.ConstraintUniqIndex:
		unique = true
	case ast.ConstraintKey,
		ast.ConstraintIndex:
		unique = false
	default:
		return nil, nil
	}
	indexType := "BTREE"
	if idx.Option != nil {
		indexType = idx.Option.Tp.String()
	}

	cols := make([]*schema.Column, 0, len(idx.Keys))
	for _, idxCol := range idx.Keys {
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

func formatExprNode(e ast.ExprNode) string {
	if e == nil {
		return ""
	}
	b := &strings.Builder{}
	e.Format(b)
	return b.String()
}
