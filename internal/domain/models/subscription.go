package models

import "github.com/recodextech/api-definitions/events"

// ContractorAccount is the account data subscriptions need about a contractor.
type ContractorAccount struct {
	ContractorID  string
	AccountID     string
	AccountNumber int64 // number in the contractor's bank transfer reference
	UserID        string
	Currency      string // account preference, empty when not set
	Name          string
	Email         string
	Phone         string
}

// Plan is a monthly plan or add-on pack, priced in Currency when known.
type Plan struct {
	ID       string
	Code     string
	Kind     string
	Status   string
	Name     string
	Price    float64
	Currency string
	JobLimit int
}

// PlanInput is the data to create or update a plan.
type PlanInput struct {
	Code     string
	Kind     string
	Status   string
	Name     string
	Price    float64
	JobLimit int
}

// SubscriptionRow is one fix.subscriptions row: a monthly plan or an add-on pack.
type SubscriptionRow struct {
	ID            string
	ContractorID  string
	Kind          string
	Status        string
	Plan          events.PlanSnapshot
	PeriodStart   int64
	PeriodEnd     int64
	RemainingJobs int
	PaymentID     string
	EndReason     string
	// PaymentReference is what the contractor writes on bank transfers, e.g. AA1111111
	PaymentReference string
}

// ContractorSubscription is what a contractor can use now: the monthly plan and active add-on packs.
type ContractorSubscription struct {
	Plan   SubscriptionRow
	Addons []SubscriptionRow
	// RemainingJobs is what the contractor can still post: the plan's and the add-ons' remaining jobs
	RemainingJobs int
	// AddonJobs is the jobs the active add-on packs came with
	AddonJobs        int
	PaymentReference string
	// PrepaidPeriods are paid months waiting to start when the current plan month ends
	PrepaidPeriods int
}

// SubscriptionPayment is a recorded payment for a plan or add-on pack.
type SubscriptionPayment struct {
	ID           string
	ContractorID string
	Provider     string
	ProviderRef  string
	Status       string
	PlanCode     string
	Amount       float64
	Currency     string
	Attributes   map[string]string
	CreatedBy    string
	CreatedAt    int64
}

// BankTransferDetails tells a contractor where to transfer a payment and which reference to use.
type BankTransferDetails struct {
	BankName      string
	AccountName   string
	AccountNumber string
	Branch        string
	Reference     string
	Currency      string
}

// PaymentInput is a bank transfer an admin records, for a contractor picked in the admin portal or
// found by the reference written on the transfer.
type PaymentInput struct {
	ContractorID  string
	Reference     string
	PlanCode      string
	Amount        float64
	BankReference string
	Note          string
}

// QuotaResult is the contractor's job quota after a consume or release.
type QuotaResult struct {
	// Consumed is 1 when this call took the process's slot, 0 when the process already held one.
	Consumed  int
	Remaining int
}

// ReconcileResult summarises a daily sweep.
type ReconcileResult struct {
	Checked int
	Changed int
	Failed  int
}
