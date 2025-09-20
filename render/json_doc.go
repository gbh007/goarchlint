package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/gbh007/goarchlint/model"
)

func (r Render) DumpJSON(pkgs []model.Package) error {
	f, err := os.Create(path.Join(r.BasePath, "scheme.json"))
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	err = enc.Encode(pkgs)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}
