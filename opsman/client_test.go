package opsman_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	. "github.com/starkandwayne/om-tiler/opsman"
	"github.com/starkandwayne/om-tiler/pattern"
	"github.com/starkandwayne/om-tiler/tiler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	It("Implements the tiler.Opsman interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var client tiler.OpsmanClient
		client, err := NewClient(Config{}, logger)
		Expect(err).ToNot(HaveOccurred())
		logger.Println(client) // use client so it compiles
	})

	var (
		server     *ghttp.Server
		client     *Client
		apiAddress string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL()
		header := http.Header{}
		header.Add("Content-Type", "application/json")
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/uaa/oauth/token"),
				ghttp.RespondWith(http.StatusOK, `
                                      {"access_token":"token","token_type":"bearer","expires_in":"3600"}
                                `, header),
			),
		)

		logger := log.New(GinkgoWriter, "", 0)
		var err error
		client, err = NewClient(Config{
			Target:               apiAddress,
			Username:             "admin",
			Password:             "password",
			DecryptionPassphrase: "decrypt",
			SkipSSLVerification:  true,
		}, logger)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("FilesUploaded", func() {
		var (
			products []string
			product  func(string)
			tile     pattern.Tile
		)

		BeforeEach(func() {
			products = []string{}
			product = func(p string) {
				products = append(products, p)
			}

			tile = pattern.Tile{
				Name:    "cf",
				Version: "2.4.4-build.2",
				Product: pattern.PivnetFile{
					Slug:    "elastic-runtime",
					Version: "2.4.4",
					Glob:    "srt-2.4.4-build.2.pivotal",
				},
				Stemcell: pattern.PivnetFile{
					Slug:    "stemcells-ubuntu-xenial",
					Version: "170.39",
					Glob:    "*vsphere*.tgz",
				},
			}
		})

		JustBeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v0/stemcell_assignments"),
					ghttp.RespondWith(http.StatusOK,
						fmt.Sprintf(`{"products":[%s]}`,
							strings.Join(products, ","),
						)),
				),
			)
		})

		Context("When product has already been uploaded", func() {
			BeforeEach(func() {
				product(`{"identifier":"cf","staged_product_version":"2.4.4-build.2",
                                          "available_stemcell_versions":["170.39"]}`)
			})

			It("accepts the EULA for a given release and product Version", func() {
				ok, err := client.FilesUploaded(context.Background(), tile)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(Equal(true))
			})
		})

		Context("When the stemcell has not been uploaded", func() {
			BeforeEach(func() {
				product(`{"type":"cf","product_version":"2.4.4-build.2",
                                          "available_stemcell_versions":[]}`)
			})

			It("accepts the EULA for a given release and product Version", func() {
				ok, err := client.FilesUploaded(context.Background(), tile)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(Equal(false))
			})
		})

		Context("When the product not been uploaded", func() {
			BeforeEach(func() {
				product(`{"type":"foo","product_version":"2.4.4-build.2",
                                          "available_stemcell_versions":["170.39"]}`)
			})

			It("accepts the EULA for a given release and product Version", func() {
				ok, err := client.FilesUploaded(context.Background(), tile)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(Equal(false))
			})
		})

	})

})
