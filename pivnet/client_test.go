package pivnet_test

import (
	"fmt"
	"log"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/starkandwayne/om-tiler/mover"
	"github.com/starkandwayne/om-tiler/pattern"
	. "github.com/starkandwayne/om-tiler/pivnet"
)

var _ = Describe("Client", func() {
	It("Implements the mover.Pivnet interface", func() {
		logger := log.New(GinkgoWriter, "", 0)
		var client mover.PivnetClient
		client = NewClient(Config{}, logger)
		logger.Println(client) // use client so it compiles
	})

	var (
		server     *ghttp.Server
		client     *Client
		token      string
		apiAddress string
		apiPrefix  string
		userAgent  string
		acceptEULA bool
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL()
		apiPrefix = "/api/v2"
		token = "my-auth-token"
		userAgent = "pivnet-resource/0.1.0 (some-url)"
		logger := log.New(GinkgoWriter, "", 0)
		client = NewClient(Config{
			Host:       apiAddress,
			Token:      token,
			UserAgent:  userAgent,
			AcceptEULA: acceptEULA,
		}, logger)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("AcceptEULA", func() {
		var (
			releaseVersion    string
			productSlug       string
			releaseID         int
			releasesURL       string
			EULAAcceptanceURL string
		)

		BeforeEach(func() {
			productSlug = "banana-slug"
			releaseVersion = "3.2"
			releaseID = 42
			releasesURL = fmt.Sprintf(apiPrefix+"/products/%s/releases", productSlug)
			EULAAcceptanceURL = fmt.Sprintf(apiPrefix+"/products/%s/releases/%d/pivnet_resource_eula_acceptance", productSlug, releaseID)
		})

		It("accepts the EULA for a given release and product Version", func() {
			response := fmt.Sprintf(`{"releases":[{"id":40,"version":"3.3.0"},{"id":%d,"version":"%s"}]}`, releaseID, releaseVersion)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", releasesURL),
					ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
					ghttp.RespondWith(http.StatusOK, response),
				),
			)

			response = fmt.Sprintf(`{"accepted_at": "2016-01-11","_links":{}}`)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", EULAAcceptanceURL),
					ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
					ghttp.VerifyHeaderKV("User-Agent", userAgent),
					ghttp.VerifyJSON(`{}`),
					ghttp.RespondWith(http.StatusOK, response),
				),
			)

			Expect(client.AcceptEULA(pattern.PivnetFile{
				Slug:    productSlug,
				Version: releaseVersion,
				Glob:    "*.tgz",
			})).To(Succeed())
		})
	})

})
