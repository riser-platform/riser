module riser

replace github.com/tshak/riser/sdk => ./sdk

go 1.12

require (
	github.com/alexeyco/simpletable v0.0.0-20190222165044-2eb48bcee7cf
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	// Pinning to 25d852a  until they get their release act together or until go-yaml/yaml.v2 supports json tags in which case we can remove this dep
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/sanity-io/litter v1.1.0
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tshak/riser-server/api/v1/model v0.0.0-20190902150527-37f7f5496fb0
	github.com/tshak/riser/sdk v0.0.0-00010101000000-000000000000
	github.com/wzshiming/ctc v1.2.0
	github.com/wzshiming/winseq v0.0.0-20181031094240-8a45cfbfe1c2 // indirect
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381 // indirect
)
