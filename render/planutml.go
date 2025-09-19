package render

import (
	"fmt"
	"io"

	"github.com/gbh007/goarchlint/model"
)

func RenderPlantUMLScheme(w io.Writer, pkgInfos []model.Package, onlyInner, preferInnerNames bool) error {
	_, err := io.WriteString(w, "@startuml \"goarchlint\"\n")
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for _, pkg := range pkgInfos {
		if onlyInner && !pkg.Inner {
			continue
		}

		for _, imp := range pkg.Imports {
			if onlyInner && !imp.Inner {
				continue
			}

			from := pkg.RelativePath
			if preferInnerNames && pkg.InnerPath != "" {
				from = pkg.InnerPath
			}

			to := imp.RelativePath
			if preferInnerNames && imp.InnerPath != "" {
				to = imp.InnerPath
			}

			_, err = fmt.Fprintf(w, "\"%s\" }|--|| \"%s\" : x%d\n", from, to, len(imp.Files))
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
