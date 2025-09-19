package model

type File struct {
	Path   string `json:"path,omitempty"`
	Pos    int    `json:"pos,omitempty"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

func CompactFiles(files ...[]File) []File {
	tmp := make(map[File]struct{}, len(files))

	for _, arr := range files {
		for _, e := range arr {
			tmp[e] = struct{}{}
		}
	}

	result := make([]File, 0, len(tmp))

	for k := range tmp {
		result = append(result, k)
	}

	return result
}
