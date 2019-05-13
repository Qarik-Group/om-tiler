package templates

//go:generate go run -tags=dev generate.go

import (
	"github.com/starkandwayne/om-tiler/pattern"
)

func GetPattern(vars map[string]interface{}, varsStore string, expectAllKeys bool) (pattern.Pattern, error) {
	return pattern.NewPattern(pattern.Template{
		Store:    Templates,
		Manifest: "pattern.yml",
		Vars:     vars,
	}, varsStore, expectAllKeys)
}
