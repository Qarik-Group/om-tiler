package tiler

import (
	"fmt"
	"log"
	"net/http"
)

type Tiler struct {
	client        OpsmanClient
	logger        *log.Logger
	templateStore http.FileSystem
}

func NewTiler(ts http.FileSystem,
	client OpsmanClient, l *log.Logger) (*Tiler, error) {

	l.SetPrefix(fmt.Sprintf("%s[OM Tiler] ", l.Prefix()))

	tiler := Tiler{
		client: client, logger: l, templateStore: ts,
	}
	return &tiler, nil
}
