package carp

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/thinkgos/carp/utils"
	"gorm.io/gorm/schema"
)

var rowScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
var rowValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

func ParseModel(v any) (*EntityDescriptor, error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Pointer && value.IsNil() {
		return nil, nil
	}
	vt := indirect(value.Type())
	for vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s is not a struct", vt.String())
	}
	return &EntityDescriptor{
		Name:        utils.SnakeCase(vt.Name()),
		Comment:     "",
		Table:       nil,
		Fields:      structToFielders(vt),
		Indexes:     nil,
		ForeignKeys: nil,
	}, nil
}

func structToFielders(vt reflect.Type) []*FieldDescriptor {
	fields := make([]*FieldDescriptor, 0, vt.NumField())
	for i := 0; i < vt.NumField(); i++ {
		fv := vt.Field(i)
		if !fv.IsExported() { // ignore unexported field
			continue
		}
		tag := fv.Tag.Get("gorm")
		if tag == "-" { // ignore field
			continue
		}
		if fv.Anonymous { // an embedded field
			fvt := fv.Type
			for fvt.Kind() == reflect.Ptr {
				fvt = fv.Type.Elem()
			}
			if fvt.Kind() != reflect.Struct {
				continue
			}
			fields = append(
				fields,
				structToFielders(fvt)...,
			)
		} else {
			t, skip := intoGoTypeType(fv.Type, tag)
			if skip {
				continue
			}
			fields = append(
				fields,
				&FieldDescriptor{
					ColumnName: utils.SnakeCase(fv.Name),
					Comment:    "",
					Nullable:   false,
					Column:     nil,
					Type:       newGoType(t, fv.Type),
					GoName:     fv.Name,
					GoPointer:  false,
					Tags:       []string{},
				},
			)
		}
	}
	return fields
}

func intoGoTypeType(origTyp reflect.Type, tag string) (t Type, skip bool) {
	typ := indirect(origTyp)
	switch ident := typ.String(); typ.Kind() { // nolint: exhaustive
	case reflect.Bool:
		t = TypeBool
	case reflect.Int:
		t = TypeInt
	case reflect.Int8:
		t = TypeInt8
	case reflect.Int16:
		t = TypeInt16
	case reflect.Int32:
		t = TypeInt32
	case reflect.Int64:
		t = TypeInt64
	case reflect.Uint:
		t = TypeUint
	case reflect.Uint8:
		t = TypeUint8
	case reflect.Uint16:
		t = TypeUint16
	case reflect.Uint32:
		t = TypeUint32
	case reflect.Uint64:
		t = TypeUint64
	case reflect.Float32:
		t = TypeFloat32
	case reflect.Float64:
		t = TypeFloat64
	case reflect.String:
		typeValue := schema.ParseTagSetting(tag, ";")["TYPE"]
		if v := strings.ToUpper(typeValue); strings.Contains(v, "DECIMAL") || strings.Contains(v, "NUMERIC") {
			t = TypeDecimal
		} else {
			t = TypeString
		}
	case reflect.Struct:
		switch ident {
		case "time.Time",
			"sql.NullTime",
			"datatypes.Date":
			t = TypeTime
		case "sql.NullBool":
			t = TypeBool
		case "sql.NullByte":
			t = TypeBytes
		case "sql.NullString":
			t = TypeString
		case "sql.NullFloat64":
			t = TypeFloat64
		case "sql.NullInt16":
			t = TypeInt16
		case "sql.NullInt32":
			t = TypeInt32
		case "sql.NullInt64":
			t = TypeInt64
		default:
			t = TypeOther
			skip = !reflect.PointerTo(typ).Implements(rowScanner) || !typ.Implements(rowValuer)
		}
	case reflect.Slice:
		switch ident {
		case "json.RawMessage", "datatypes.JSON":
			t = TypeJSON
		case "[]uint8", "[]byte":
			t = TypeBytes
		default:
			skip = true
		}
	case reflect.Array:
		// TODO: ...
		t = TypeBytes
	default:
		t = TypeOther
	}
	return t, skip
}
