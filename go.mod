module riser

replace github.com/tshak/riser/sdk => ./sdk

go 1.12

require (
	github.com/alexeyco/simpletable v0.0.0-20190222165044-2eb48bcee7cf
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/sanity-io/litter v1.1.0
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	github.com/tshak/riser-server/api/v1/model v0.0.0-20190717133547-bd9d68203c33
	github.com/tshak/riser/sdk v0.0.0-20190705130421-2e250dea579a
	github.com/wzshiming/ctc v1.2.0
	github.com/wzshiming/winseq v0.0.0-20181031094240-8a45cfbfe1c2 // indirect
	golang.org/x/sys v0.0.0-20190712062909-fae7ac547cb7 // indirect
	k8s.io/apimachinery v0.0.0-20190717022731-0bb8574e0887 // indirect
	k8s.io/klog v0.3.3 // indirect
)
