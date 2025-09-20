package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gbh007/goarchlint/parser"
	"github.com/gbh007/goarchlint/render"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Generator struct {
		Out              string
		OnlyInner        bool `mapstructure:"only_inner"`
		PreferInnerNames bool `mapstructure:"prefer_inner_names"`
		Format           string
		SchemeFile       string `mapstructure:"scheme_file"`
		Clean            bool
	}
}

func getVersion() string {
	info, ok := debug.ReadBuildInfo()
	if ok {
		return info.Main.Version
	}

	return ""
}

var (
	rootCmd = &cobra.Command{
		Use:          "goarchlint",
		Short:        "Go architecture linter and doc generator",
		SilenceUsage: true,
		Version:      getVersion(),
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

			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			dumpJSON, err := cmd.Flags().GetBool("dump-json")
			if err != nil {
				return err
			}

			viper.SetDefault("generator.out", outPath)
			viper.SetConfigFile(configPath)

			err = viper.ReadInConfig()
			if err != nil {
				return fmt.Errorf("read config: %w", err)
			}

			cfg := Config{}

			err = viper.Unmarshal(&cfg)
			if err != nil {
				return fmt.Errorf("unmarshal config: %w", err)
			}

			pkgInfos, module, err := parser.Parse(projectPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			r := render.Render{
				OnlyInner:        cfg.Generator.OnlyInner,
				PreferInnerNames: cfg.Generator.PreferInnerNames,
				MarkdownMode:     false,
				Format:           render.FormatFrom(cfg.Generator.Format),
				BasePath:         cfg.Generator.Out,
				SchemeFileFormat: render.FormatFrom(cfg.Generator.SchemeFile),
				CleanDir:         cfg.Generator.Clean,
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
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "generate config",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			return viper.WriteConfigAs(configPath)
		},
	}
)

func init() {
	viper.SetDefault("generator.out", "docs/arch")
	viper.SetDefault("generator.only_inner", true)
	viper.SetDefault("generator.prefer_inner_names", true)
	viper.SetDefault("generator.format", "mermaid")
	viper.SetDefault("generator.scheme_file", "plantuml")
	viper.SetDefault("generator.clean", false)

	rootCmd.PersistentFlags().StringP("config", "c", "goarchlint.toml", "path to config")
	rootCmd.PersistentFlags().StringP("path", "p", ".", "path to project")

	rootCmd.AddCommand(runCmd)

	generateCmd.Flags().Bool("dump-json", false, "Dump json")
	generateCmd.Flags().StringP("out", "o", "docs/arch", "path to doc output")
	rootCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(configCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
