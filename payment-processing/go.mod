module spenmo/payment-processing/payment-processing

go 1.13

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jmoiron/sqlx v1.3.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/rs/xid v1.3.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.8.1
	gitlab.com/opaper/goutils/httpmiddleware v1.0.0 // indirect
	gitlab.com/opaper/goutils/log v1.1.0
)

replace github.com/spf13/afero => github.com/spf13/afero v1.5.1
