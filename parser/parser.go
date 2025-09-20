package parser

import (
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
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

func Parse(projectPath string) ([]model.Package, error) {
	pPath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, fmt.Errorf("get absolute path: %w", err)
	}

	modFilename := pPath + "/go.mod"

	modFile, err := os.Open(modFilename)
	if err != nil {
		return nil, fmt.Errorf("open go.mod file: %w", err)
	}

	defer modFile.Close()

	modFileData, err := io.ReadAll(modFile)
	if err != nil {
		return nil, fmt.Errorf("read go.mod file: %w", err)
	}

	mod, err := modfile.Parse(modFilename, modFileData, nil)
	if err != nil {
		return nil, fmt.Errorf("parse go.mod file: %w", err)
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
		return nil, fmt.Errorf("walk project files: %w", err)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.LoadImports,
		Tests: true,
		Dir:   pPath,
	}, dirs...)
	if err != nil {
		return nil, fmt.Errorf("load package infos: %w", err)
	}

	pkgInfos := []model.Package{}

	for _, pkg := range pkgs {
		for _, filename := range pkg.GoFiles {
			// Отсекаем не файлы проекта, например .cache/go-build
			if !strings.HasPrefix(filename, pPath) {
				continue
			}

			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
			if err != nil {
				return nil, fmt.Errorf("parse file %s: %w", filename, err)
			}

			pkgPath, _ := path.Split(filename)
			pkgPath = strings.TrimSuffix(pkgPath, "/")

			pkgInfo := model.Package{
				RelativePath: pkgPath,
				Inner:        isInnerPath(pkgPath),
			}

			if pkgInfo.Inner {
				pkgInfo.RelativePath = strings.ReplaceAll(pkgPath, pPath, coreModulePath)
				pkgInfo.InnerPath = strings.ReplaceAll(pkgPath, pPath, "")

				if pkgInfo.InnerPath == "" {
					pkgInfo.InnerPath = "/"
				}
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

				if importInfo.Inner {
					importInfo.InnerPath = strings.ReplaceAll(importInfo.RelativePath, coreModulePath, "")

					if importInfo.InnerPath == "" {
						importInfo.InnerPath = "/"
					}
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

	return pkgInfos, nil
}
