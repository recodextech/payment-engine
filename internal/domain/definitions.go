package domain

const (
	TripCreateHandler = `http.handler.trip-create`
)

type ContextKey string

const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
)

const (
	ContextKeyAccountID ContextKey = `account-id`
	ContextKeyTraceID   ContextKey = `trace-id`
	ContextKeyUserID    ContextKey = `user-id`
	ContextKeyUserType  ContextKey = `triggered-user-type`
	ContextKeyStreamID  ContextKey = "stream-id"
	ContextKeyTimeZone  ContextKey = "time-zone"
)

func (c ContextKey) String() string {
	return string(c)
}

// process types
type ProcessType string

func (p ProcessType) String() string {
	return string(p)
}

const (
	RideProcess    ProcessType = `RIDE`
	HailingProcess ProcessType = `HAILING`
)

// payment status
type PaymentStatus string

func (p PaymentStatus) String() string {
	return string(p)
}

const (
	PaymentPending PaymentStatus = `IN_PROGRESS`
	PaymentSuccess PaymentStatus = `SUCCESS`
)

const (
	WalletStatusActive string = "ACTIVE"
	WalletStatusHold   string = "HOLD"
)

type WalletType string

func (w WalletType) String() string {
	return string(w)
}

const (
	WalletCash   WalletType = "CASH"
	WalletPoints WalletType = "POINTS"
	WalletCredit WalletType = "CREDIT"
)
