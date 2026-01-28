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

// process status
type ProcessStatus string

func (p ProcessStatus) String() string {
	return string(p)
}

const (
	ProcessPending   ProcessStatus = `PENDING`
	ProcessStarted   ProcessStatus = `STARTED`
	ProcessCompleted ProcessStatus = `COMPLETED`
	ProcessCancelled ProcessStatus = `CANCELLED`
)

// job types
type JobType string

func (j JobType) String() string {
	return string(j)
}

const (
	JobRide JobType = `RIDE`
)

// job status
type JobStatus string

func (j JobStatus) String() string {
	return string(j)
}

const (
	JobPending   JobStatus = `PENDING`
	JobAssigned  JobStatus = `ASSIGNED`
	JobAccepeted JobStatus = `ACCEPTED`
	JobCancelled JobStatus = `CANCELLED`
	JobCompleted JobStatus = `COMPLETED`
)

// location types
type JobLocationType string

func (j JobLocationType) String() string {
	return string(j)
}

const (
	Pickup JobLocationType = `PICKUP`
	Drop   JobLocationType = `DROP`
)

// job asssignment status
type JobAssignmentStatus string

func (j JobAssignmentStatus) String() string {
	return string(j)
}

const (
	JobAssignmentPending   JobAssignmentStatus = `PENDING`
	JobAssignmentAccepeted JobAssignmentStatus = `ACCEPTED`
)
