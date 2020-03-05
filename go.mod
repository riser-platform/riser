module riser

replace github.com/riser-platform/riser/sdk => ./sdk

go 1.13

require (
	github.com/AlecAivazis/survey/v2 v2.0.4
	// Pinning to 25d852a  until they get their release act together or until go-yaml/yaml.v2 supports json tags in which case we can remove this dep
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/strfmt v0.19.3 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/hashicorp/go-version v1.2.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-runewidth v0.0.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.7-0.20200303170754-4348ed7c0671
	github.com/riser-platform/riser/sdk v0.0.0-00010101000000-000000000000
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/wzshiming/ctc v1.2.0
	github.com/wzshiming/winseq v0.0.0-20181031094240-8a45cfbfe1c2 // indirect
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13 // indirect
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.17.3
)
