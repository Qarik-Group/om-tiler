package opsman

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
	"github.com/starkandwayne/om-configurator/config"
	"github.com/starkandwayne/om-configurator/configurator"
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

func (c *Client) DownloadProduct(a *configurator.DownloadProductArgs) error {
	args := []string{
		fmt.Sprintf("--output-directory=%s", a.OutputDirectory),
		fmt.Sprintf("--pivnet-api-token=%s", c.config.PivnetToken),
		fmt.Sprintf("--pivnet-product-slug=%s", a.PivnetProductSlug),
		fmt.Sprintf("--pivnet-product-version=%s", a.PivnetProductVersion),
		fmt.Sprintf("--pivnet-product-glob=%s", a.PivnetProductGlob),
		fmt.Sprintf("--stemcell-iaas=%s", a.StemcellIaas),
		"--download-stemcell",
	}
	pivnetFactory := commands.DefaultPivnetFactory
	stower := commands.DefaultStow{}
	cmd := commands.NewDownloadProduct(os.Environ, c.log, c.log, pivnetFactory, stower)
	return cmd.Execute(args)
}

func (c *Client) UploadProduct(a *configurator.UploadProductArgs) error {
	args := []string{
		fmt.Sprintf("--product=%s", a.ProductFilePath),
		fmt.Sprintf("--product-version=%s", a.PivnetProductVersion),
		fmt.Sprintf("--polling-interval=%s", pollingIntervalSec),
	}
	form := formcontent.NewForm()
	cmd := commands.NewUploadProduct(form, metadataExtractor, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) UploadStemcell(a *configurator.UploadStemcellArgs) error {
	args := []string{
		fmt.Sprintf("--stemcell=%s", a.StemcellFilePath),
		"--floating",
	}
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, c.api, c.log)
	return uscmd.Execute(args)
}

func (c *Client) ConfigureProduct() error {

}

func (c *Client) ApplyChanges() error {

}
