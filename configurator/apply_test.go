package configurator_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"

	. "github.com/starkandwayne/om-configurator/configurator"
	"github.com/starkandwayne/om-configurator/configurator/configuratorfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply", func() {
	var (
		fakeOpsman *configuratorfakes.FakeOpsmanClient
		deployment Deployment
	)

	assetsDir := func() string {
		_, filename, _, _ := runtime.Caller(0)
		return filepath.Join(filepath.Dir(filename), "assets")
	}

	loadYAMLAsset := func(file string, out interface{}) {
		data, err := ioutil.ReadFile(filepath.Join(assetsDir(), file))
		Expect(err).ToNot(HaveOccurred())
		err = yaml.Unmarshal(data, out)
		Expect(err).ToNot(HaveOccurred())
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
			newOpsman := func(_ Opsman, _ *log.Logger) (OpsmanClient, error) {
				return fakeOpsman, nil
			}
			loadYAMLAsset("deployment_with_tiles.yml", &deployment)

			logger := log.New(GinkgoWriter, "", 0)
			templateStore := http.Dir(assetsDir())
			configurator, err := NewConfigurator(&deployment, templateStore, newOpsman, logger)
			Expect(err).ToNot(HaveOccurred())
			err = configurator.Apply()
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
