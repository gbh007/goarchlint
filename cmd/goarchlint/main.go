package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gbh007/goarchlint/config"
	"github.com/gbh007/goarchlint/linter"
	"github.com/gbh007/goarchlint/parser"
	"github.com/gbh007/goarchlint/render"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:          "goarchlint",
		Short:        "Go architecture linter and doc generator",
		SilenceUsage: true,
		Version:      config.GetVersion(),
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
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("read config: %w", err)
			}

			cfg := config.Config{}

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
		Use:     "run",
		Aliases: []string{"lint"},
		Short:   "run lint",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectPath, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			configPath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			silentLax, err := cmd.Flags().GetBool("silent-lax")
			if err != nil {
				return err
			}

			err = setLinterRulesFromPreset(cmd)
			if err != nil {
				return err
			}

			viper.SetConfigFile(configPath)

			err = viper.ReadInConfig()
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("read config: %w", err)
			}

			cfg := config.Config{}

			err = viper.Unmarshal(&cfg)
			if err != nil {
				return fmt.Errorf("unmarshal config: %w", err)
			}

			pkgInfos, _, err := parser.Parse(projectPath)
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}

			result, err := linter.Validate(
				os.Stdout,
				pkgInfos,
				lo.Map(cfg.Linter.Rules, func(v config.LinterRule, _ int) linter.Rule {
					res := linter.Rule{
						Path:        v.Path,
						Allow:       v.Allow,
						Deny:        v.Deny,
						OnlyInner:   v.OnlyInner,
						Description: v.Description,
					}

					switch v.Type {
					case "strict":
						res.Type = linter.RuleTypeStrict
					case "lax":
						res.Type = linter.RuleTypeLax
					}

					return res
				}),
				silentLax,
			)
			if err != nil {
				return fmt.Errorf("validate: %w", err)
			}

			_, err = fmt.Fprintf(os.Stdout, "Violations:\n\tstrict: %d\n\tlax: %d\n", result.Strict, result.Lax)
			if err != nil {
				return fmt.Errorf("write result: %w", err)
			}

			if result.Strict > 0 {
				return errors.New("has strict violation")
			}

			return nil
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

			err = setLinterRulesFromPreset(cmd)
			if err != nil {
				return err
			}

			return viper.WriteConfigAs(configPath)
		},
	}
)

func setLinterRulesFromPreset(cmd *cobra.Command) error {
	linterPreset, err := cmd.Flags().GetString("linter")
	if err != nil {
		return err
	}

	switch linterPreset {
	case "hex":
		viper.SetDefault("linter.rules", config.TransformRulesToViper(config.LinterPresetHex()))
	case "hexlite":
		viper.SetDefault("linter.rules", config.TransformRulesToViper(config.LinterPresetHexLite()))
	case "clean":
		viper.SetDefault("linter.rules", config.TransformRulesToViper(config.LinterPresetClean()))
	}

	return nil
}

func init() {
	viper.SetDefault("generator.out", "docs/arch")
	viper.SetDefault("generator.only_inner", true)
	viper.SetDefault("generator.prefer_inner_names", true)
	viper.SetDefault("generator.format", "mermaid")
	viper.SetDefault("generator.scheme_file", "plantuml")
	viper.SetDefault("generator.clean", false)

	rootCmd.PersistentFlags().StringP("config", "c", "goarchlint.toml", "path to config")
	rootCmd.PersistentFlags().StringP("path", "p", ".", "path to project")
	rootCmd.PersistentFlags().String("linter", "", "linter preset, can be: hex, hexlite, clean")

	runCmd.Flags().Bool("silent-lax", false, "don't show lax violation")
	rootCmd.AddCommand(runCmd)

	generateCmd.Flags().Bool("dump-json", false, "Dump json")
	generateCmd.Flags().StringP("out", "o", "docs/arch", "path to doc output")

	generateCmd.Flags().Bool("generator.only_inner", true, "show on scheme only inner packages")
	generateCmd.Flags().Bool("generator.prefer_inner_names", true, "prefer to display inner names")
	generateCmd.Flags().String("generator.format", "mermaid", "format scheme in md")
	generateCmd.Flags().String("generator.scheme_file", "plantuml", "format general scheme")
	generateCmd.Flags().Bool("generator.clean", false, "clean all files in target directory")

	_ = viper.BindPFlag("generator.out", generateCmd.Flag("out"))
	_ = viper.BindPFlag("generator.only_inner", generateCmd.Flag("generator.only_inner"))
	_ = viper.BindPFlag("generator.prefer_inner_names", generateCmd.Flag("generator.prefer_inner_names"))
	_ = viper.BindPFlag("generator.format", generateCmd.Flag("generator.format"))
	_ = viper.BindPFlag("generator.scheme_file", generateCmd.Flag("generator.scheme_file"))
	_ = viper.BindPFlag("generator.clean", generateCmd.Flag("generator.clean"))

	rootCmd.AddCommand(generateCmd)

	rootCmd.AddCommand(configCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
