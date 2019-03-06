package configurator

import (
	"fmt"
	"log"
	"net/http"

	"github.com/starkandwayne/om-configurator/config"
)

type Configurator struct {
	client        Opsman
	deployment    *config.Deployment
	logger        *log.Logger
	templateStore *http.FileSystem
}

func NewConfigurator(d *config.Deployment,
	templateStore *http.FileSystem,
	newOpsman func(*config.Opsman, *log.Logger) (Opsman, error),
	logger *log.Logger) (*Configurator, error) {

	client, err := newOpsman(&d.Opsman, logger)
	if err != nil {
		return &Configurator{}, err
	}

	logger.SetPrefix(fmt.Sprintf("%s[OM Configurator] ", logger.Prefix()))

	configurator := Configurator{
		client:        client,
		deployment:    d,
		logger:        logger,
		templateStore: templateStore,
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
