package enst

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	boolType           = NewGoType(TypeBool, false)
	intType            = NewGoType(TypeInt, int(0))
	int8Type           = NewGoType(TypeInt8, int8(0))
	int16Type          = NewGoType(TypeInt16, int16(0))
	int32Type          = NewGoType(TypeInt32, int32(0))
	int64Type          = NewGoType(TypeInt64, int64(0))
	uintType           = NewGoType(TypeUint, uint(0))
	uint8Type          = NewGoType(TypeUint8, uint8(0))
	uint16Type         = NewGoType(TypeUint16, uint16(0))
	uint32Type         = NewGoType(TypeUint32, uint32(0))
	uint64Type         = NewGoType(TypeUint64, uint64(0))
	float32Type        = NewGoType(TypeFloat32, float32(0))
	float64Type        = NewGoType(TypeFloat64, float64(0))
	stringType         = NewGoType(TypeString, "")
	decimalType        = NewGoType(TypeDecimal, "")
	enumType           = NewGoType(TypeEnum, "")
	timeType           = NewGoType(TypeTime, time.Time{})
	bytesType          = NewGoType(TypeBytes, []byte(nil))
	sqlNullBoolType    = NewGoType(TypeBool, sql.NullBool{})
	sqlNullByteType    = NewGoType(TypeUint8, sql.NullByte{})
	sqlNullFloat64Type = NewGoType(TypeFloat64, sql.NullFloat64{})
	sqlNullInt16Type   = NewGoType(TypeInt16, sql.NullInt16{})
	sqlNullInt32Type   = NewGoType(TypeInt32, sql.NullInt32{})
	sqlNullInt64Type   = NewGoType(TypeInt64, sql.NullInt64{})
	sqlNullStringType  = NewGoType(TypeString, sql.NullString{})
	sqlNullTimeType    = NewGoType(TypeTime, sql.NullTime{})
	jSONRawMessageType = NewGoType(TypeJSON, json.RawMessage{})
	softDeleteType     = NewGoType(TypeInt64, soft_delete.DeletedAt(0))
	gormDeletedAtType  = NewGoType(TypeInt64, gorm.DeletedAt{})
	datatypesDateType  = NewGoType(TypeTime, datatypes.Date{})
	datatypesJSONType  = NewGoType(TypeJSON, datatypes.JSON{})
)

var sqlNullValueGoType = map[Type]GoType{
	TypeBool:    sqlNullBoolType,
	TypeUint8:   sqlNullByteType,
	TypeInt16:   sqlNullInt16Type,
	TypeInt32:   sqlNullInt32Type,
	TypeInt64:   sqlNullInt64Type,
	TypeFloat64: sqlNullFloat64Type,
	TypeString:  sqlNullStringType,
	TypeTime:    sqlNullTimeType,
}

func getSQLNullValueGoType(t Type) (GoType, bool) {
	v, ok := sqlNullValueGoType[t]
	return v, ok
}

func BoolType() GoType           { return boolType }
func IntType() GoType            { return intType }
func Int8Type() GoType           { return int8Type }
func Int16Type() GoType          { return int16Type }
func Int32Type() GoType          { return int32Type }
func Int64Type() GoType          { return int64Type }
func UintType() GoType           { return uintType }
func Uint8Type() GoType          { return uint8Type }
func Uint16Type() GoType         { return uint16Type }
func Uint32Type() GoType         { return uint32Type }
func Uint64Type() GoType         { return uint64Type }
func Float32Type() GoType        { return float32Type }
func Float64Type() GoType        { return float64Type }
func StringType() GoType         { return stringType }
func DecimalType() GoType        { return decimalType }
func EnumType() GoType           { return enumType }
func TimeType() GoType           { return timeType }
func BytesType() GoType          { return bytesType }
func SQLNullBoolType() GoType    { return sqlNullBoolType }
func SQLNullByteType() GoType    { return sqlNullByteType }
func SQLNullFloat64Type() GoType { return sqlNullFloat64Type }
func SQLNullInt16Type() GoType   { return sqlNullInt16Type }
func SQLNullInt32Type() GoType   { return sqlNullInt32Type }
func SQLNullInt64Type() GoType   { return sqlNullInt64Type }
func SQLNullStringType() GoType  { return sqlNullStringType }
func SQLNullTimeType() GoType    { return sqlNullTimeType }
func JSONRawMessageType() GoType { return jSONRawMessageType }
func SoftDeleteType() GoType     { return softDeleteType }
func GormDeletedAtType() GoType  { return gormDeletedAtType }
func DatatypesDateType() GoType  { return datatypesDateType }
func DatatypesJSONType() GoType  { return datatypesJSONType }
