package constant

const (
	LoggerConfig      = "logger"
	ApplicationConfig = "application"
	DatabaseConfig    = "database"
	AWSConfig         = "aws"
	ExternalConfig    = "external"
)

// Database constant
const (
	URLConfigKey                         = "url"
	HTTPConnectTimeoutInMillisKey        = "http.connectTimeoutInMillis"
	HTTPKeepAliveDurationInMillisKey     = "http.keepAliveDurationInMillis"
	HTTPMaxIdleConnectionsKey            = "http.maxIdleConnections"
	HTTPIdleConnectionTimeoutInMillisKey = "http.idleConnectionTimeoutInMillis"
	HTTPTlsHandshakeTimeoutInMillisKey   = "http.tlsHandshakeTimeoutInMillis"
	HTTPExpectContinueTimeoutInMillisKey = "http.expectContinueTimeoutInMillis"
	HTTPTimeoutInMillisKey               = "http.timeoutInMillis"

	SqlServerDriverName                       = "mssql"
	MySqlServerDriverName                     = "mysql"
	PGDriverName                              = "postgres"
	DatabaseMaxOpenConnectionsKey             = "maxOpenConnections"
	DatabaseMaxIdleConnectionsKey             = "maxIdleConnections"
	DatabaseConnectionMaxLifetimeInSecondsKey = "connectionMaxLifetimeInSeconds"
	DatabaseConnectionMaxIdleTimeInSecondsKey = "connectionMaxIdleTimeInSeconds"
	DatabasePortConfigKey                     = "port"

	SqlserverDB               = "sqlserver"
	DMATSqlserver200          = "dmatsqlserver200"
	DMATSqlserver201          = "dmatsqlserver201"
	DMATSqlserver202          = "dmatsqlserver202"
	MySqlserverDB             = "mysqlserver"
	PgserverDB                = "postgres"
	DatabaseServerConfigKey   = "server"
	DatabaseUsernameConfigKey = "username"
	DatabasePasswordConfigKey = "password"
	DatabaseNameConfigKey     = "databaseName"
	WhiteListOriginUrl        = "white_list_origin_urls"
)

const (
	LogLevelConfigKey     = "level"
	ConsoleLoggingEnabled = "ConsoleLoggingEnabled"
	EncodeLogsAsJson      = "EncodeLogsAsJson"
	FileLoggingEnabled    = "FileLoggingEnabled"
	Directory             = "Directory"
	Filename              = "Filename"
	MaxSize               = "MaxSize"
	MaxBackups            = "MaxBackups"
	MaxAge                = "MaxAge"
)
const (
	ProfileService   = "profile"
	EdisService      = "verifyedis"
	PortfolioService = "portfolio"
	CNSService       = "sns"
	CheckTPin        = "checkedis"
)
