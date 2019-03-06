package configurator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
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

		err = c.configureProduct(tile)
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
	dir, err := ioutil.TempDir("", p.Slug)
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

	if err = c.client.UploadProduct(tile); err != nil {
		return err
	}

	stemcell, err := findFileInDir(dir, "*.tgz")
	if err != nil {
		return err
	}

	return c.client.UploadStemcell(stemcell)
}

func (c *Configurator) configureProduct(t Tile) error {
	ts := templateStore{base: t.Product.Slug, store: c.templateStore}

	templateFile, err := ts.lookup("", "product")
	if err != nil {
		return err
	}

	var opsFiles []io.Reader
	if err = ts.batchLookup("features", t.Features, &opsFiles, false); err != nil {
		return err
	}

	if err = ts.batchLookup("optional", t.Optional, &opsFiles, false); err != nil {
		return err
	}

	if err = ts.batchLookup("resource", t.Resource, &opsFiles, false); err != nil {
		return err
	}

	if t.Network != "" {
		network, err := ts.lookup("network", t.Network)
		if err != nil {
			return err
		}
		opsFiles = append(opsFiles, network)
	}

	var varsFiles []io.Reader
	if err != ts.batchLookup("", []string{
		"errand-vars", "product-default-vars", "resource-vars",
	}, &varsFiles, true) {
		return err
	}

	vars, err := yaml.Marshal(t.Vars)
	if err != nil {
		return err
	}

	varsFiles = append(varsFiles, bytes.NewReader(vars))

	ic := interpolateConfig{
		TemplateFile: templateFile,
		OpsFiles:     opsFiles,
		VarsFiles:    varsFiles,
	}

	tpl, err := ic.evaluate()
	if err != nil {
		return err
	}

	return c.client.ConfigureProduct(tpl)
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
