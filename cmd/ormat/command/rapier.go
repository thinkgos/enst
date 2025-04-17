package command

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thinkgos/ens/rapier"
	"github.com/thinkgos/ens/utils"
)

type rapierOpt struct {
	source

	// output directory
	OutputDir string

	// codegen
	PackageName       string // required, proto 包名
	ModelImportPath   string // required, model导入路径
	DisableDocComment bool   // 禁用doc注释
	EnableInt         bool   // 使能int8,uint8,int16,uint16,int32,uint32输出为int,uint
	EnableBoolInt     bool   // 使能bool输出int
}

type rapierCmd struct {
	cmd *cobra.Command
	rapierOpt
}

func newRapierCmd() *rapierCmd {
	root := &rapierCmd{}

	cmd := &cobra.Command{
		Use:     "rapier",
		Short:   "Generate rapier from database/file",
		Example: "ormat rapier",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			rapierSchemaes := sc.IntoRapier()
			packageName := cmp.Or(root.PackageName, utils.GetPkgName(root.OutputDir))
			for _, entity := range rapierSchemaes.Entities {
				codegen := &rapier.CodeGen{
					Entities:          []*rapier.Struct{entity},
					ByName:            "ormat",
					Version:           version,
					PackageName:       packageName,
					ModelImportPath:   root.ModelImportPath,
					DisableDocComment: root.DisableDocComment,
					EnableInt:         root.EnableInt,
					EnableBoolInt:     root.EnableBoolInt,
				}
				data, err := codegen.Gen().FormatSource()
				if err != nil {
					return err
				}
				filename := joinFilename(root.OutputDir, entity.TableName, ".rapier.gen.go")
				err = WriteFile(filename, data)
				if err != nil {
					return fmt.Errorf("%v: %w", entity.TableName, err)
				}
				slog.Info("👉 " + filename)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")

	// database url
	cmd.Flags().StringVarP(&root.URL, "url", "u", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.Flags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table(仅url时有效)")
	cmd.Flags().StringSliceVarP(&root.Exclude, "exclude", "e", nil, "exclude table pattern(仅url时有效)")

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./repository", "out directory")

	cmd.Flags().StringVar(&root.PackageName, "package", "", "proto package name")
	cmd.Flags().StringVar(&root.ModelImportPath, "modelImportPath", "", "model导入路径")
	cmd.Flags().BoolVar(&root.DisableDocComment, "disableDocComment", false, "禁用文档注释")
	cmd.Flags().BoolVar(&root.EnableInt, "enableInt", false, "使能int8,uint8,int16,uint16,int32,uint32输出为int,uint")
	cmd.Flags().BoolVar(&root.EnableBoolInt, "enableBoolInt", false, "使能bool输出int")

	cmd.MarkFlagsOneRequired("url", "input")

	root.cmd = cmd
	return root
}
