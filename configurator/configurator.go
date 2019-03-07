package configurator

import (
	"fmt"
	"log"
	"net/http"
)

type Configurator struct {
	client        OpsmanClient
	logger        *log.Logger
	templateStore http.FileSystem
}

func NewConfigurator(ts http.FileSystem,
	client OpsmanClient, l *log.Logger) (*Configurator, error) {

	l.SetPrefix(fmt.Sprintf("%s[OM Configurator] ", l.Prefix()))

	configurator := Configurator{
		client: client, logger: l, templateStore: ts,
	}
	return &configurator, nil
}
