package configurator_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/starkandwayne/om-configurator/configurator"
	"github.com/starkandwayne/om-configurator/configurator/configuratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply", func() {
	var (
		fakeOpsman *configuratorfakes.FakeOpsmanClient
	)

	assetsDir := func() string {
		_, filename, _, _ := runtime.Caller(0)
		return filepath.Join(filepath.Dir(filename), "assets")
	}

	readAsset := func(file string) []byte {
		data, err := ioutil.ReadFile(filepath.Join(assetsDir(), file))
		Expect(err).ToNot(HaveOccurred())
		return data
	}

	Context("Given a deployment with products", func() {
		BeforeEach(func() {
			fakeOpsman = &configuratorfakes.FakeOpsmanClient{
				DownloadProductStub: func(c DownloadProductArgs) error {
					_, err := os.Create(filepath.Join(
						c.OutputDirectory,
						fmt.Sprintf("%s-%s.pivotal",
							c.PivnetProductSlug,
							c.PivnetProductVersion,
						),
					))
					Expect(err).ToNot(HaveOccurred())
					_, err = os.Create(filepath.Join(
						c.OutputDirectory,
						fmt.Sprintf("stemcell-%s.tgz",
							c.StemcellIaas,
						),
					))
					Expect(err).ToNot(HaveOccurred())
					return nil
				},
			}

			logger := log.New(GinkgoWriter, "", 0)
			templateStore := http.Dir(assetsDir())
			config := Config{
				Target:               "example.com",
				Username:             "username",
				Password:             "password",
				DecryptionPassphrase: "decrypt",
				PivnetToken:          "token",
			}
			configurator, err := NewConfigurator(&config, templateStore, fakeOpsman, logger)
			Expect(err).ToNot(HaveOccurred())
			err = configurator.Apply("deployment_with_tiles.yml")
			Expect(err).ToNot(HaveOccurred())
		})

		It("Downloads the tiles and stemcells from Pivotal Network", func() {
			args := fakeOpsman.DownloadProductArgsForCall(0)
			Expect(args.PivnetProductSlug).To(Equal("p-healthwatch"))
			Expect(args.PivnetProductVersion).To(Equal("1.2.3"))
			Expect(args.PivnetProductGlob).To(Equal("*.pivotal"))
			Expect(args.StemcellIaas).To(Equal("vsphere"))

			args = fakeOpsman.DownloadProductArgsForCall(1)
			Expect(args.PivnetProductSlug).To(Equal("elastic-runtime"))
			Expect(args.PivnetProductVersion).To(Equal("3.2.1"))
			Expect(args.PivnetProductGlob).To(Equal("srt*.pivotal"))
			Expect(args.StemcellIaas).To(Equal("gcp"))
		})

		It("Uploads the tiles and stemcells to Ops Manager", func() {
			Expect(fakeOpsman.UploadProductArgsForCall(0)).
				To(HaveSuffix("p-healthwatch-1.2.3.pivotal"))
			Expect(fakeOpsman.UploadStemcellArgsForCall(0)).
				To(HaveSuffix("stemcell-vsphere.tgz"))
			Expect(fakeOpsman.UploadProductArgsForCall(1)).
				To(HaveSuffix("elastic-runtime-3.2.1.pivotal"))
			Expect(fakeOpsman.UploadStemcellArgsForCall(1)).
				To(HaveSuffix("stemcell-gcp.tgz"))
		})

		It("Configures the products", func() {
			config := fakeOpsman.ConfigureProductArgsForCall(0)
			Expect(config).To(MatchYAML(readAsset("p-healthwatch-merged.yml")))
		})
	})

})
