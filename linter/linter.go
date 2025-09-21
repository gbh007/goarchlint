package linter

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/gbh007/goarchlint/model"
	"github.com/samber/lo"
)

type RuleType byte

func (r RuleType) String() string {
	switch r {
	case RuleTypeStrict:
		return "strict"
	case RuleTypeLax:
		return "lax"
	default:
		return "unknown-" + strconv.Itoa(int(r))
	}
}

const (
	RuleTypeUnknown RuleType = iota
	RuleTypeStrict
	RuleTypeLax
)

type Rule struct {
	Path        string
	Allow       []string
	Deny        []string
	Type        RuleType
	OnlyInner   bool
	Description string

	pathPattern   *regexp.Regexp
	allowPatterns []*regexp.Regexp
	denyPatterns  []*regexp.Regexp
}

type Result struct {
	Strict int
	Lax    int
}

type violation struct {
	Type        RuleType
	PackagePath string
	Description string
	IsDeny      bool
}

func Validate(
	w io.Writer,
	pkgs []model.Package,
	rules []Rule,
	silentLax bool,
) (Result, error) {
	violations := make(map[model.File][]violation)

	rules = lo.Filter(rules, func(v Rule, _ int) bool {
		return v.Type == RuleTypeLax || v.Type == RuleTypeStrict
	})

	for i, rule := range rules {
		rules[i].pathPattern = compilePattern(rule.Path)

		rules[i].allowPatterns = lo.Map(rule.Allow, func(v string, _ int) *regexp.Regexp {
			return compilePattern(v)
		})

		rules[i].denyPatterns = lo.Map(rule.Deny, func(v string, _ int) *regexp.Regexp {
			return compilePattern(v)
		})
	}

	for _, pkg := range pkgs {
		if !pkg.Inner {
			continue
		}

		for _, rule := range rules {
			if !rule.pathPattern.MatchString(pkg.InnerPath) {
				continue
			}

			for _, imp := range pkg.Imports {
				if rule.OnlyInner && !imp.Inner {
					continue
				}

				path := imp.RelativePath
				if imp.Inner { // TODO: возможно такое стоит вынести в конфиг
					path = imp.InnerPath
				}

				found := false

				for _, pattern := range rule.denyPatterns {
					if pattern.MatchString(path) {
						found = true

						break
					}
				}

				if found {
					for _, file := range imp.Files {
						violations[file] = append(violations[file], violation{
							Type:        rule.Type,
							PackagePath: path,
							IsDeny:      true,
							Description: rule.Description,
						})
					}

					continue
				}

				if len(rule.allowPatterns) > 0 {
					found = true
				}

				for _, pattern := range rule.allowPatterns {
					if pattern.MatchString(path) {
						found = false

						break
					}
				}

				if found {
					for _, file := range imp.Files {
						violations[file] = append(violations[file], violation{
							Type:        rule.Type,
							PackagePath: path,
							Description: rule.Description,
						})
					}
				}
			}
		}
	}

	res := Result{}

	for file, vs := range violations {
		v := compactViolations(vs)

		switch v.Type {
		case RuleTypeStrict:
			res.Strict++
		case RuleTypeLax:
			res.Lax++

			if silentLax {
				continue
			}
		}

		_, err := io.WriteString(w, fmt.Sprintf(
			"%s:%d [%s] %s%s\n",
			file.Path,
			file.Line,
			v.Type.String(),
			lo.Ternary(
				v.IsDeny,
				"Import \""+v.PackagePath+"\" is deny",
				"Import \""+v.PackagePath+"\" not allowed",
			),
			lo.Ternary(v.Description != "", " > "+v.Description, ""),
		))
		if err != nil {
			return Result{}, fmt.Errorf("write violation: %w", err)
		}
	}

	return res, nil
}

func compilePattern(s string) *regexp.Regexp {
	s = regexp.QuoteMeta(s)
	s = strings.ReplaceAll(s, `\*\*`, `.+`)
	s = strings.ReplaceAll(s, `\*`, `[^/]+`)

	return regexp.MustCompile("^" + s + "$")
}

func compactViolations(vs []violation) violation {
	res := violation{}

	for _, v := range vs {
		if v.IsDeny {
			res = v

			break
		}

		if res.Type == RuleTypeUnknown {
			res = v

			continue
		}
	}

	return res
}
