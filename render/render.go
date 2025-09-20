package render

import (
	"errors"
	"fmt"
	"io"

	"github.com/gbh007/goarchlint/model"
)

type Format byte

const (
	FormatUnknown Format = iota
	FormatMermaid
	FormatPlantUML
)

type Render struct {
	OnlyInner        bool
	PreferInnerNames bool
	MarkdownMode     bool
	Format           Format
}

func (r Render) Render(w io.Writer, pkgInfos []model.Package) error {
	var (
		renderFunc func(w io.Writer, pkgInfos []model.Package) error
		mdType     string
	)

	switch r.Format {
	case FormatMermaid:
		renderFunc = r.renderMermaidScheme
		mdType = "mermaid"
	case FormatPlantUML:
		renderFunc = r.renderPlantUMLScheme
		mdType = "plantuml"
	default:
		return errors.New("unsupported format")
	}

	if r.MarkdownMode {
		_, err := io.WriteString(w, "```"+mdType+"\n")
		if err != nil {
			return fmt.Errorf("write code block header: %w", err)
		}
	}

	err := renderFunc(w, pkgInfos)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	if r.MarkdownMode {
		_, err := io.WriteString(w, "```")
		if err != nil {
			return fmt.Errorf("write code block footer: %w", err)
		}
	}

	return nil
}
