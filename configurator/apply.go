package configurator

import (
	"fmt"
	"log"

	"github.com/starkandwayne/om-configurator/config"
	"github.com/starkandwayne/om-configurator/opsman"
)

type Configurator struct {
	client     *opsman.Client
	deployment *config.Deployment
	logger     *log.Logger
}

func NewConfigurator(d *config.Deployment, logger *log.Logger) (*Configurator, error) {
	client, err := opsman.NewClient(d.Opsman, logger)
	if err != nil {
		return &Configurator{}, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM Configurator] ", logger.Prefix()))

	configurator := Configurator{
		client:     client,
		deployment: d,
		logger:     logger,
	}
	return &configurator, nil
}

func (c *Configurator) Apply() error {
	err := c.deployment.Validate()
	if err != nil {
		return err
	}

	err = c.client.ConfigureAuthentication()
	if err != nil {
		return err
	}

	for _, tile := range c.deployment.Tiles {
		err = c.client.UploadProduct(tile.Product)
		if err != nil {
			return err
		}

		err = c.client.ConfigureProduct(tile)
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
