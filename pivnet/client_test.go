package pivnet_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	})

	JustBeforeEach(func() {
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

		JustBeforeEach(func() {
			response := fmt.Sprintf(`{"releases":[{"id":40,"version":"3.3.0"},{"id":%d,"version":"%s"}]}`, releaseID, releaseVersion)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", releasesURL),
					ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
					ghttp.RespondWith(http.StatusOK, response),
				),
			)

			response = fmt.Sprintf(`{"accepted_at": "2016-01-11","_links":{}}`)
			handlers := []http.HandlerFunc{
				ghttp.VerifyRequest("POST", EULAAcceptanceURL),
				ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
				ghttp.VerifyJSON(`{}`),
				ghttp.RespondWith(http.StatusOK, response),
			}
			if userAgent != "" {
				handlers = append(handlers, ghttp.VerifyHeaderKV("User-Agent", userAgent))
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(handlers...),
			)

		})

		Context("given a normal user who needs to accept manually", func() {
			BeforeEach(func() {
				productSlug = "banana-slug"
				releaseVersion = "3.2"
				releaseID = 42
				releasesURL = fmt.Sprintf(apiPrefix+"/products/%s/releases", productSlug)
				EULAAcceptanceURL = fmt.Sprintf(apiPrefix+"/products/%s/releases/%d/pivnet_resource_eula_acceptance", productSlug, releaseID)
			})

			It("accepts the EULA for a given release and product Version", func() {

				Expect(client.AcceptEULA(context.Background(), pattern.PivnetFile{
					Slug:    productSlug,
					Version: releaseVersion,
					Glob:    "*.tgz",
				})).To(Succeed())
			})
		})

		Context("given an agent which is allowed to accept EULA's automatically", func() {
			BeforeEach(func() {
				productSlug = "banana-slug"
				releaseVersion = "3.2"
				releaseID = 42
				releasesURL = fmt.Sprintf(apiPrefix+"/products/%s/releases", productSlug)
				EULAAcceptanceURL = fmt.Sprintf(apiPrefix+"/products/%s/releases/%d/eula_acceptance", productSlug, releaseID)
				userAgent = "special client"
				acceptEULA = true
			})

			It("accepts the EULA for a given release and product Version", func() {
				Expect(client.AcceptEULA(context.Background(), pattern.PivnetFile{
					Slug:    productSlug,
					Version: releaseVersion,
					Glob:    "*.tgz",
				})).To(Succeed())
			})
		})

	})

	Describe("DownloadFile", func() {
		var (
			releaseVersion  string
			releaseID       int
			releasesURL     string
			productFilesURL string
			fileID          int
			fileName        string
			fileURL         string
			fileContent     string
			downloadFileURL string
			file            pattern.PivnetFile
			downloadDir     string
		)

		Context("Given a PivnetFile with valid Slug, Version and Glob", func() {
			JustBeforeEach(func() {
				response := fmt.Sprintf(`{"releases":[{"id":40,"version":"3.3.0"},{"id":%d,"version":"%s"}]}`, releaseID, releaseVersion)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", releasesURL),
						ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
						ghttp.RespondWith(http.StatusOK, response),
					),
				)
				response = fmt.Sprintf(`{"product_files":[{"id":40,"aws_object_key":"foo.pivotal"},{"id":%d,"aws_object_key":"%s"}]}`, fileID, fileName)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", productFilesURL),
						ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
						ghttp.RespondWith(http.StatusOK, response),
					),
				)
				response = fmt.Sprintf(`{"product_file":{"_links":{"download":{"href": "%s"}}}}`, fileURL)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", fmt.Sprintf("%s/%d", productFilesURL, fileID)),
						ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
						ghttp.RespondWith(http.StatusOK, response),
					),
				)
				responseHeader := http.Header{"Location": []string{fmt.Sprintf("%s%s", server.URL(), downloadFileURL)}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", fmt.Sprintf("/api/v2%s", fileURL)),
						ghttp.VerifyHeaderKV("Authorization", fmt.Sprintf("Token %s", token)),
						ghttp.RespondWith(http.StatusFound, "", responseHeader),
					),
				)
				responseHeader = http.Header{"Content-Length": []string{fmt.Sprintf("%d", len(fileContent))}}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("HEAD", downloadFileURL),
						ghttp.RespondWith(http.StatusOK, "", responseHeader),
					),
				)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", downloadFileURL),
						ghttp.RespondWith(http.StatusOK, fileContent),
					),
				)
			})

			BeforeEach(func() {
				acceptEULA = false
				releasesURL = "/api/v2/products/elastic-runtime/releases"
				productFilesURL = "/api/v2/products/elastic-runtime/releases/22/product_files"
				releaseID = 22
				releaseVersion = "2.4"
				fileID = 18
				fileName = "foo-srt-2.4.pivotal"
				fileURL = "/products/elastic-runtime/releases/22/product_files/18/download"
				downloadFileURL = "/s3/download/foo-srt-2.4.pivotal"
				fileContent = "c"
				file = pattern.PivnetFile{
					Slug:    "elastic-runtime",
					Version: releaseVersion,
					Glob:    "*srt*.pivotal",
				}
				var err error
				downloadDir, err = ioutil.TempDir("", "downloadDir")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(downloadDir)
			})

			It("Downloads tile from pivnet", func() {
				_, err := client.DownloadFile(context.Background(), file, downloadDir)
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("Given a PivnetFile with download url", func() {
			JustBeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", downloadFileURL),
						ghttp.RespondWith(http.StatusOK, fileContent),
					),
				)
			})

			BeforeEach(func() {
				downloadFileURL = "/cache/foo-srt-2.4.pivotal"
				fileContent = "c"

				file = pattern.PivnetFile{
					Slug:    "elastic-runtime",
					Version: releaseVersion,
					Glob:    "*srt*.pivotal",
					URL:     fmt.Sprintf("%s%s", server.URL(), downloadFileURL),
				}
				var err error
				downloadDir, err = ioutil.TempDir("", "downloadDir")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(downloadDir)
			})

			It("Downloads tile directly", func() {
				_, err := client.DownloadFile(context.Background(), file, downloadDir)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
