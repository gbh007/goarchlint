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

type ImportConfig struct {
	RelativePath bool
	ShowInner    bool
}

func (r Render) RenderImportTable(w io.Writer, pkg model.Package, cfg ImportConfig) error {
	tw := tablewriter.NewTable(
		w,
		tablewriter.WithRenderer(
			renderer.NewMarkdown(),
		),
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)

	row := []any{"Name", "Path", "Count"}
	if cfg.ShowInner {
		row = []any{"Name", "Path", "Inner", "Count"}
	}

	tw.Header(row...)

	for _, imp := range pkg.Imports {
		row := []any{}

		_, name := path.Split(imp.RelativePath)
		row = append(row, name)

		link := r.getImportPath(imp)
		if cfg.RelativePath {
			link = r.resolvePath(r.getPackagePath(pkg), link)
		}

		innerValue := ""

		if r.PreferInnerNames && imp.Inner {
			row = append(row, r.mdLink(imp.InnerPath, link))
			innerValue = "✅"
		} else {
			row = append(row, r.mdLink(imp.RelativePath, link))
			innerValue = "❌"
		}

		if cfg.ShowInner {
			row = append(row, innerValue)
		}

		row = append(row, len(imp.Files))

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
