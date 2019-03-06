package configurator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/starkandwayne/om-configurator/config"
)

type DownloadProductArgs struct {
	OutputDirectory      string
	PivnetProductSlug    string
	PivnetProductVersion string
	PivnetProductGlob    string
	StemcellIaas         string
}

type UploadProductArgs struct {
	ProductFilePath      string
	PivnetProductVersion string
}

type UploadStemcellArgs struct {
	StemcellFilePath string
}

//go:generate counterfeiter . Opsman
type Opsman interface {
	ConfigureAuthentication() error
	DownloadProduct(DownloadProductArgs) error
	UploadProduct(UploadProductArgs) error
	UploadStemcell(UploadStemcellArgs) error
	ConfigureProduct() error
	ApplyChanges() error
}

func (c *Configurator) downloadAndUploadProduct(p config.Product) error {
	dir, err := ioutil.TempDir("", p.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	err = c.client.DownloadProduct(DownloadProductArgs{
		OutputDirectory:      dir,
		PivnetProductSlug:    p.Slug,
		PivnetProductVersion: p.Version,
		PivnetProductGlob:    p.Glob,
		StemcellIaas:         p.StemcellIaas,
	})
	if err != nil {
		return err
	}

	tile, err := findFileInDir(dir, "*pivotal")
	if err != nil {
		return err
	}

	err = c.client.UploadProduct(UploadProductArgs{
		ProductFilePath:      tile,
		PivnetProductVersion: p.Version,
	})
	if err != nil {
		return err
	}

	stemcell, err := findFileInDir(dir, "*tgz")
	if err != nil {
		return err
	}

	return c.client.UploadStemcell(UploadStemcellArgs{
		StemcellFilePath: stemcell,
	})
}

func findFileInDir(dir, glob string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, glob))
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		return "", fmt.Errorf("no file found for %s in %s", glob, dir)
	}
	return files[0], nil
}
