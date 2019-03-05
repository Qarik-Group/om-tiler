package config

type Deployment struct {
	Opsman Opsman `yaml:"opsman"`
	Tiles  []Tile `yaml:"tiles"`
}

type Opsman struct {
	Target               string `yaml:"target"`
	Username             string `yaml:"username"`
	Password             string `yaml:"password"`
	DecryptionPassphrase string `yaml:"decryption_passphrase"`
	SkipSSLVerification  bool   `yaml:"skip_ssl_verification"`
	PivnetToken          string `yaml:"pivnet_token"`
}

type Tile struct {
	Product  Product                `yaml:"product"`
	Features []string               `yaml:"features"`
	Network  string                 `yaml:"network"`
	Optional []string               `yaml:"optional"`
	Resource []string               `yaml:"resource"`
	Vars     map[string]interface{} `yaml:"vars"`
}

type Product struct {
	Name         string `yaml:"name"`
	Slug         string `yaml:"slug"`
	Version      string `yaml:"version"`
	Glob         string `yaml:"glob"`
	StemcellIaas string `yaml:"stemcell_iaas"`
}

func (d *Deployment) Validate() error {
	return nil
}
