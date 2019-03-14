package mover_test

import (
	"log"

	. "github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/mover/moverfakes"
	"github.com/starkandwayne/om-tiler/tiler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mover", func() {
	It("Implements the tiler.Mover interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var mover tiler.Mover
		client := moverfakes.FakePivnetClient{}
		mover, err := NewMover(&client, "", logger)
		Expect(err).ToNot(HaveOccurred())
		logger.Println(mover) // use mover so it compiles
	})
})
