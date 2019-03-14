package pivnet

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/starkandwayne/om-tiler/pattern"
)

const (
	retryAttempts = 5 // How many times to retry downloading a tile from PivNet
	retryDelay    = 5 // How long wait in between download retries
)

type Config struct {
	Host       string
	Token      string
	UserAgent  string
	Logger     *log.Logger
	AcceptEULA bool
}

type Client struct {
	logger     *log.Logger
	client     gopivnet.Client
	acceptEULA bool
}

func NewClient(c Config, logger *log.Logger) *Client {
	host := c.Host
	if c.Host == "" {
		host = gopivnet.DefaultHost
	}

	client := gopivnet.NewClient(gopivnet.ClientConfig{
		Host:      host,
		Token:     c.Token,
		UserAgent: c.UserAgent,
	}, logshim.NewLogShim(c.Logger, c.Logger, false))
	return &Client{client: client, logger: logger,
		acceptEULA: c.AcceptEULA}
}

func (c *Client) DownloadFile(f pattern.PivnetFile, path string) (file *os.File, err error) {
	if c.acceptEULA {
		if err = c.AcceptEULA(f); err != nil {
			return
		}
	}
	for i := 0; i < retryAttempts; i++ {
		file, err = c.downloadFile(f, path)

		// Success or recoverable error
		if err == nil || err != io.ErrUnexpectedEOF {
			return
		}

		c.logger.Printf("download tile failed, retrying in %d seconds", retryDelay)
		time.Sleep(time.Duration(retryDelay) * time.Second)
	}

	return nil, fmt.Errorf("download tile failed after %d attempts", retryAttempts)
}

func (c *Client) GetEULA(f pattern.PivnetFile) (string, error) {
	release, err := c.lookupRelease(f)
	if err != nil {
		return "", err
	}

	eula, err := c.client.EULA.Get(release.EULA.Slug)
	if err != nil {
		return "", err
	}

	return eula.Content, nil
}

func (c *Client) AcceptEULA(f pattern.PivnetFile) error {
	release, err := c.lookupRelease(f)
	if err != nil {
		return err
	}

	return c.client.EULA.Accept(f.Slug, release.ID)
}

func (c *Client) downloadFile(f pattern.PivnetFile, path string) (file *os.File, err error) {
	if path == "" {
		file, err = ioutil.TempFile("", "tile")
	} else {
		file, err = os.Create(path)
	}
	if err != nil {
		return nil, err
	}

	// Delete the file if we're returning an error
	defer func() {
		if err != nil {
			os.Remove(file.Name())
		}
	}()

	productFile, release, err := c.lookupProductFile(f)
	if err != nil {
		return nil, err
	}

	return file, c.client.ProductFiles.DownloadForRelease(file, f.Slug, release.ID, productFile.ID, os.Stdout)
}

func normalizeReleaseVersion(v string) (string, error) {
	version, err := goversion.NewVersion(v)
	if err != nil {
		return "", fmt.Errorf("Unable to parse version %s: %s", v, err)
	}
	return strings.Trim(strings.Replace(fmt.Sprint(version.Segments()), " ", ".", -1), "[]"), nil
}

func (c *Client) lookupRelease(f pattern.PivnetFile) (gopivnet.Release, error) {
	version, err := normalizeReleaseVersion(f.Version)
	if err != nil {
		return gopivnet.Release{}, err
	}

	releases, err := c.client.Releases.List(f.Slug)
	if err != nil {
		return gopivnet.Release{}, err
	}

	for _, r := range releases {
		rv, _ := normalizeReleaseVersion(r.Version)
		if rv == version {
			return r, nil
		}
	}

	return gopivnet.Release{}, fmt.Errorf(
		"release not found for %s with version: '%s'", f.Slug, version,
	)
}

func (c *Client) lookupProductFile(f pattern.PivnetFile) (gopivnet.ProductFile, gopivnet.Release, error) {
	release, err := c.lookupRelease(f)
	files, err := c.client.ProductFiles.ListForRelease(f.Slug, release.ID)
	if err != nil {
		return gopivnet.ProductFile{}, gopivnet.Release{}, err
	}

	for _, file := range files {
		if file.FileVersion == f.Version {
			return file, release, err
		}
	}

	return gopivnet.ProductFile{}, gopivnet.Release{}, fmt.Errorf(
		"file not found for %s/%s with version: '%s'", f.Slug, release.Version, f.Version,
	)
}
