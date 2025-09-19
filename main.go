package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gbh007/goarchlint/model"
	"github.com/gbh007/goarchlint/render"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

func main() {
	projectPath := flag.String("p", ".", "path to project")
	flag.Parse()

	pPath, err := filepath.Abs(*projectPath)
	if err != nil {
		panic(err)
	}

	modFilename := pPath + "/go.mod"

	modFile, err := os.Open(modFilename)
	if err != nil {
		panic(err)
	}

	defer modFile.Close()

	modFileData, err := io.ReadAll(modFile)
	if err != nil {
		panic(err)
	}

	mod, err := modfile.Parse(modFilename, modFileData, nil)
	if err != nil {
		panic(err)
	}

	var (
		coreModulePath string
		deps           = make(map[string]struct{}, len(mod.Require))
	)

	if mod.Module != nil {
		coreModulePath = mod.Module.Mod.Path
	}

	for _, module := range mod.Require {
		deps[module.Mod.Path] = struct{}{}
	}

	isInnerPath := func(p string) bool {
		_, ok := deps[p]
		if ok {
			return false
		}

		return strings.HasPrefix(p, pPath)
	}

	isInnerPkg := func(p string) bool {
		_, ok := deps[p]
		if ok {
			return false
		}

		return strings.HasPrefix(p, coreModulePath)
	}

	dirs := []string{}

	err = filepath.WalkDir(pPath, func(p string, d fs.DirEntry, err error) error {
		if d != nil && d.Name() == ".git" {
			return filepath.SkipDir
		}

		if d != nil && d.IsDir() {
			dirs = append(dirs, p)
		}

		return err
	})
	if err != nil {
		panic(err)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadImports,
		Tests: true,
		Dir:   pPath,
	}, dirs...)
	if err != nil {
		panic(err)
	}

	pkgInfos := []model.Package{}

	for _, pkg := range pkgs {
		for _, filename := range pkg.GoFiles {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
			if err != nil {
				panic(err)
			}

			pkgPath, _ := path.Split(filename)
			pkgPath = strings.TrimSuffix(pkgPath, "/")

			pkgInfo := model.Package{
				RelativePath: pkgPath,
				Inner:        isInnerPath(pkgPath),
			}

			if pkgInfo.Inner {
				pkgInfo.RelativePath = strings.ReplaceAll(pkgInfo.RelativePath, pPath, coreModulePath)
			}

			if f.Name != nil {
				pkgInfo.Name = f.Name.Name
				pkgInfo.IsMain = pkgInfo.Name == "main"
			}

			for _, node := range f.Imports {
				fpos := fset.Position(node.Pos())

				fileInfo := model.File{
					Path:   filename,
					Pos:    int(node.Pos()),
					Line:   fpos.Line,
					Column: fpos.Column,
				}

				importPath := ""

				if node.Path != nil {
					importPath = strings.Trim(node.Path.Value, "\"")
				}

				importInfo := model.Import{
					RelativePath: importPath,
					Files:        []model.File{fileInfo},
					Inner:        isInnerPkg(importPath),
				}

				pkgInfo.Imports = append(pkgInfo.Imports, importInfo)
			}

			pkgInfos = append(pkgInfos, pkgInfo)
		}
	}

	pkgInfos = model.CompactPackages(pkgInfos)

	slices.SortFunc(pkgInfos, func(a, b model.Package) int {
		return strings.Compare(a.RelativePath, b.RelativePath)
	})

	data, _ := json.MarshalIndent(pkgInfos, "", "  ")
	fmt.Println(string(data))

	outMD, _ := os.Create("out.md")
	defer outMD.Close()

	_ = render.RenderMermaidScheme(outMD, pkgInfos, true)

	outPuml, _ := os.Create("out.puml")
	defer outPuml.Close()

	_ = render.RenderPlantUMLScheme(outPuml, pkgInfos, true)
}
