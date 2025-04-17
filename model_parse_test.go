package enst_test

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/plugin/soft_delete"
)

type Anonymous struct {
	AnonymousField int
}

type TestData struct {
	Id          int64                 `gorm:"column:id;not null;autoIncrement:true;primaryKey;comment:主键"`
	StrVal      string                `gorm:"column:str_val;type:varchar(255);not null"`
	StrNullVal1 sql.NullString        `gorm:"column:str_null_val1;type:varchar(2048);null"`
	StrNullVal2 *string               `gorm:"column:str_null_val2;type:varchar(2048);null"`
	Value1      float64               `gorm:"column:value1;type:double;not null;default:'';uniqueIndex:uk_key_value,priority:1"`
	Value2      *float64              `gorm:"column:value2;type:double;null;default:'';uniqueIndex:uk_key_value,priority:1"`
	Priority    uint                  `gorm:"column:priority;type:unsigned int(10) unsigned;not null;default:255"`
	Visible     bool                  `gorm:"column:visible;type:tinyint(1) unsigned;not null;default:0"`
	Time1       time.Time             `gorm:"column:time1;type:datetime;not null"`
	Time2       *time.Time            `gorm:"column:time2;type:datetime;null"`
	Time3       sql.NullTime          `gorm:"column:time3;type:datetime;null"`
	DeletedAt   soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint(20);not null;default:0"`
	Js          json.RawMessage
	Bs          []byte
	// unexported             int
	OmitField              int `gorm:"-"`
	Unsupported            []time.Time
	NotImplScanerAndValuer Anonymous
	Anonymous
}

// func Test(t *testing.T) {
// 	entity, _ := enst.ParseModel(TestData{})
// 	des := entity.Build(nil)
// 	g := codegen.
// 		New(
// 			[]*enst.EntityDescriptor{des},
// 			codegen.WithPackageName("ens"),
// 		).
// 		GenRapier("")

// 	bytes, err := g.FormatSource()
// 	if err != nil {
// 		t.Log(err)
// 		t.Log(string(g.Bytes()))
// 	} else {
// 		os.WriteFile(des.Name+".rapier.go", bytes, 0644)
// 	}
// }
