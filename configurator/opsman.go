package configurator

type DownloadProductArgs struct {
	OutputDirectory      string
	PivnetProductSlug    string
	PivnetProductVersion string
	PivnetProductGlob    string
	StemcellIaas         string
}

type StageProductArgs struct {
	ProductName    string
	ProductVersion string
}

//go:generate counterfeiter . OpsmanClient
type OpsmanClient interface {
	ConfigureAuthentication() error
	DownloadProduct(DownloadProductArgs) error
	UploadProduct(string) error
	UploadStemcell(string) error
	StageProduct(StageProductArgs) error
	ConfigureDirector([]byte) error
	ConfigureProduct([]byte) error
	ApplyChanges() error
}
