package opsman_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/starkandwayne/om-configurator/configurator"
	. "github.com/starkandwayne/om-configurator/opsman"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("PivnetEULAAccepter", func() {
	var (
		server     *ghttp.Server
		accepter   PivnetEulaAccepter
		token      string
		apiAddress string
		apiPrefix  string
		userAgent  string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		apiAddress = server.URL()
		apiPrefix = "/api/v2"
		token = "my-auth-token"
		userAgent = "pivnet-resource/0.1.0 (some-url)"
		logger := log.New(GinkgoWriter, "", 0)
		accepter = PivnetEulaAccepter{
			Endpoint:  apiAddress,
			Token:     token,
			UserAgent: userAgent,
			Logger:    logger,
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Accept", func() {
		var (
			productVersion    string
			productSlug       string
			releaseID         int
			authURL           string
			releasesURL       string
			EULAAcceptanceURL string
		)

		BeforeEach(func() {
			productSlug = "banana-slug"
			productVersion = "3.2.1"
			releaseID = 42
			authURL = apiPrefix + "/authentication"
			releasesURL = fmt.Sprintf(apiPrefix+"/products/%s/releases", productSlug)
			EULAAcceptanceURL = fmt.Sprintf(apiPrefix+"/products/%s/releases/%d/pivnet_resource_eula_acceptance", productSlug, releaseID)
		})

		It("accepts the EULA for a given release and product Version", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", authURL),
					ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
					ghttp.RespondWith(http.StatusOK, `{}`),
				),
			)

			response := fmt.Sprintf(`{"releases":[{"id":40,"version":"3.2.0"},{"id":%d,"version":"%s"}]}`, releaseID, productVersion)
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

			Expect(accepter.Accept(configurator.DownloadProductArgs{
				PivnetProductSlug:    productSlug,
				PivnetProductVersion: productVersion,
			})).To(Succeed())
		})
	})

})
