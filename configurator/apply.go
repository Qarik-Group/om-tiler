package configurator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (c *Configurator) Apply() error {
	err := c.client.ConfigureAuthentication()
	if err != nil {
		return err
	}

	for _, tile := range c.deployment.Tiles {
		err = c.downloadAndUploadProduct(tile.Product)
		if err != nil {
			return err
		}

		err = c.client.ConfigureProduct()
		if err != nil {
			return err
		}
	}

	err = c.client.ApplyChanges()
	if err != nil {
		return err
	}

	return nil
}

func (c *Configurator) downloadAndUploadProduct(p Product) error {
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

	tile, err := findFileInDir(dir, "*.pivotal")
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

	stemcell, err := findFileInDir(dir, "*.tgz")
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
