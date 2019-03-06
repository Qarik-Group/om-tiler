package configurator

import (
	"fmt"

	validator "gopkg.in/go-playground/validator.v9"
)

type Deployment struct {
	Opsman Opsman `yaml:"opsman" validate:"required,dive"`
	Tiles  []Tile `yaml:"tiles" validate:"required,dive"`
}

type Opsman struct {
	Target               string `yaml:"target" validate:"required"`
	Username             string `yaml:"username" validate:"required"`
	Password             string `yaml:"password" validate:"required"`
	DecryptionPassphrase string `yaml:"decryption_passphrase" validate:"required"`
	SkipSSLVerification  bool   `yaml:"skip_ssl_verification"`
	PivnetToken          string `yaml:"pivnet_token" validate:"required"`
}

type Tile struct {
	Product  Product                `yaml:"product" validate:"required,dive"`
	Features []string               `yaml:"features"`
	Optional []string               `yaml:"optional"`
	Resource []string               `yaml:"resource"`
	Network  string                 `yaml:"network"`
	Vars     map[string]interface{} `yaml:"vars"`
}

type Product struct {
	Name         string `yaml:"name" validate:"required"`
	Slug         string `yaml:"slug" validate:"required"`
	Version      string `yaml:"version" validate:"required"`
	Glob         string `yaml:"glob"`
	StemcellIaas string `yaml:"stemcell_iaas" validate:"required"`
}

func (d *Deployment) Validate() error {
	err := validator.New().Struct(d)
	if err != nil {
		return fmt.Errorf("Deployment has error(s):\n%+v\n", err)
	}
	return nil
}
