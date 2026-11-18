package domain

import "github.com/recodextech/api-definitions/pkg/uuid"

// SystemUserID is recorded as meta.created_by for changes without a calling user
// (the daily sweep, a month rolling over when the contractor opens the app).
var SystemUserID = func() uuid.UUID {
	id, err := uuid.Parse("00000000-0000-0000-0000-00000000fe01")
	if err != nil {
		panic(err)
	}
	return id
}()

// DefaultCurrency is used when an account has no currency preference yet.
const DefaultCurrency = "LKR"

// FreePlanCode is the plan every contractor falls back to.
const FreePlanCode = "FREE"

// PeriodMonths is how long a monthly plan or an add-on pack lasts from the day it starts.
const PeriodMonths = 1

// Plan kinds (also the kind of a fix.subscriptions row) and plan statuses
const (
	PlanKindSubscription = "SUBSCRIPTION"
	PlanKindAddon        = "ADDON"

	PlanStatusActive   = "ACTIVE"
	PlanStatusArchived = "ARCHIVED"
)

// Subscription row statuses and why a row ended
const (
	SubscriptionActive = "ACTIVE"
	SubscriptionEnded  = "ENDED"

	EndExpired     = "EXPIRED"
	EndPlanChanged = "PLAN_CHANGED"
)

// Payment providers and statuses. PAID is paid but not used yet (a month paid in advance);
// APPLIED has opened a subscription row.
const (
	ProviderBankTransfer = "BANK_TRANSFER"
	ProviderAdmin        = "ADMIN" // given by an admin without payment

	SubPaymentPending = "PENDING"
	SubPaymentPaid    = "PAID"
	SubPaymentApplied = "APPLIED"
	SubPaymentFailed  = "FAILED"

	// payment attribute keys
	AttrReference     = "reference" // the contractor's reference written on the transfer, e.g. AA1111111
	AttrBankReference = "bank_reference"
	AttrNote          = "note"
)

// Quota ledger: one slot per process
const (
	QuotaConsumed = "CONSUMED"
	QuotaReleased = "RELEASED"

	ReleaseJobCancelled       = "JOB_CANCELLED"
	ReleaseJobExpired         = "JOB_EXPIRED"
	ReleaseProcessCreateError = "PROCESS_CREATE_FAILED"
	ReleaseRescheduleError    = "RESCHEDULE_FAILED"
)

// Push notification types sent to the contractor (Firebase via notification-service)
const (
	NotifySubscriptionActivated   = "subscription_activated"
	NotifySubscriptionRenewed     = "subscription_renewed"
	NotifySubscriptionPaymentRecv = "subscription_payment_received"
	NotifySubscriptionExpired     = "subscription_expired"
	NotifySubscriptionAddon       = "subscription_addon_purchased"
	NotifyQuotaLow                = "subscription_quota_low"
	NotifyQuotaExhausted          = "subscription_quota_exhausted"
)
