package opsman

import (
	"fmt"
	"log"

	gopivnet "github.com/pivotal-cf/go-pivnet"
	"github.com/pivotal-cf/go-pivnet/logshim"
	"github.com/starkandwayne/om-configurator/configurator"
)

type PivnetEulaAccepter struct {
	Endpoint  string
	Token     string
	UserAgent string
	Logger    *log.Logger
}

func (c *PivnetEulaAccepter) Accept(a configurator.DownloadProductArgs) error {
	client := c.getClient()

	ok, err := client.Auth.Check()
	if !ok {
		return fmt.Errorf("authorizing pivnet credentials: %v", err)
	}

	releaseID, err := releaseIDForVersion(&client, a.PivnetProductSlug, a.PivnetProductVersion)
	if err != nil {
		return err
	}
	return client.EULA.Accept(a.PivnetProductSlug, releaseID)
}

func releaseIDForVersion(c *gopivnet.Client, slug, version string) (int, error) {
	releases, err := c.Releases.List(slug)
	if err != nil {
		return 0, err
	}

	for _, r := range releases {
		if r.Version == version {
			return r.ID, nil
		}
	}

	return 0, fmt.Errorf(
		"release not found for %s with version: '%s'", slug, version,
	)
}

func (c *PivnetEulaAccepter) getClient() gopivnet.Client {
	endpoint := c.Endpoint
	if c.Endpoint == "" {
		endpoint = gopivnet.DefaultHost
	}

	return gopivnet.NewClient(gopivnet.ClientConfig{
		Host:      endpoint,
		Token:     c.Token,
		UserAgent: c.UserAgent,
	}, logshim.NewLogShim(c.Logger, c.Logger, false))
}
