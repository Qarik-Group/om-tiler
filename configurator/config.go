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
	Tiles []Tile `yaml:"tiles" validate:"required,dive"`
}

func (d *Deployment) Validate() error {
	return validate("Deployment", d)
}

type Tile struct {
	PivnetMeta PivnetMeta             `yaml:"pivnet" validate:"required,dive"`
	OpsmanMeta OpsmanMeta             `yaml:"opsman" validate:"required,dive"`
	Features   []string               `yaml:"features"`
	Optional   []string               `yaml:"optional"`
	Resource   []string               `yaml:"resource"`
	Network    string                 `yaml:"network"`
	Vars       map[string]interface{} `yaml:"vars"`
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
