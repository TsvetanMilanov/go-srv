package srv

const (
	// AppLoggerName is the name of the app logger instance.
	AppLoggerName = "appLogger"

	// ReqLoggerName is the name of the req logger instance.
	ReqLoggerName = "reqLogger"

	// AppDIName is the name of the app di container.
	AppDIName = "appDi"

	// ReqDIName is the name of the req di container.
	ReqDIName = "reqDi"

	// TraceIDName is the name of the trace id property.
	TraceIDName = "traceId"

	// TraceIDReqHeaderName is the name of the header which will be used
	// to acquire trace id.
	TraceIDReqHeaderName = "X-Trace-Id"

	defaultMetricsServerAddr = ":80"
)
