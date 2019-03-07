package opsman

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gosuri/uilive"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/pivotal-cf/om/api"
	"github.com/pivotal-cf/om/commands"
	"github.com/pivotal-cf/om/extractor"
	"github.com/pivotal-cf/om/formcontent"
	"github.com/pivotal-cf/om/network"
	"github.com/pivotal-cf/om/progress"
	"github.com/starkandwayne/om-configurator/configurator"
)

type Client struct {
	api    api.Api
	log    *log.Logger
	config configurator.Opsman
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingIntervalSec = time.Duration(10) * time.Second
)

func NewClient(c configurator.Opsman, logger *log.Logger) (*Client, error) {
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
	cmd := commands.NewConfigureAuthentication(c.api, c.log)
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
	pivnetLogWriter := logshim.NewLogShim(c.log, c.log, false)
	cmd := commands.NewDownloadProduct(os.Environ, pivnetLogWriter, os.Stdout, pivnetFactory, stower)
	return cmd.Execute(args)
}

func (c *Client) UploadProduct(productFile string) error {
	args := []string{
		fmt.Sprintf("--product=%s", productFile),
		fmt.Sprintf("--polling-interval=%s", pollingIntervalSec),
	}
	form := formcontent.NewForm()
	metadataExtractor := extractor.MetadataExtractor{}
	cmd := commands.NewUploadProduct(form, metadataExtractor, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) UploadStemcell(stemcell string) error {
	args := []string{
		fmt.Sprintf("--stemcell=%s", stemcell),
		"--floating",
	}
	form := formcontent.NewForm()
	cmd := commands.NewUploadStemcell(form, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) ConfigureProduct(config []byte) error {
	configFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return err
	}
	defer os.Remove(configFile.Name())

	if _, err = configFile.Write(config); err != nil {
		return err
	}

	if err = configFile.Close(); err != nil {
		return err
	}

	args := []string{
		fmt.Sprintf("--config=%s", configFile.Name()),
	}
	cmd := commands.NewConfigureProduct(
		os.Environ, c.api, c.config.Target, c.log)
	return cmd.Execute(args)
}

func (c *Client) ApplyChanges() error {
	return nil
}
