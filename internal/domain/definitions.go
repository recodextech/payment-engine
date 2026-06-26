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
	PaymentPending   PaymentStatus = `IN_PROGRESS`
	PaymentSuccess   PaymentStatus = `SUCCESS`
	PaymentCancelled PaymentStatus = `CANCELLED`
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
	WalletCard   WalletType = "CARD"
)

// ReasonCode classifies why an inter-account value transfer is occurring.
// Every payment must carry one of these codes so auditing is always possible.
type ReasonCode string

func (r ReasonCode) String() string {
	return string(r)
}

const (
	// ReasonJobCompleted is the standard payout from contractor to worker after
	// a job finishes.
	ReasonJobCompleted ReasonCode = "JOB_COMPLETED"
	// ReasonFeeDeduction covers platform or service fees taken from a wallet.
	ReasonFeeDeduction ReasonCode = "FEE_DEDUCTION"
	// ReasonRefund reverses a prior payment back to the originating wallet.
	ReasonRefund ReasonCode = "REFUND"
	// ReasonPointsRedemption is the only authorised path for converting Points
	// wallet value into a payout — must go through the dedicated redemption service.
	ReasonPointsRedemption ReasonCode = "POINTS_REDEMPTION"
)

// job status
type JobStatus string

func (j JobStatus) String() string {
	return string(j)
}

const (
	JobAccepted  JobStatus = `ACCEPTED`
	JobStarted   JobStatus = `STARTED`
	JobCancelled JobStatus = `CANCELLED`
	JobEnded     JobStatus = `ENDED`
)
