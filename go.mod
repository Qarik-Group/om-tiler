module github.com/starkandwayne/om-tiler

require (
	cloud.google.com/go v0.44.3 // indirect
	code.cloudfoundry.org/clock v0.0.0-20180518195852-02e53af36e6c // indirect
	code.cloudfoundry.org/workpool v0.0.0-20170718174546-99757edba735 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/bmatcuk/doublestar v1.1.5 // indirect
	github.com/cheggaaa/pb v2.0.7+incompatible // indirect
	github.com/cloudfoundry/bosh-agent v2.247.0+incompatible // indirect
	github.com/cloudfoundry/bosh-cli v6.0.0+incompatible
	github.com/cloudfoundry/bosh-davcli v0.0.34 // indirect
	github.com/cloudfoundry/bosh-gcscli v0.0.6 // indirect
	github.com/cloudfoundry/bosh-s3cli v0.0.81 // indirect
	github.com/cloudfoundry/bosh-utils v0.0.0-20190817100217-701ecef9262d
	github.com/cloudfoundry/config-server v0.1.18
	github.com/cloudfoundry/go-socks5 v0.0.0-20180221174514-54f73bdb8a8e // indirect
	github.com/cloudfoundry/socks5-proxy v0.2.0 // indirect
	github.com/cppforlife/go-patch v0.2.0
	github.com/cppforlife/go-semi-semantic v0.0.0-20160921010311-576b6af77ae4 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/google/uuid v1.1.1
	github.com/graymeta/stow v0.2.3 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/kamildrazkiewicz/go-flow v0.0.0-20170919072006-2a1d885f8ce4
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/onsi/ginkgo v1.9.0
	github.com/onsi/gomega v1.6.0
	github.com/pivotal-cf/go-pivnet v1.0.3
	github.com/pivotal-cf/om v0.0.0-20190816215002-d607995f0947
	github.com/pivotal-cf/paraphernalia v0.0.0-20180203224945-a64ae2051c20 // indirect
	github.com/pivotal-cf/pivnet-cli v0.0.60
	github.com/square/certstrap v1.1.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/tedsuo/ifrit v0.0.0-20180802180643-bea94bb476cc // indirect
	github.com/vito/go-interact v1.0.0 // indirect
	golang.org/x/mobile v0.0.0-20190814143026-e8b3e6111d02 // indirect
	golang.org/x/tools v0.0.0-20190820033707-85edb9ef3283 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55 // indirect
	google.golang.org/grpc v1.23.0 // indirect
	gopkg.in/VividCortex/ewma.v1 v1.1.1 // indirect
	gopkg.in/cheggaaa/pb.v2 v2.0.7 // indirect
	gopkg.in/fatih/color.v1 v1.7.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/mattn/go-runewidth.v0 v0.0.4 // indirect
	gopkg.in/yaml.v2 v2.2.2
	honnef.co/go/tools v0.0.1-2019.2.2 // indirect
)

replace (
	github.com/cheggaaa/pb => github.com/cheggaaa/pb v1.0.28 // from bosh-cli Gopkg.lock
	github.com/jessevdk/go-flags => github.com/cppforlife/go-flags v0.0.0-20170707010757-351f5f310b26
)
