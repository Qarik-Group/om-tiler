package tiler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/starkandwayne/om-tiler/pattern"
)

func (c *Tiler) Apply(t pattern.Template) error {
	p, err := pattern.NewPattern(t)
	if err != nil {
		return err
	}

	if err = p.Validate(true); err != nil {
		return err
	}

	err = c.client.ConfigureAuthentication()
	if err != nil {
		return err
	}

	err = c.configureDirector(p.Director)
	if err != nil {
		return err
	}

	for _, tile := range p.Tiles {
		err = c.downloadAndUploadProduct(tile.PivnetMeta)
		if err != nil {
			return err
		}

		err = c.client.StageProduct(StageProductArgs{
			ProductName:    tile.Name,
			ProductVersion: tile.Version,
		})
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

func (c *Tiler) downloadAndUploadProduct(p pattern.PivnetMeta) error {
	dir, err := ioutil.TempDir("", p.Slug)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	glob := p.Glob
	if glob == "" {
		glob = "*.pivotal"
	}

	err = c.client.DownloadProduct(DownloadProductArgs{
		OutputDirectory:      dir,
		PivnetProductSlug:    p.Slug,
		PivnetProductVersion: p.Version,
		PivnetProductGlob:    glob,
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

func (c *Tiler) configureProduct(t pattern.Tile) error {
	tpl, err := t.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureProduct(tpl)
}

func (c *Tiler) configureDirector(d pattern.Director) error {
	tpl, err := d.ToTemplate().Evaluate(true)
	if err != nil {
		return err
	}

	return c.client.ConfigureDirector(tpl)
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
