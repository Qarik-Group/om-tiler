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

type Config struct {
	Target               string
	Username             string
	Password             string
	DecryptionPassphrase string
	SkipSSLVerification  bool
	PivnetToken          string
	PivnetUserAgent      string
}

type Client struct {
	api    api.Api
	log    *log.Logger
	config Config
}

const (
	connectTimeout     = time.Duration(5) * time.Second
	requestTimeout     = time.Duration(1800) * time.Second
	pollingIntervalSec = "10"
)

func NewClient(c Config, logger *log.Logger) (*Client, error) {
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

func (c *Client) DownloadProduct(a configurator.DownloadProductArgs) error {
	if c.config.PivnetUserAgent != "" {
		piv := PivnetEulaAccepter{
			Token:     c.config.PivnetToken,
			UserAgent: c.config.PivnetUserAgent,
			Logger:    c.log,
		}
		err := piv.Accept(a)
		if err != nil {
			return err
		}
	}

	args := []string{
		fmt.Sprintf("--output-directory=%s", a.OutputDirectory),
		fmt.Sprintf("--pivnet-api-token=%s", c.config.PivnetToken),
		fmt.Sprintf("--pivnet-product-slug=%s", a.PivnetProductSlug),
		fmt.Sprintf("--pivnet-file-glob=%s", a.PivnetProductGlob),
		fmt.Sprintf("--product-version=%s", a.PivnetProductVersion),
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

func (c *Client) StageProduct(a configurator.StageProductArgs) error {
	args := []string{
		fmt.Sprintf("--product-name=%s", a.ProductName),
		fmt.Sprintf("--product-version=%s", a.ProductVersion),
	}
	cmd := commands.NewStageProduct(c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) ConfigureProduct(config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return err
	}
	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureProduct(
		os.Environ, c.api, c.config.Target, c.log)
	return cmd.Execute(args)
}

func (c *Client) ConfigureDirector(config []byte) error {
	configFile, err := tmpConfigFile(config)
	if err != nil {
		return err
	}
	args := []string{
		fmt.Sprintf("--config=%s", configFile),
	}
	cmd := commands.NewConfigureDirector(os.Environ, c.api, c.log)
	return cmd.Execute(args)
}

func (c *Client) ApplyChanges() error {
	return nil
}

func tmpConfigFile(config []byte) (string, error) {
	configFile, err := ioutil.TempFile("", "config")
	if err != nil {
		return "", err
	}
	defer os.Remove(configFile.Name())

	if _, err = configFile.Write(config); err != nil {
		return "", err
	}

	if err = configFile.Close(); err != nil {
		return "", err
	}

	return configFile.Name(), nil
}
