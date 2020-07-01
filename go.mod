module riser

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.0.4
	// Pinning to 25d852a  until they get their release act together or until go-yaml/yaml.v2 supports json tags in which case we can remove this dep
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/strfmt v0.19.3 // indirect
	github.com/go-ozzo/ozzo-validation/v3 v3.8.1
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.14
	github.com/riser-platform/riser-server/pkg/sdk v0.0.37-0.20200630175822-14d8484e6fbd
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	github.com/whilp/git-urls v0.0.0-20191001220047-6db9661140c0
	github.com/wzshiming/ctc v1.2.0
	github.com/wzshiming/winseq v0.0.0-20181031094240-8a45cfbfe1c2 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
	golang.org/x/tools v0.0.0-20190729092621-ff9f1409240a // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.18.5
)
