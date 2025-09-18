package model

type Package struct {
	Name         string
	Inner        bool
	RelativePath string

	IsMain  bool
	Imports []Import
}

type Import struct {
	Inner        bool
	RelativePath string

	Files []File
}

type File struct {
	Path   string
	Pos    int
	Line   int
	Column int
}

func CompactPackages(s ...[]Package) []Package {
	tmp := make(map[string]Package)

	for _, imps := range s {
		for _, imp := range imps {
			other, ok := tmp[imp.RelativePath]
			if !ok {
				tmp[imp.RelativePath] = imp

				continue
			}

			other.Imports = CompactImports(other.Imports, imp.Imports)
			tmp[imp.RelativePath] = other
		}
	}

	result := make([]Package, 0, len(tmp))

	for _, k := range tmp {
		result = append(result, k)
	}

	return result
}

func CompactImports(s ...[]Import) []Import {
	tmp := make(map[string]Import)

	for _, imps := range s {
		for _, imp := range imps {
			other, ok := tmp[imp.RelativePath]
			if !ok {
				tmp[imp.RelativePath] = imp

				continue
			}

			other.Files = CompactFiles(other.Files, imp.Files)
			tmp[imp.RelativePath] = other
		}
	}

	result := make([]Import, 0, len(tmp))

	for _, k := range tmp {
		result = append(result, k)
	}

	return result
}

func CompactFiles(f ...[]File) []File {
	return Unique(f...)
}

func Unique[V comparable](s ...[]V) []V {
	tmp := make(map[V]struct{}, len(s))

	for _, arr := range s {
		for _, e := range arr {
			tmp[e] = struct{}{}
		}
	}

	result := make([]V, 0, len(tmp))

	for k := range tmp {
		result = append(result, k)
	}

	return result
}
