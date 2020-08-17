module riser

go 1.15

require (
	github.com/AlecAivazis/survey/v2 v2.0.8
	// Pinning to 25d852a  until they get their release act together or until go-yaml/yaml.v2 supports json tags in which case we can remove this dep
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/go-ozzo/ozzo-validation/v3 v3.8.1
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/go-version v1.2.1
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/kr/pty v1.1.8 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/riser-platform/riser-server/api/v1/model v0.0.20
	github.com/riser-platform/riser-server/pkg/sdk v0.0.46
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/whilp/git-urls v0.0.0-20191001220047-6db9661140c0
	github.com/wzshiming/ctc v1.2.3
	golang.org/x/tools v0.0.0-20200717024301-6ddee64345a6 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apimachinery v0.18.8
)
