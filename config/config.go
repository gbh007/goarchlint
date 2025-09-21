package config

import (
	"runtime/debug"

	"github.com/samber/lo"
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
	Linter struct {
		Rules []LinterRule
	}
}

type LinterRule struct {
	Path        string
	Allow       []string
	Deny        []string
	Type        string
	OnlyInner   bool `mapstructure:"only_inner"`
	Description string
}

func (lr LinterRule) ToMap() map[string]any {
	return map[string]any{
		"path":        lr.Path,
		"allow":       lr.Allow,
		"deny":        lr.Deny,
		"type":        lr.Type,
		"only_inner":  lr.OnlyInner,
		"description": lr.Description,
	}
}

func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if ok {
		return info.Main.Version
	}

	return ""
}

func TransformRulesToViper(rules []LinterRule) []map[string]any {
	return lo.Map(rules, func(lr LinterRule, _ int) map[string]any {
		return lr.ToMap()
	})
}
