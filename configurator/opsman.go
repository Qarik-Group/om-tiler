package configurator

type DownloadProductArgs struct {
	OutputDirectory      string
	PivnetProductSlug    string
	PivnetProductVersion string
	PivnetProductGlob    string
	StemcellIaas         string
}

//go:generate counterfeiter . OpsmanClient
type OpsmanClient interface {
	ConfigureAuthentication() error
	DownloadProduct(DownloadProductArgs) error
	UploadProduct(string) error
	UploadStemcell(string) error
	ConfigureProduct([]byte) error
	ApplyChanges() error
}
