package render

import (
	"fmt"
	"io"
	"path"

	"github.com/gbh007/goarchlint/model"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

type PackageConfig struct {
	RelativePath bool
	CurrentPath  string
	ShowInner    bool
}

func (r Render) RenderPackageTable(w io.Writer, pkgs []model.Package, cfg PackageConfig) error {
	tw := tablewriter.NewTable(
		w,
		tablewriter.WithRenderer(
			renderer.NewMarkdown(),
		),
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)

	row := []any{"Name", "Path"}
	if cfg.ShowInner {
		row = []any{"Name", "Path", "Inner"}
	}

	tw.Header(row...)

	for _, pkg := range pkgs {
		row := []any{}

		_, name := path.Split(pkg.RelativePath)
		row = append(row, name)

		link := r.getPackagePath(pkg)
		if cfg.RelativePath {
			link = r.resolvePath(cfg.CurrentPath, link)
		}

		innerValue := ""

		if r.PreferInnerNames && pkg.Inner {
			row = append(row, r.mdLink(pkg.InnerPath, link))
			innerValue = "✅"
		} else {
			row = append(row, r.mdLink(pkg.RelativePath, link))
			innerValue = "❌"
		}

		if cfg.ShowInner {
			row = append(row, innerValue)
		}

		err := tw.Append(row...)
		if err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}

	err := tw.Render()
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	return nil
}
