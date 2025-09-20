package render

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/gbh007/goarchlint/model"
)

type Format byte

const (
	FormatNone Format = iota
	FormatMermaid
	FormatPlantUML
)

type Render struct {
	OnlyInner        bool
	PreferInnerNames bool
	MarkdownMode     bool
	Format           Format
	BasePath         string
	SchemeFileFormat Format
}

func (r Render) getPackagePath(p model.Package) string {
	s := p.RelativePath
	if r.PreferInnerNames {
		s = p.InnerPath
	}

	s, f := path.Split(s)
	s = strings.Trim(s, "/")

	return path.Join(s, f+".md")
}

func (r Render) getImportPath(p model.Import) string {
	s := p.RelativePath
	if r.PreferInnerNames {
		s = p.InnerPath
	}

	s, f := path.Split(s)
	s = strings.Trim(s, "/")

	return path.Join(s, f+".md")
}

func (Render) resolvePath(current, target string) string {
	s, err := filepath.Rel(current, target)
	if err != nil {
		panic(err)
	}

	return s
}

func (Render) mdLink(name, uri string) string {
	if uri == ".md" { // FIXME: избавится от этого костыля
		return name
	}

	return fmt.Sprintf("[%s](%s)", name, uri)
}
