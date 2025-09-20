package model

import "slices"

type Import struct {
	Inner        bool   `json:"inner,omitempty"`
	RelativePath string `json:"relative_path,omitempty"`
	InnerPath    string `json:"inner_path,omitempty"`
	Files        []File `json:"files,omitempty"`
}

func (i Import) Clone() Import {
	i.Files = slices.Clone(i.Files)

	return i
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
