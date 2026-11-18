package repositories

const (
	columnParamKey       = "key"
	columnParamMeta      = "meta"
	columnParamPayload   = "payload"
	columnParamProcessId = "process_id"
	columnParamStatus    = "status"
	columnParamJobType   = "job_type"
	columnParamDelete    = "delete"
	columnParamUpdatedAt = "updated_at"

	columnParamJobID  = "job_id"
	columnParamAmount = "amount"
	columnParamType   = "type"

	columnParamPaymentID = "payment_id"
	columnParamPayeeID   = "payee_id"

	// Account wallet column params
	columnParamAccountID = "account_id"
	columnParamBalance   = "balance"
	columnParamDeleted   = "deleted"

	// Subscription plan column params
	columnParamCode      = "code"
	columnParamKind      = "kind"
	columnParamDeletedAt = "deleted_at"

	// Subscription column params
	columnParamContractorID  = "contractor_id"
	columnParamPlanID        = "plan_id"
	columnParamPeriodEnd     = "period_end"
	columnParamRemainingJobs = "remaining_jobs"

	// Subscription payment column params
	columnParamProvider    = "provider"
	columnParamProviderRef = "provider_ref"
)
