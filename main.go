package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

func main() {
	projectPath := flag.String("p", ".", "path to project")
	flag.Parse()

	modFilename := *projectPath + "go.mod"

	modFile, err := os.Open(*projectPath + "/go.mod")
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

	// mod.Module.Mod.Path // Информация о корневом пакете
	// mod.Require[0].Mod.Path // Информация о зависимости

	_ = mod

	dirs := []string{}

	err = filepath.WalkDir(*projectPath, func(p string, d fs.DirEntry, err error) error {
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
		Dir:   *projectPath,
	}, dirs...)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, filename := range pkg.GoFiles {
			f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ImportsOnly)
			if err != nil {
				panic(err)
			}

			if f.Name != nil {
				fmt.Println("package", f.Name.Name, "file", filename)
			}

			for _, node := range f.Imports {
				localName := ""
				if node.Name != nil {
					localName = node.Name.Name
				}

				importPath := ""

				if node.Path != nil {
					importPath = node.Path.Value
				}

				fmt.Println(filename, ">", localName, "=", importPath)
			}
		}
	}
}
