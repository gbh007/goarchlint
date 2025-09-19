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

	outJSON, _ := os.Create("out.json")
	defer outJSON.Close()

	enc := json.NewEncoder(outJSON)
	enc.SetIndent("", "  ")

	err = enc.Encode(pkgInfos)
	if err != nil {
		panic(err)
	}

	outMD, err := os.Create("out.md")
	if err != nil {
		panic(err)
	}

	defer outMD.Close()

	err = render.RenderMermaidScheme(outMD, pkgInfos, true, true)
	if err != nil {
		panic(err)
	}

	outPuml, err := os.Create("out.puml")
	if err != nil {
		panic(err)
	}

	defer outPuml.Close()

	err = render.RenderPlantUMLScheme(outPuml, pkgInfos, true, true)
	if err != nil {
		panic(err)
	}
}
