package log

import (
	"context"
	"io"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartpet/websocket/constant"
)

type Level string
type ContextKey string

const (
	partyCodeCtxKey = "_ctx_party_code"
	requestIdCtxKey = "_ctx_request_id"
)

const (
	PartyCodeLogKey = "PartyCode"
	LogType         = "log_type"
	AccessLog       = "access"
	ApplicationLog  = "application"
	MetricsLog      = "metrics"
	TransactionsLog = "transaction"
)

func addCategoryLog(event *zerolog.Event, logCategory string) *zerolog.Event {
	return event.Str(LogType, logCategory)
}

func (l Level) zeroLogLevel() zerolog.Level {
	switch l {
	case constant.TraceLevel:
		return zerolog.TraceLevel
	case constant.DebugLevel:
		return zerolog.DebugLevel
	case constant.InfoLevel:
		return zerolog.InfoLevel
	case constant.WarnLevel:
		return zerolog.WarnLevel
	case constant.ErrorLevel:
		return zerolog.ErrorLevel
	case constant.FatalLevel:
		return zerolog.FatalLevel
	case constant.PanicLevel:
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}

// InitLogger is used to initialize logger
func InitLogger(level Level) {
	zerolog.ErrorStackMarshaler = getErrorStackMarshaller()
	zerolog.SetGlobalLevel(level.zeroLogLevel())
	log.Logger = log.With().Caller().Logger()
}

// Trace is the for trace log
func Trace(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Trace())
}

// Debug is the for debug log
func Debug(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Debug())
}

// Info is the for info log
func Info(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Info())
}

// Warn is the for warn log
func Warn(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Warn())
}

// Error is the for error log
func Error(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Error().Stack())
}

// Panic is the for panic log
func Panic(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Panic().Stack())
}

// Fatal is the for fatal log
func Fatal(ctx context.Context) *zerolog.Event {
	return withIDAndPath(ctx, log.Fatal().Stack())
}

func getErrorStackMarshaller() func(err error) interface{} {
	return func(err error) interface{} {
		if err != nil {
			if e, ok := err.(*LogError); ok {
				return map[string]interface{}{
					constant.CodeLogParam:    e.Code,
					constant.MessageLogParam: e.Message,
					constant.DetailsLogParam: e.Details,
					constant.TraceLogParam:   e.GetTrace(),
				}
			}
		}
		return string(debug.Stack())
	}
}

func withIDAndPath(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if ctx == nil {
		return event
	}
	id := ctx.Value(constant.IDLogParam)
	if id != nil {
		event.Interface(constant.IDLogParam, id)
	}
	path := ctx.Value(constant.PathLogParam)
	if path != nil {
		event.Interface(constant.PathLogParam, path)
	}
	correlationId := ctx.Value(constant.CorrelationLogParam)
	if correlationId != nil {
		event.Interface(constant.CorrelationLogParam, correlationId)
	}
	return event
}

func withUserData(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if ctx == nil {
		return event
	}

	userData := ctx.Value(constant.UserData)
	if userMap, ok := userData.(map[string]interface{}); ok {
		event.Interface(constant.UserId, userMap[constant.UserId])
		event.Interface(constant.Source, userMap[constant.Source])
		event.Interface(constant.AppId, userMap[constant.AppId])
	}

	// print request uin
	requestUINHeader := ctx.Value(constant.RequestUIN)
	if requestUINHeader != nil {
		event.Interface(constant.RequestUIN, requestUINHeader)
	}

	transactionType := ctx.Value(constant.TransactionType)
	if transactionType != nil {
		event.Interface(constant.TransactionType, transactionType)
	}

	// resource: https://pkg.go.dev/context#example-WithValue
	orderId := ctx.Value(ContextKey(constant.UserOrderId))
	if orderId != nil {
		event.Interface(constant.UserOrderId, orderId)
	}

	requestUIN := ctx.Value(ContextKey(constant.RequestUIN))
	if requestUIN != nil {
		event.Interface(constant.RequestUIN, requestUIN)
	}

	requestedExchange := ctx.Value(ContextKey(constant.RequestedExchange))
	if requestedExchange != nil {
		event.Interface(constant.RequestedExchange, requestedExchange)
	}

	return event
}

// ***********************APPLICATION METRIC LOGS******************************

// Debug is the for debug log
func MetricDebug(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Debug(ctx))
	return addCategoryLog(event, MetricsLog)
}

// Info is the for info log
func MetricInfo(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Info(ctx))
	return addCategoryLog(event, MetricsLog)
}

// Warn is the for warn log
func MetricWarn(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Warn(ctx))
	return addCategoryLog(event, MetricsLog)
}

// Error is the for error log
func MetricError(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Error(ctx))
	return addCategoryLog(event, MetricsLog)
}

// Fatal is the for fatal log
func MetricFatal(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Fatal(ctx))
	return addCategoryLog(event, MetricsLog)
}

// ***********************APPLICATION OTHER LOGS******************************

// Debug is the for debug log
func ApplicationDebug(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Debug(ctx))
	return addCategoryLog(event, ApplicationLog)
}

// Info is the for info log
func ApplicationInfo(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Info(ctx))
	return addCategoryLog(event, ApplicationLog)
}

// Warn is the for warn log
func ApplicationWarn(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Warn(ctx))
	return addCategoryLog(event, ApplicationLog)
}

// Error is the for error log
func ApplicationError(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Error(ctx))
	return addCategoryLog(event, ApplicationLog)
}

// Fatal is the for fatal log
func ApplicationFatal(ctx context.Context) *zerolog.Event {
	event := addPartyCodeToLog(ctx, Fatal(ctx))
	return addCategoryLog(event, ApplicationLog)
}

func addPartyCodeToLog(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	partyCode, ok := ctx.Value(partyCodeCtxKey).(string)
	if ok {
		return event.Str(PartyCodeLogKey, partyCode)
	}
	return event
}

const (
	LogTypeKey = "log_type"

	// Will be defaulted to Application Log Type if no log type is specified.
	LogTypeAccess                 = "access"
	LogTypeAudit                  = "audit"
	LogTypeTransaction            = "transaction"
	LogTypeApplication            = "application"
	LogTypeApplicationPerformance = "application-performance"
	LogTypeProcessAudit           = "process-audit"
)

// InitLoggerWithWriter is used to initialize logger with a writer
func InitLoggerWithWriter(level Level, w io.Writer) {
	zerolog.ErrorStackMarshaler = getErrorStackMarshaller()
	zerolog.SetGlobalLevel(level.zeroLogLevel())
	log.Logger = zerolog.New(w).With().Caller().Timestamp().Logger()
}

// ErrorWarn checks for the error object.
// In case it is corresponding to a 4XX status code, it logs it as warning.
// Otherwise, it logs it as an error.

type logTypeEvent struct {
	logType string
}

func Access() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeAccess}
}

func Audit() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeAudit}
}

func Transaction() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeTransaction}
}

func Application() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeApplication}
}

func ApplicationPerformance() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeApplicationPerformance}
}

func ProcessAudit() *logTypeEvent {
	return &logTypeEvent{logType: LogTypeProcessAudit}
}
