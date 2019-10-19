module riser

replace github.com/tshak/riser/sdk => ./sdk

go 1.13

require (
	github.com/AlecAivazis/survey/v2 v2.0.4
	github.com/alexeyco/simpletable v0.0.0-20190222165044-2eb48bcee7cf
	github.com/c-bata/go-prompt v0.2.3
	// Pinning to 25d852a  until they get their release act together or until go-yaml/yaml.v2 supports json tags in which case we can remove this dep
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/hashicorp/go-version v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.4.0
	github.com/tshak/riser-server/api/v1/model v0.0.0-20191006090436-fbaa2a96f3a7
	github.com/tshak/riser/sdk v0.0.0-00010101000000-000000000000
	github.com/wzshiming/ctc v1.2.0
	github.com/wzshiming/winseq v0.0.0-20181031094240-8a45cfbfe1c2 // indirect
	golang.org/x/sys v0.0.0-20190902133755-9109b7679e13 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4
)
