package opsman_test

import (
	"log"

	. "github.com/starkandwayne/om-tiler/opsman"
	"github.com/starkandwayne/om-tiler/tiler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("Implements the tiler.Opsman interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var client tiler.OpsmanClient
		client, err := NewClient(Config{}, logger)
		Expect(err).ToNot(HaveOccurred())
		logger.Println(client) // use client so it compiles
	})
})
