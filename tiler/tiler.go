package tiler

import (
	"fmt"
	"log"
)

type Tiler struct {
	client OpsmanClient
	logger *log.Logger
}

func NewTiler(c OpsmanClient, l *log.Logger) (*Tiler, error) {
	l.SetPrefix(fmt.Sprintf("%s[OM Tiler] ", l.Prefix()))
	return &Tiler{client: c, logger: l}, nil
}
