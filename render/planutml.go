package render

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) renderPlantUMLScheme(w io.Writer, pkgInfos []model.Package) error {
	_, err := io.WriteString(w, "@startuml \"goarchlint\"\nskinparam componentStyle rectangle\n")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	pkgInfos = slices.Clone(pkgInfos)
	slices.SortStableFunc(pkgInfos, func(a, b model.Package) int {
		return strings.Compare(a.RelativePath, b.RelativePath)
	})

	for _, pkg := range pkgInfos {
		if r.OnlyInner && !pkg.Inner {
			continue
		}

		pkg.Imports = slices.Clone(pkg.Imports)
		slices.SortStableFunc(pkg.Imports, func(a, b model.Import) int {
			return strings.Compare(a.RelativePath, b.RelativePath)
		})

		for _, imp := range pkg.Imports {
			if r.OnlyInner && !imp.Inner {
				continue
			}

			from := pkg.RelativePath
			if r.PreferInnerNames && pkg.InnerPath != "" {
				from = pkg.InnerPath
			}

			to := imp.RelativePath
			if r.PreferInnerNames && imp.InnerPath != "" {
				to = imp.InnerPath
			}

			_, err = fmt.Fprintf(w, "[%s] --> [%s]\n", from, to)
			if err != nil {
				return fmt.Errorf("write line: %w", err)
			}
		}
	}

	_, err = io.WriteString(w, "@enduml")
	if err != nil {
		return fmt.Errorf("write footer: %w", err)
	}

	return nil
}
