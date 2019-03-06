package configurator

type DownloadProductArgs struct {
	OutputDirectory      string
	PivnetProductSlug    string
	PivnetProductVersion string
	PivnetProductGlob    string
	StemcellIaas         string
}

type UploadProductArgs struct {
	ProductFilePath      string
	PivnetProductVersion string
}

type UploadStemcellArgs struct {
	StemcellFilePath string
}

//go:generate counterfeiter . OpsmanClient
type OpsmanClient interface {
	ConfigureAuthentication() error
	DownloadProduct(DownloadProductArgs) error
	UploadProduct(UploadProductArgs) error
	UploadStemcell(UploadStemcellArgs) error
	ConfigureProduct(string) error
	ApplyChanges() error
}
