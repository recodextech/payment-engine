package repositories

import (
	"context"

	"payment-engine/internal/domain/models"

	"github.com/recodextech/api-definitions/events"
)

// PlanRepository reads and writes subscription plan business data.
type PlanRepository interface {
	ListPlans(ctx context.Context, kind, status string) ([]events.SubscriptionPlanPayload, error)
	GetPlanByCode(ctx context.Context, code string) (events.SubscriptionPlanPayload, error)
	GetPlanByID(ctx context.Context, id string) (events.SubscriptionPlanPayload, error)
	CreatePlan(ctx context.Context, plan events.SubscriptionPlanPayload) error
	UpdatePlan(ctx context.Context, plan events.SubscriptionPlanPayload) error
}

// SubscriptionRepository stores subscription rows, payments and the job quota ledger.
// Methods that lock rows must run inside an InTransaction callback.
type SubscriptionRepository interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	// LockContractor serialises all subscription changes of one contractor.
	LockContractor(ctx context.Context, contractorID string) error

	// ListActiveRows returns a contractor's ACTIVE rows, soonest ending first.
	ListActiveRows(ctx context.Context, contractorID string) ([]events.SubscriptionEvent, error)
	GetRowByID(ctx context.Context, id string) (events.SubscriptionEvent, bool, error)
	ListRows(ctx context.Context, status, kind string, limit int) ([]events.SubscriptionEvent, error)
	// ListDueContractors returns contractors with an ACTIVE row that has ended.
	ListDueContractors(ctx context.Context, nowMillis int64, limit int) ([]string, error)
	CreateRow(ctx context.Context, row events.SubscriptionEvent) error
	EndRow(ctx context.Context, row events.SubscriptionEvent) error
	// TakeJob takes one job from the row if it has one left, in one statement.
	TakeJob(ctx context.Context, rowID string) (remaining int, ok bool, err error)
	// ReturnJob gives one job back to the row while it is ACTIVE and has not ended.
	ReturnJob(ctx context.Context, rowID string, nowMillis int64) (ok bool, err error)

	// CreatePayment fails with RepositoryDuplicateError when the provider reference was recorded before.
	CreatePayment(ctx context.Context, payment events.SubscriptionPaymentEvent) error
	SetPaymentStatus(ctx context.Context, payment events.SubscriptionPaymentEvent) error
	// ListPaidPayments returns paid payments not used yet, oldest first.
	ListPaidPayments(ctx context.Context, contractorID string) ([]events.SubscriptionPaymentEvent, error)
	ListPayments(ctx context.Context, contractorID string, limit int) ([]events.SubscriptionPaymentEvent, error)

	GetLedgerByProcessForUpdate(ctx context.Context, processID string) (events.QuotaLedgerEvent, bool, error)
	// SaveLedger inserts the process's slot, or takes a RELEASED one again.
	SaveLedger(ctx context.Context, entry events.QuotaLedgerEvent) error
}

// ContractorAccountRepository resolves the account behind a contractor.
type ContractorAccountRepository interface {
	GetContractorAccount(ctx context.Context, contractorID string) (models.ContractorAccount, error)
	GetContractorAccountByNumber(ctx context.Context, accountNumber int64) (models.ContractorAccount, error)
	ContractorBelongsToAccount(ctx context.Context, contractorID, accountID string) (bool, error)
}
