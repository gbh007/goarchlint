package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gbh007/goarchlint/parser"
	"github.com/gbh007/goarchlint/render"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "goarchlint",
		Short:        "Go architecture linter and doc generator",
		SilenceUsage: true,
	}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			outPath, err := cmd.Flags().GetString("out")
			if err != nil {
				return err
			}

			// configPath, err := cmd.Flags().GetString("config")
			// if err != nil {
			// 	return err
			// }

			dumpJSON, err := cmd.Flags().GetBool("dump-json")
			if err != nil {
				return err
			}
			pkgInfos, module, err := parser.Parse(projectPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			r := render.Render{
				OnlyInner:        true,
				PreferInnerNames: true,
				MarkdownMode:     false,
				Format:           render.FormatMermaid,
				BasePath:         outPath,
				SchemeFileFormat: render.FormatPlantUML,
				CleanDir:         true,
			}

			err = r.RenderDocs(module, pkgInfos)
			if err != nil {
				return fmt.Errorf("render: %w", err)
			}

			if dumpJSON {
				err = r.DumpJSON(pkgInfos)
				if err != nil {
					return fmt.Errorf("dump: %w", err)
				}
			}

			return nil
		},
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "run lint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("unimplemented yet")
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "goarchlint.toml", "path to config")
	rootCmd.PersistentFlags().StringP("path", "p", ".", "path to project")

	rootCmd.AddCommand(runCmd)

	generateCmd.Flags().Bool("dump-json", false, "Dump json")
	generateCmd.Flags().StringP("out", "o", "docs/arch", "path to doc output")
	rootCmd.AddCommand(generateCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
