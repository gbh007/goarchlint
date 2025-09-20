package main

import (
	"flag"

	"github.com/gbh007/goarchlint/parser"
	"github.com/gbh007/goarchlint/render"
)

func main() {
	projectPath := flag.String("p", ".", "path to project")
	flag.Parse()

	pkgInfos, err := parser.Parse(*projectPath)
	if err != nil {
		panic(err)
	}

	r := render.Render{
		OnlyInner:        true,
		PreferInnerNames: true,
		MarkdownMode:     false,
		Format:           render.FormatMermaid,
		BasePath:         "out",
		SchemeFileFormat: render.FormatPlantUML,
	}

	err = r.RenderDocs("TEST", pkgInfos)
	if err != nil {
		panic(err)
	}

	err = r.DumpJSON(pkgInfos)
	if err != nil {
		panic(err)
	}
}
