package render

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) renderMermaidScheme(w io.Writer, pkgInfos []model.Package) error {
	_, err := io.WriteString(w, "erDiagram\n")
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

			_, err = fmt.Fprintf(w, "    \"%s\" ||--|{ \"%s\" : x%d\n", from, to, len(imp.Files))
			if err != nil {
				return fmt.Errorf("write line: %w", err)
			}
		}
	}

	return nil
}
