package render

import (
	"fmt"
	"io"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) renderMermaidScheme(w io.Writer, pkgInfos []model.Package) error {
	_, err := io.WriteString(w, "erDiagram\n")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, pkg := range pkgInfos {
		if r.OnlyInner && !pkg.Inner {
			continue
		}

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
