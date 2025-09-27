package render

import (
	"fmt"
	"os"
	"path"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) RenderSchemeDoc(pkgs []model.Package) error {
	rCopy := r
	rCopy.MarkdownMode = false
	rCopy.Format = rCopy.SchemeFileFormat

	var filename string

	switch rCopy.Format {
	case FormatMermaid:
		filename = "scheme.mmd"
	case FormatPlantUML:
		filename = "scheme.puml"
	default:
		return nil
	}

	f, err := os.Create(path.Join(r.BasePath, filename))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer f.Close()

	err = rCopy.RenderScheme(f, pkgs)
	if err != nil {
		return fmt.Errorf("write scheme: %w", err)
	}

	return nil
}
