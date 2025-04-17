package command

import (
	"cmp"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thinkgos/ens"
	"github.com/thinkgos/ens/utils"
)

type modelOpt struct {
	source

	OutputDir string

	PackageName string // 包名

	ens.Option
	DisableCommentTag bool              // 禁用注释放入tag标签中
	DisableDocComment bool              // 禁用文档注释
	CustomFieldIdent  map[string]string // 自定义字段类型, 格式: TableName.ColumnName->Ident
	Merge             bool
	MergeFilename     string
}

type modelCmd struct {
	cmd *cobra.Command
	modelOpt
}

func newModelCmd() *modelCmd {
	root := &modelCmd{}

	cmd := &cobra.Command{
		Use:     "model",
		Short:   "Generate model from database",
		Example: "ormat model",
		RunE: func(*cobra.Command, []string) error {
			schemaes, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			packageName := cmp.Or(root.PackageName, utils.GetPkgName(root.OutputDir))
			customFieldIdent := map[string]map[string]string{}
			for key, val := range root.CustomFieldIdent {
				ks := strings.Split(key, ".")
				if len(ks) != 2 || val == "" {
					continue
				}
				tb := ks[0]
				field := ks[1]
				if tb == "" || field == "" {
					continue
				}
				fields, ok := customFieldIdent[tb]
				if !ok {
					fields = map[string]string{}
				}
				fields[field] = val
				customFieldIdent[tb] = fields
			}

			if root.Merge {
				g := ens.CodeGen{
					Entities:          schemaes.Entities,
					ByName:            "ormat",
					Version:           version,
					PackageName:       packageName,
					DisableDocComment: root.DisableDocComment,
					CustomFieldIdent:  customFieldIdent,
					Option:            root.Option,
				}
				data, err := g.Gen().FormatSource()
				if err != nil {
					return err
				}
				filename := joinFilename(root.OutputDir, root.MergeFilename, ".go")
				err = WriteFile(filename, data)
				if err != nil {
					return err
				}
				slog.Info("👉 " + filename)
			} else {
				for _, entity := range schemaes.Entities {
					g := &ens.CodeGen{
						Entities:          []*ens.EntityDescriptor{entity},
						ByName:            "ormat",
						Version:           version,
						PackageName:       packageName,
						DisableDocComment: root.DisableDocComment,
						CustomFieldIdent:  customFieldIdent,
						Option:            root.Option,
					}
					data, err := g.Gen().FormatSource()
					if err != nil {
						return fmt.Errorf("%v: %v", entity.Name, err)
					}
					filename := joinFilename(root.OutputDir, entity.Name, ".go")
					err = WriteFile(filename, data)
					if err != nil {
						return fmt.Errorf("%v: %v", entity.Name, err)
					}
					slog.Info("👉 " + filename)
				}
			}
			slog.Info("😄 generate success !!!")
			return nil
		},
	}
	// input file
	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")
	// database url
	cmd.Flags().StringVarP(&root.URL, "url", "u", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.Flags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table")
	cmd.Flags().StringSliceVarP(&root.Exclude, "exclude", "e", nil, "exclude table pattern")

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	cmd.Flags().StringVar(&root.PackageName, "package", "", "package name")
	cmd.Flags().BoolVar(&root.DisableDocComment, "disableDocComment", false, "禁用文档注释")

	cmd.Flags().BoolVar(&root.IgnoreOmitempty, "ignoreOmitempty", false, "忽略tags标签的 omitempty 标签")
	cmd.Flags().StringToStringVar(&root.Tags, "tags", map[string]string{"json": utils.StyleSnakeCase}, "tags标签,类型支持[smallCamelCase,pascalCase,snakeCase,kebab]")
	cmd.Flags().BoolVar(&root.EnableInt, "enableInt", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	cmd.Flags().BoolVar(&root.EnableBoolInt, "enableBoolInt", false, "使能bool输出int")
	cmd.Flags().BoolVar(&root.DisableNullToPoint, "disableNullToPoint", false, "禁用字段为null时输出指针类型,将输出为sql.Nullxx")
	cmd.Flags().BoolVar(&root.DisableCommentTag, "disableCommentTag", false, "禁用注释放入tag标签中")
	cmd.Flags().BoolVar(&root.EnableForeignKey, "enableForeignKey", false, "使用外键")
	cmd.Flags().StringSliceVar(&root.EscapeName, "escapeName", nil, "escape name list")
	cmd.Flags().StringToStringVar(&root.CustomFieldIdent, "customFieldIdent", map[string]string{}, "自定义字段类型, 格式: TableName.ColumnName=Ident")

	cmd.Flags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.Flags().StringVar(&root.MergeFilename, "filename", "", "merge filename")

	cmd.MarkFlagsOneRequired("url", "input")

	root.cmd = cmd
	return root
}
