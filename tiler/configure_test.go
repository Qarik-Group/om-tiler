package tiler_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/starkandwayne/om-tiler/pattern"
	. "github.com/starkandwayne/om-tiler/tiler"
	"github.com/starkandwayne/om-tiler/tiler/tilerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configure", func() {
	var (
		fakeOpsman *tilerfakes.FakeOpsmanClient
		fakeMover  *tilerfakes.FakeMover
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
			fakeOpsman = &tilerfakes.FakeOpsmanClient{}
			fakeMover = &tilerfakes.FakeMover{
				GetStub: func(f pattern.PivnetFile) (*os.File, error) {
					return ioutil.TempFile("", f.Slug)
				},
			}
			logger := log.New(GinkgoWriter, "", 0)
			templateStore := http.Dir(assetsDir())
			tiler, err := NewTiler(fakeOpsman, fakeMover, logger)
			Expect(err).ToNot(HaveOccurred())
			p, err := pattern.NewPattern(pattern.Template{
				Manifest: "pattern.yml",
				Vars: map[string]interface{}{
					"iaas-configuration_project":   "example-project",
					"iaas-configuration_name":      "bar",
					"real-iaas-configuration_name": "foo",
					"network_name":                 "network1",
				},
				Store: templateStore,
			})
			Expect(err).ToNot(HaveOccurred())
			err = tiler.Configure(p)
			Expect(err).ToNot(HaveOccurred())
		})

		It("Configures the director", func() {
			config := fakeOpsman.ConfigureDirectorArgsForCall(0)
			Expect(config).To(MatchYAML(readAsset("results/director-config.yml")))
		})

		It("Downloads the tiles and stemcells from Pivotal Network", func() {
			args := fakeMover.GetArgsForCall(0)
			Expect(args.Slug).To(Equal("p-healthwatch"))
			Expect(args.Version).To(Equal("1.2.3"))
			Expect(args.Glob).To(Equal("*.pivotal"))

			args = fakeMover.GetArgsForCall(1)
			Expect(args.Slug).To(Equal("stemcells-ubuntu-xenial"))
			Expect(args.Version).To(Equal("170.38"))
			Expect(args.Glob).To(Equal("*vsphere*.tgz"))

			args = fakeMover.GetArgsForCall(2)
			Expect(args.Slug).To(Equal("elastic-runtime"))
			Expect(args.Version).To(Equal("3.2.1"))
			Expect(args.Glob).To(Equal("srt*.pivotal"))

			args = fakeMover.GetArgsForCall(3)
			Expect(args.Slug).To(Equal("stemcells-ubuntu-trusty"))
			Expect(args.Version).To(Equal("170.50"))
			Expect(args.Glob).To(Equal("*gcp*.tgz"))
		})

		It("Uploads the tiles and stemcells to Ops Manager", func() {
			Expect(fakeOpsman.UploadProductArgsForCall(0).Name()).
				To(ContainSubstring("p-healthwatch"))
			Expect(fakeOpsman.UploadStemcellArgsForCall(0).Name()).
				To(ContainSubstring("stemcells-ubuntu-xenial"))
			Expect(fakeOpsman.UploadProductArgsForCall(1).Name()).
				To(ContainSubstring("elastic-runtime"))
			Expect(fakeOpsman.UploadStemcellArgsForCall(1).Name()).
				To(ContainSubstring("stemcells-ubuntu-trusty"))
		})

		It("Stages the products", func() {
			args := fakeOpsman.StageProductArgsForCall(0)
			Expect(args.Name).To(Equal("p-healthwatch"))
			Expect(args.Version).To(Equal("1.2.3-build.1"))

			args = fakeOpsman.StageProductArgsForCall(1)
			Expect(args.Name).To(Equal("cf"))
			Expect(args.Version).To(Equal("3.2.1"))
		})

		It("Configures the products", func() {
			config := fakeOpsman.ConfigureProductArgsForCall(0)
			Expect(config).To(MatchYAML(readAsset("results/p-healthwatch.yml")))
		})
	})

})
