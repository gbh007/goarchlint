package model

type Package struct {
	Name         string   `json:"name,omitempty"`
	Inner        bool     `json:"inner,omitempty"`
	RelativePath string   `json:"relative_path,omitempty"`
	InnerPath    string   `json:"inner_path,omitempty"`
	IsMain       bool     `json:"is_main,omitempty"`
	Imports      []Import `json:"imports,omitempty"`
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
