package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"examples/templates"

	"github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/opsman"
	"github.com/starkandwayne/om-tiler/pivnet"
	"github.com/starkandwayne/om-tiler/tiler"
)

func main() {

	logger := log.New(os.Stdout, "", 0)
	workDir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	cacheDir := filepath.Join(workDir, "cache")
	varsStore := filepath.Join(workDir, "creds.yml")
	mover, err := mover.NewMover(
		pivnet.NewClient(pivnet.Config{
			Token: os.Getenv("PIVNET_TOKEN")}, logger),
		cacheDir,
		logger,
	)
	if err != nil {
		logger.Fatal(err)
	}
	om, err := opsman.NewClient(opsman.Config{
		Target:               os.Getenv("OPSMAN_TARGET"),
		Username:             "admin",
		Password:             os.Getenv("OPSMAN_PASSWORD"),
		DecryptionPassphrase: os.Getenv("OPSMAN_PASSWORD"),
		SkipSSLVerification:  true,
	}, logger)
	if err != nil {
		logger.Fatal(err)
	}

	t := tiler.NewTiler(om, mover, logger)

	vars := map[string]interface{}{
		"domain": "pcf.example.com",
	}
	pattern, err := templates.GetPattern(vars, varsStore, true)
	if err != nil {
		logger.Fatal(err)
	}

	skipApplyChanges := false
	err = t.Build(context.Background(), pattern, skipApplyChanges)
	if err != nil {
		logger.Fatal(err)
	}
}
