package pattern

import (
	"fmt"
	"net/http"

	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

type Pattern struct {
	Director Director `yaml"director validate:"required,dive"`
	Tiles    []Tile   `yaml:"tiles" validate:"required,dive"`
}

func NewPattern(t Template) (p Pattern, err error) {
	db, err := t.Evaluate(false)
	if err != nil {
		return Pattern{}, err
	}

	if err = yaml.Unmarshal(db, &p); err != nil {
		return Pattern{}, err
	}

	mergeVars(p.Director.Vars, t.Vars)
	p.Director.Store = t.Store

	for i, _ := range p.Tiles {
		mergeVars(p.Tiles[i].Vars, t.Vars)
		p.Tiles[i].Store = t.Store
	}

	return p, err
}

func mergeVars(target map[string]interface{}, source map[string]interface{}) {
	if target == nil {
		target = source
		return
	}
	for k, v := range source {
		if _, ok := target[k]; !ok {
			target[k] = v
		}
	}
}

func (p *Pattern) Validate(expectAllKeys bool) error {
	err := validator.New().Struct(p)
	if err != nil {
		return fmt.Errorf("pattern.Pattern has error(s):\n%+v\n", err)
	}

	_, err = p.Director.ToTemplate().Evaluate(expectAllKeys)
	if err != nil {
		return fmt.Errorf("Director interpolation error(s):\n%+v\n", err)
	}

	for _, tile := range p.Tiles {
		_, err = tile.ToTemplate().Evaluate(expectAllKeys)
		if err != nil {
			return fmt.Errorf("Tile %s interpolation error(s):\n%+v\n", err, tile.Name)
		}
	}

	return nil
}

type Template struct {
	Manifest  string                 `yaml:"manifest"`
	OpsFiles  []string               `yaml:"ops_files"`
	VarsFiles []string               `yaml:"vars_files"`
	Vars      map[string]interface{} `yaml:"vars"`
	Store     http.FileSystem
}

type Director Template

func (d *Director) ToTemplate() *Template {
	return &Template{
		Manifest:  d.Manifest,
		OpsFiles:  d.OpsFiles,
		VarsFiles: d.VarsFiles,
		Vars:      d.Vars,
		Store:     d.Store,
	}
}

type Tile struct {
	Name       string     `yaml:"name" validate:"required"`
	Version    string     `yaml:"version" validate:"required"`
	PivnetMeta PivnetMeta `yaml:"pivnet" validate:"required,dive"`
	Template   `yaml:",inline"`
}

func (t *Tile) ToTemplate() *Template {
	return &Template{
		Manifest:  t.Manifest,
		OpsFiles:  t.OpsFiles,
		VarsFiles: t.VarsFiles,
		Vars:      t.Vars,
		Store:     t.Store,
	}
}

type PivnetMeta struct {
	Slug         string `yaml:"slug" validate:"required"`
	Version      string `yaml:"version" validate:"required"`
	Glob         string `yaml:"glob"`
	StemcellIaas string `yaml:"stemcell_iaas" validate:"required"`
}
