package configurator

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

type Config struct {
	Target               string `yaml:"target" validate:"required"`
	Username             string `yaml:"username" validate:"required"`
	Password             string `yaml:"password" validate:"required"`
	DecryptionPassphrase string `yaml:"decryption_passphrase" validate:"required"`
	SkipSSLVerification  bool   `yaml:"skip_ssl_verification"`
	PivnetToken          string `yaml:"pivnet_token" validate:"required"`
}

func (c *Config) Validate() error {
	return validate("Config", c)
}

type Deployment struct {
	Director Director `yaml"director validate:"required,dive"`
	Tiles    []Tile   `yaml:"tiles" validate:"required,dive"`
}

func (d *Deployment) Validate() error {
	return validate("Deployment", d)
}

type Template struct {
	Manifest  string                 `yaml:"manifest"`
	OpsFiles  []string               `yaml:"ops_files"`
	VarsFiles []string               `yaml:"vars_files"`
	Vars      map[string]interface{} `yaml:"vars"`
}

type Director Template

func (d *Director) ToTemplate() Template {
	return Template{
		Manifest:  d.Manifest,
		OpsFiles:  d.OpsFiles,
		VarsFiles: d.VarsFiles,
		Vars:      d.Vars,
	}
}

type Tile struct {
	PivnetMeta PivnetMeta `yaml:"pivnet" validate:"required,dive"`
	OpsmanMeta OpsmanMeta `yaml:"opsman" validate:"required,dive"`
	Template   `yaml:",inline"`
}

func (t *Tile) ToTemplate() Template {
	return Template{
		Manifest:  t.Manifest,
		OpsFiles:  t.OpsFiles,
		VarsFiles: t.VarsFiles,
		Vars:      t.Vars,
	}
}

type PivnetMeta struct {
	Slug         string `yaml:"slug" validate:"required"`
	Version      string `yaml:"version" validate:"required"`
	Glob         string `yaml:"glob"`
	StemcellIaas string `yaml:"stemcell_iaas" validate:"required"`
}

type OpsmanMeta struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required"`
}

func validate(name string, s interface{}) error {
	err := validator.New().Struct(s)
	if err != nil {
		return fmt.Errorf("%s has error(s):\n%+v\n", name, err)
	}
	return nil
}
