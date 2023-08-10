package codegen

import (
	"fmt"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

func (g *CodeGen) GenAssist(modelImportPath string) *CodeGen {
	pkgQualifierPrefix := ""
	if p := ens.PkgName(modelImportPath); p != "" {
		pkgQualifierPrefix = p + "."
	}
	if !g.disableDocComment {
		g.
			P("// Code generated by ", g.byName, ". DO NOT EDIT.").
			P("// version: ", g.version).
			P()
	}
	g.
		P("package ", g.packageName).
		P()

	//* import
	g.P("import (")
	if pkgQualifierPrefix != "" {
		g.P(`"`, modelImportPath, `"`)
		g.P()
	}
	g.P(`assist "github.com/things-go/gorm-assist"`)
	g.P(`"gorm.io/gorm"`)
	g.P(")")

	//* struct
	constFieldFn := func(stName, fieldName string) string {
		return fmt.Sprintf(`xx_%s_%s`, stName, fieldName)
	}
	constFieldWithTableNameFn := func(stName, fieldName string) string {
		return fmt.Sprintf(`xx_%[1]s_%[2]s_WithTableName`, stName, fieldName)
	}
	for _, et := range g.entities {
		structName := utils.CamelCase(et.Name)
		tableName := et.Name

		constTableName := fmt.Sprintf("xx_%s_TableName", structName)
		{ //* const field
			g.P("const (")
			g.P("// hold model `", structName, "` table name")
			g.P(constTableName, ` = "`, tableName, `"`)
			g.P("// hold model `", structName, "` column name")
			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				g.P(constFieldFn(structName, fieldName), ` = "`, field.Name, `"`)
			}
			g.P("// hold model `", structName, "` column name with table name(`", tableName, "`) prefix")
			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				g.P(
					constFieldWithTableNameFn(structName, fieldName),
					" = ",
					constTableName, ` + "_" + `, constFieldFn(structName, fieldName),
				)
			}
			g.P(")")
			g.P()
		}

		varNativeModel := fmt.Sprintf(`xxx_%s_Native_Model`, structName)
		varModel := fmt.Sprintf(`xxx_%s_Model`, structName)
		fnInnerNew := fmt.Sprintf(`new_%s`, structName)
		{ //* var field
			g.
				P("var ", varNativeModel, " = ", fnInnerNew, `("")`).
				P("var ", varModel, " = ", fnInnerNew, "(", constTableName, ")").
				P()
		}
		typeNative := fmt.Sprintf("%s_Native", structName)
		//* type
		{
			g.P("type ", typeNative, " struct {")
			g.P("xAlias string")
			g.P("ALL assist.Asterisk")
			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				g.P(fieldName, ` assist.`, field.AssistDataType)
			}
			g.P("}")
			g.P()
		}
		//* function X_xxx
		{
			g.
				P("// X_", structName, " model with TableName `", tableName, "`.").
				P("func X_", structName, "() ", typeNative, " {").
				P("return ", varModel).
				P("}").
				P()
		}
		//* function new_xxx
		{
			g.
				P("func ", fnInnerNew, "(xAlias string) ", typeNative, " {").
				P("return ", typeNative, " {").
				P("xAlias: xAlias,").
				P("ALL:  assist.NewAsterisk(xAlias),")
			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				g.P(fieldName, ": assist.New", field.AssistDataType, "(xAlias, ", constFieldFn(structName, fieldName), "),")
			}
			g.
				P("}").
				P("}").
				P()
		}
		//* function New_xxxx
		{
			g.
				P("// New_", structName, " new instance.").
				P("func New_", structName, "(xAlias string) ", typeNative, " {").
				P("switch xAlias {").
				P(`case "":`).
				P("return ", varNativeModel).
				P("case ", constTableName, ":").
				P("return ", varModel).
				P("default:").
				P("return ", fnInnerNew, "(xAlias)").
				P("}").
				P("}").
				P()
		}
		//* method As
		{
			g.
				P("// As alias").
				P("func (*", typeNative, ") As(alias string) ", typeNative, " {").
				P("return New_", structName, "(alias)").
				P("}").
				P()
		}
		//* method X_Alias
		{
			g.
				P("// X_Alias hold table name when call New_", structName, " or ", structName, "_Impl.As that you defined.").
				P("func (x *", typeNative, ") X_Alias() string {").
				P("return x.xAlias").
				P("}").
				P()
		}
		// table and column field
		{
			//* method TableName
			g.P("// TableName hold model `", structName, "` table name returns `", tableName, "`.").
				P("func (*", typeNative, ") TableName() string {").
				P("return ", constTableName).
				P("}").
				P()

			for _, field := range et.Fields {
				fieldName := utils.CamelCase(field.Name)
				columnName := field.Name
				//* method Field_xxx
				g.
					P("// Field_", fieldName, " hold model `", structName, "` column name.").
					P("// if prefixes not exist returns `", columnName, "`, others `{prefixes[0]}_", columnName, "`").
					P("func (*", typeNative, ") Field_", fieldName, "(prefixes ...string) string {").
					P("if len(prefixes) == 0 {").
					P("return ", constFieldFn(structName, fieldName)).
					P("}").
					P("if prefixes[0] == ", constTableName, " {").
					P("return ", constFieldWithTableNameFn(structName, fieldName)).
					P("}").
					P(`return prefixes[0] + "_" + `, constFieldFn(structName, fieldName)).
					P("}").
					P()
			}
		}
		{ // other method
			genAssistOtherImpl(g, et, typeNative, structName, pkgQualifierPrefix)
		}
	}
	return g
}

