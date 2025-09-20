package render

import (
	"fmt"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) RenderDocs(name string, pkgs []model.Package) error {
	err := r.RenderMainDoc("TEST", pkgs)
	if err != nil {
		return fmt.Errorf("main doc: %w", err)
	}

	if r.SchemeFileFormat != FormatNone {
		err = r.RenderSchemeDoc(pkgs)
		if err != nil {
			return fmt.Errorf("scheme doc: %w", err)
		}
	}

	for _, pkg := range pkgs {
		if !pkg.Inner {
			continue
		}

		err = r.RenderPackageDoc(pkg, pkgs)
		if err != nil {
			return fmt.Errorf("package doc: %w", err)
		}
	}

	return nil
}
