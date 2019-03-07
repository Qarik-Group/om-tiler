package configurator

import (
	"fmt"
	"log"
	"net/http"
)

type Configurator struct {
	client        OpsmanClient
	config        *Config
	logger        *log.Logger
	templateStore http.FileSystem
}

func NewConfigurator(c *Config, ts http.FileSystem,
	client OpsmanClient, l *log.Logger) (*Configurator, error) {

	err := c.Validate()
	if err != nil {
		return &Configurator{}, err
	}

	l.SetPrefix(fmt.Sprintf("%s[OM Configurator] ", l.Prefix()))

	configurator := Configurator{
		client: client, config: c,
		logger: l, templateStore: ts,
	}
	return &configurator, nil
}
