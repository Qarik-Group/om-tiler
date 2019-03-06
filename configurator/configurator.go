package configurator

import (
	"fmt"
	"log"
	"net/http"
)

type Configurator struct {
	client        OpsmanClient
	deployment    *Deployment
	logger        *log.Logger
	templateStore http.FileSystem
}

func NewConfigurator(d *Deployment,
	templateStore http.FileSystem,
	newOpsman func(Opsman, *log.Logger) (OpsmanClient, error),
	logger *log.Logger) (*Configurator, error) {

	err := d.Validate()
	if err != nil {
		return &Configurator{}, err
	}

	client, err := newOpsman(d.Opsman, logger)
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
