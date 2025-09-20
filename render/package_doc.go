package render

import (
	"fmt"
	"io"
	"os"
	"path"
	"slices"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) RenderPackageDoc(pkg model.Package, pkgs []model.Package) error {
	p := r.getPackagePath(pkg)
	filePath := path.Join(r.BasePath, p)

	dirPath, _ := path.Split(filePath)

	err := r.checkAndCreateDir(dirPath)
	if err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer f.Close()

	_, err = io.WriteString(f, "# "+pkg.Name+"\n\n")
	if err != nil {
		return fmt.Errorf("write h1: %w", err)
	}

	if len(pkg.Imports) > 0 {
		pkg.Imports = slices.Clone(pkg.Imports)
		slices.SortFunc(pkg.Imports, func(a, b model.Import) int {
			return len(b.Files) - len(a.Files)
		})

		_, err = io.WriteString(f, "## Imports\n\n")
		if err != nil {
			return fmt.Errorf("write imports header: %w", err)
		}

		err = r.RenderImportTable(f, pkg, ImportConfig{
			RelativePath: true,
			ShowInner:    true,
		})
		if err != nil {
			return fmt.Errorf("write imports table: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write imports footer: %w", err)
		}
	}

	usedBy := model.FilterAndCleanPackageByImport(pkgs, pkg.RelativePath)

	if len(usedBy) > 0 {
		_, err = io.WriteString(f, "## Used by\n\n")
		if err != nil {
			return fmt.Errorf("write used by header: %w", err)
		}

		err = r.RenderPackageTable(f, usedBy, PackageConfig{
			RelativePath: true,
			CurrentPath:  p,
		})
		if err != nil {
			return fmt.Errorf("write used by table: %w", err)
		}

		_, err = io.WriteString(f, "\n")
		if err != nil {
			return fmt.Errorf("write used by footer: %w", err)
		}
	}

	if r.Format != FormatNone {
		usedBy = append(usedBy, pkg)

		rCopy := r
		rCopy.MarkdownMode = true

		_, err = io.WriteString(f, "## Scheme\n\n")
		if err != nil {
			return fmt.Errorf("write scheme header: %w", err)
		}

		err = rCopy.RenderScheme(f, usedBy)
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
