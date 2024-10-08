package constant

const (
	ActionKey         = "action"
	StatusCodeKey     = "statusCode"
	LatencyKey        = "latency"
	ClientIPKey       = "clientIP"
	MethodKey         = "method"
	PathKey           = "path"
	ErrorKey          = "error"
	LoggerId          = "id"
	LogPath           = "logPath"
	HTTPConfigKey     = "httpConfig"
	DatabaseConfigKey = "databaseConfig"
	LogLevelKey       = "logLevel"
)

const (
	UserData        = "user"
	UserId          = "user_id"
	Source          = "source"
	AppId           = "app_id"
	RequestIDHeader = "X-requestId"
)
const (
	IDLogParam        = "id"
	EndTimeLogParam   = "endTime"
	StartTimeLogParam = "startTime"
	QueryLogParam     = "query"
	HeaderLogParams   = "header"
	UserOrderId       = "order_id"
	RequestUIN        = "request_uin"
	TransactionType   = "transaction_type"
	RequestedExchange = "requested_exchange"
)

// Context Params automatically added to logs
const (
	PathLogParam        = "path"
	CorrelationLogParam = "correlationId"
	ClientIDLogParam    = "clientID"
	ActionLogParam      = "action"
	S2SIssuerLogParam   = "s2sIssuer"
)

// Other additional Log Params
const (
	StatusCodeLogParam = "statusCode"
	URILogParam        = "uri"
	LatencyLogParam    = "latency"
	ClientIPLogParam   = "clientIP"
	MethodLogParam     = "method"
	ErrorLogParam      = "error"
	CodeLogParam       = "code"
	MessageLogParam    = "message"
	DetailsLogParam    = "details"
	TraceLogParam      = "trace"
)
