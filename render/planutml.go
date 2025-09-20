package render

import (
	"fmt"
	"io"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) renderPlantUMLScheme(w io.Writer, pkgInfos []model.Package) error {
	_, err := io.WriteString(w, "@startuml \"goarchlint\"\nskinparam componentStyle rectangle\n")
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
