package render

import (
	"fmt"
	"io"
	"os"
	"path"
	"slices"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) RenderMainDoc(name string, pkgs []model.Package) error {
	err := r.checkAndCreateDir(r.BasePath)
	if err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(path.Join(r.BasePath, "README.md"))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer f.Close()

	_, err = io.WriteString(f, "# "+name+"\n\n")
	if err != nil {
		return fmt.Errorf("write h1: %w", err)
	}

	var innerPackages, mainPackages []model.Package

	for _, pkg := range pkgs {
		if !pkg.Inner {
			continue
		}

		innerPackages = append(innerPackages, pkg)

		if pkg.IsMain {
			mainPackages = append(mainPackages, pkg)
		}
	}

	if len(mainPackages) > 0 {
		_, err = io.WriteString(f, "## Main packages\n\n")
		if err != nil {
			return fmt.Errorf("write main packages header: %w", err)
		}

		err = r.RenderPackageTable(f, mainPackages, PackageConfig{})
		if err != nil {
			return fmt.Errorf("write main packages table: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write main packages footer: %w", err)
		}
	}

	if len(innerPackages) > 0 {
		_, err = io.WriteString(f, "## Inner packages\n\n")
		if err != nil {
			return fmt.Errorf("write inner packages header: %w", err)
		}

		err = r.RenderPackageTable(f, innerPackages, PackageConfig{})
		if err != nil {
			return fmt.Errorf("write inner packages table: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write inner packages footer: %w", err)
		}
	}

	externalImports := []model.Import{}

	for _, pkg := range pkgs {
		for _, imp := range pkg.Imports {
			if !imp.Inner {
				externalImports = append(externalImports, imp.Clone())
			}
		}
	}

	if len(externalImports) > 0 {
		externalImports = model.CompactImports(externalImports)
		slices.SortFunc(externalImports, func(a, b model.Import) int {
			return len(b.Files) - len(a.Files)
		})

		_, err = io.WriteString(f, "## External imports\n\n")
		if err != nil {
			return fmt.Errorf("write external imports header: %w", err)
		}

		err = r.RenderImportTable(f, model.Package{Imports: externalImports}, ImportConfig{})
		if err != nil {
			return fmt.Errorf("write external imports table: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write external imports footer: %w", err)
		}
	}

	if r.Format != FormatNone {
		rCopy := r
		rCopy.MarkdownMode = true

		_, err = io.WriteString(f, "## Scheme\n\n")
		if err != nil {
			return fmt.Errorf("write scheme header: %w", err)
		}

		err = rCopy.RenderScheme(f, pkgs)
		if err != nil {
			return fmt.Errorf("write scheme: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write scheme footer: %w", err)
		}
	}

	err = r.renderFileFooter(f)
	if err != nil {
		return fmt.Errorf("write file footer: %w", err)
	}

	return nil
}
