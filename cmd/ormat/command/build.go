package command

import (
	"context"
	"log/slog"
	"os"

	"ariga.io/atlas/sql/schema"
	"github.com/spf13/cobra"
	"github.com/things-go/ens"
	"github.com/things-go/ens/driver"
)

type buildOpt struct {
	InputFile []string

	genFileOpt
}

type buildCmd struct {
	cmd    *cobra.Command
	Schema string
	buildOpt
}

func newBuildCmd() *buildCmd {
	root := &buildCmd{}

	getSchema := func() ens.Schemaer {
		innerParseFromFile := func(filename string) (ens.Schemaer, error) {
			content, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			d, err := driver.LoadDriver(root.Schema)
			if err != nil {
				return nil, err
			}
			return d.InspectSchema(context.Background(), &driver.InspectOption{
				URL:            "",
				Data:           string(content),
				InspectOptions: schema.InspectOptions{},
			})
		}

		mixin := &ens.MixinSchema{
			Name:     "",
			Entities: make([]ens.MixinEntity, 0, 128),
		}
		for _, filename := range root.InputFile {
			sc, err := innerParseFromFile(filename)
			if err != nil {
				slog.Warn("🧐 parse failed !!!", slog.String("file", filename), slog.Any("error", err))
				continue
			}
			mixin.Entities = append(mixin.Entities, sc.(*ens.MixinSchema).Entities...)
		}
		return mixin
	}

	cmd := &cobra.Command{
		Use:     "build",
		Short:   "Generate model from sql",
		Example: "ormat build",
		RunE: func(*cobra.Command, []string) error {
			sc := getSchema()
			return root.genFileOpt.GenModel(sc)
		},
	}

	cmdMapper := &cobra.Command{
		Use:     "mapper",
		Short:   "model mapper from database",
		Example: "ormat gen mapper",
		RunE: func(*cobra.Command, []string) error {
			sc := getSchema()
			return root.genFileOpt.GenMapper(sc)
		},
	}

	cmd.PersistentFlags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.PersistentFlags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser driver, [file+mysql,file+tidb]")
	cmd.PersistentFlags().StringVarP(&root.OutputDir, "out", "o", "./model", "out directory")

	InitFlagSetForConfig(cmd.PersistentFlags(), &root.View)

	cmd.PersistentFlags().BoolVar(&root.Merge, "merge", false, "merge in a file or not")
	cmd.PersistentFlags().StringVar(&root.MergeFilename, "filename", "", "merge filename")
	cmd.PersistentFlags().StringVar(&root.Template, "template", "", "use custom template")

	cmd.MarkPersistentFlagRequired("input") // nolint
	cmd.AddCommand(
		cmdMapper,
	)
	root.cmd = cmd
	return root
}