func genAssistOtherImpl(g *CodeGen, et *ens.EntityDescriptor, typeNative, structName, pkgQualifierPrefix string) {
	modelName := pkgQualifierPrefix + structName
	//* method New_Executor
	{
		g.
			P("// New_Executor new entity executor which suggest use only once.").
			P("func (*", typeNative, ") New_Executor(db *gorm.DB) *assist.Executor[", modelName, "]  {").
			P("return assist.NewExecutor[", modelName, "](db)").
			P("}").
			P()
	}
	//* method Select_Expr
	{
		g.
			P("// Select_Expr select model fields").
			P("func (x *", typeNative, ") Select_Expr() []assist.Expr {").
			P("return []assist.Expr{")
		for _, field := range et.Fields {
			g.P("x.", utils.CamelCase(field.Name), ",")
		}
		g.
			P("}").
			P("}").
			P()
	}

	//* method Select_VariantExpr
	{
		g.
			P("// Select_VariantExpr select model fields, but time.Time field convert to timestamp(int64).").
			P("func (x *", typeNative, ") Select_VariantExpr(prefixes ...string) []assist.Expr {").
			P("if len(prefixes) > 0 {").
			P("return []assist.Expr{")
		for _, field := range et.Fields {
			g.P(genAssist_SelectVariantExprField(structName, field, true))
		}
		g.
			P("}").
			P("} else {").
			P("return []assist.Expr{")
		for _, field := range et.Fields {
			g.P(genAssist_SelectVariantExprField(structName, field, false))
		}
		g.
			P("}").
			P("}").
			P("}").
			P()
	}
}

func genAssist_SelectVariantExprField(structName string, field *ens.FieldDescriptor, hasPrefix bool) string {
	fieldName := utils.CamelCase(field.Name)

	b := strings.Builder{}
	b.Grow(64)
	b.WriteString("x.")
	b.WriteString(fieldName)
	if field.Type.IsTime() {
		b.WriteString(".UnixTimestamp()")
		if field.Nullable {
			b.WriteString(".IfNull(0)")
		}
		if !hasPrefix {
			b.WriteString(".As(")
			b.WriteString(fmt.Sprintf(`xx_%s_%s`, structName, fieldName))
			b.WriteString(")")
		}
	}
	if hasPrefix {
		b.WriteString(".As(x.Field_")
		b.WriteString(fieldName)
		b.WriteString("(prefixes...))")
	}
	b.WriteString(",")
	return b.String()
}
