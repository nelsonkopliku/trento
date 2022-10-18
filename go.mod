module github.com/trento-project/trento

go 1.16

require (
	github.com/avast/retry-go/v4 v4.1.0
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-gonic/gin v1.8.1
	github.com/gomarkdown/markdown v0.0.0-20210514010506-3b9f47219fe7
	github.com/google/uuid v1.3.0
	github.com/hooklift/gowsdl v0.5.0
	github.com/lib/pq v1.10.7
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/afero v1.9.2
	github.com/spf13/cobra v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.13.0
	github.com/stretchr/testify v1.8.0
	github.com/swaggo/files v0.0.0-20220728132757-551d4a08d97a
	github.com/swaggo/gin-swagger v1.5.3
	github.com/swaggo/swag v1.8.7
	github.com/tdewolff/minify/v2 v2.12.4
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/vektra/mockery/v2 v2.12.3
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	gorm.io/datatypes v1.0.2
	gorm.io/driver/postgres v1.4.4
	gorm.io/gorm v1.23.7
)

replace github.com/trento-project/trento => ./
