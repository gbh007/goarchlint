package main

import (
	"encoding/json"
	"flag"
	"os"

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

	outJSON, _ := os.Create("out/out.json")
	defer outJSON.Close()

	enc := json.NewEncoder(outJSON)
	enc.SetIndent("", "  ")

	err = enc.Encode(pkgInfos)
	if err != nil {
		panic(err)
	}

	outMMD, err := os.Create("out/out.mmd")
	if err != nil {
		panic(err)
	}

	defer outMMD.Close()

	r := render.Render{
		OnlyInner:        true,
		PreferInnerNames: true,
		MarkdownMode:     false,
		Format:           render.FormatMermaid,
		BasePath:         "out",
		SchemeFileFormat: render.FormatPlantUML,
	}

	err = r.RenderScheme(outMMD, pkgInfos)
	if err != nil {
		panic(err)
	}

	err = r.RenderMainDoc("TEST", pkgInfos)
	if err != nil {
		panic(err)
	}

	r.Format = render.FormatPlantUML

	outPuml, err := os.Create("out/out.puml")
	if err != nil {
		panic(err)
	}

	defer outPuml.Close()

	err = r.RenderScheme(outPuml, pkgInfos)
	if err != nil {
		panic(err)
	}

	outMD, err := os.Create("out/out.md")
	if err != nil {
		panic(err)
	}

	defer outMD.Close()

	err = r.RenderPackageTable(outMD, pkgInfos, render.PackageConfig{})
	if err != nil {
		panic(err)
	}
}
