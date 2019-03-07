package opsman_test

import (
	"log"

	"github.com/starkandwayne/om-configurator/configurator"
	. "github.com/starkandwayne/om-configurator/opsman"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("Implements the configurator.Opsman interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var client configurator.OpsmanClient
		client, err := NewClient(Config{}, logger)
		Expect(err).ToNot(HaveOccurred())
		logger.Println(client) // use client so it compiles
	})
})
