package opsman

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
	"github.com/starkandwayne/om-configurator/config"
)

type Client struct {
	api    api.Api
	log    *log.Logger
	config config.Opsman
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingIntervalSec = time.Duration(10) * time.Second
)

func NewClient(c config.Opsman, logger *log.Logger) (*Client, error) {
	oauthClient, err := network.NewOAuthClient(
		c.Target, c.Username, c.Password, "", "",
		c.SkipSSLVerification, true,
		requestTimeout, connectTimeout,
	)
	if err != nil {
		return &Client{}, err
	}

	unauthenticatedClient := network.NewUnauthenticatedClient(
		c.Target, c.SkipSSLVerification,
		requestTimeout, connectTimeout,
	)

	logger.SetPrefix(fmt.Sprintf("%s[OM] ", logger.Prefix()))

	live := uilive.New()
	live.Out = os.Stderr

	client := Client{
		api: api.New(api.ApiInput{
			Client:         oauthClient,
			UnauthedClient: unauthenticatedClient,
			ProgressClient: network.NewProgressClient(
				oauthClient, progress.NewBar(), live),
			UnauthedProgressClient: network.NewProgressClient(
				unauthenticatedClient, progress.NewBar(), live),
			Logger: logger,
		}),
		log:    logger,
		config: c,
	}

	return &client, nil
}

func (c *Client) ConfigureAuthentication() error {
	args := []string{
		fmt.Sprintf("--username=%s", c.config.Username),
		fmt.Sprintf("--password=%s", c.config.Password),
		fmt.Sprintf("--decryption-passphrase=%s", c.config.DecryptionPassphrase),
	}
	cmd := commands.NewConfigureAuthentication(c.api, c.log, c.log)
	return cmd.Execute(args)
}

func (c *Client) UploadProduct(p *config.Product) error {
	dir, err := ioutil.TempDir("", p.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	dargs := []string{
		fmt.Sprintf("--output-directory=%s", dir),
		fmt.Sprintf("--pivnet-api-token=%s", c.config.PivnetToken),
		fmt.Sprintf("--pivnet-product-slug=%s", p.Slug),
		fmt.Sprintf("--pivnet-product-version=%s", p.Version),
		fmt.Sprintf("--pivnet-product-glob=%s", p.Glob),
		fmt.Sprintf("--stemcell-iaas=%s", p.StemcellIaas),
		"--download-stemcell",
	}
	pivnetFactory := commands.DefaultPivnetFactory
	stower := commands.DefaultStow{}
	dcmd := commands.NewDownloadProduct(os.Environ, c.log, c.log, pivnetFactory, stower)
	err := dcmd.Execute(dargs)
	if err != nil {
		return err
	}

	pFile, err := filepath.Glob(filepath.Join(dir, "*.pivotal"))
	if err != nil {
		return err
	}
	if len(pFile) != 1 {
		return error.New(fmt.Sprintf("No tile found for %s in %s", p.Name, dir))
	}

	upargs := []string{
		fmt.Sprintf("--product=%s", pFile[0]),
		fmt.Sprintf("--product-version=%s", p.Version),
		fmt.Sprintf("--polling-interval=%s", pollingIntervalSec),
	}
	pform := formcontent.NewForm()
	upcmd := commands.NewUploadProduct(pform, metadataExtractor, c.api, c.log)
	err = upcmd.Execute(upargs)
	if err != nil {
		return err
	}

	sFile, err := filepath.Glob(filepath.Join(dir, "*.tgz"))
	if err != nil {
		return err
	}
	if len(sFile) != 1 {
		return error.New(fmt.Sprintf("No stemcell found for %s in %s", p.Name, dir))
	}

	usargs := []string{
		fmt.Sprintf("--stemcell=%s", sFile[0]),
		"--floating",
	}
	sform := formcontent.NewForm()
	uscmd := commands.NewUploadStemcell(sform, c.api, c.log)
	err = uscmd.Execute(usargs)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ConfigureProduct(*config.Tile) error {

}

func (c *Client) ApplyChanges() error {

}
