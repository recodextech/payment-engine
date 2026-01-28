package request

type ContextParam string

func (c ContextParam) String() string {
	return string(c)
}

type Header string

func (r Header) String() string {
	return string(r)
}

// Request header parameters
const (
	HeaderAccountID Header = `account-id`
	HeaderTraceID   Header = `trace-id`
	HeaderUserID    Header = `user-id`
)

// Custom param types
const (
	ParamTypeAppUUID   = "app-uuid"
	ParamTypeAppInt64  = "app-int64"
	ParamTypeAppString = "string-param"
)

const (
	UserID    ContextParam = "user-id"
	AccountID ContextParam = "account-id"
	TraceID   ContextParam = "trace-id"
	XHeaders  ContextParam = "x-headers"
)

// Request path parameters
const (
	PathParamWorkerID     ContextParam = "worker_id"
	PathParamContractorID ContextParam = "contractor_id"
)
