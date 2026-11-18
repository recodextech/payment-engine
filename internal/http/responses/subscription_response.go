package responses

import (
	"payment-engine/internal/domain/models"

	"github.com/recodextech/api-definitions/events"
)

// PlanResponse is a plan as shown to the app and the admin portal.
type PlanResponse struct {
	ID       string  `json:"id"`
	Code     string  `json:"code"`
	Kind     string  `json:"kind"`
	Status   string  `json:"status"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency,omitempty"`
	JobLimit int     `json:"job_limit"`
}

func NewPlanResponse(p models.Plan) PlanResponse {
	return PlanResponse{ID: p.ID, Code: p.Code, Kind: p.Kind, Status: p.Status, Name: p.Name, Price: p.Price,
		Currency: p.Currency, JobLimit: p.JobLimit}
}

func NewPlanListResponse(plans []models.Plan) []PlanResponse {
	out := make([]PlanResponse, 0, len(plans))
	for _, p := range plans {
		out = append(out, NewPlanResponse(p))
	}
	return out
}

// PlanSnapshotResponse is the plan a subscription was taken on.
type PlanSnapshotResponse struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	JobLimit int     `json:"job_limit"`
}

func planSnapshot(p events.PlanSnapshot) PlanSnapshotResponse {
	return PlanSnapshotResponse{Code: p.Code, Price: p.Price, Currency: p.Currency, JobLimit: p.JobLimit}
}

// SubscriptionRowResponse is one monthly plan or add-on pack row.
type SubscriptionRowResponse struct {
	ID               string               `json:"id"`
	ContractorID     string               `json:"contractor_id"`
	Kind             string               `json:"kind"`
	Status           string               `json:"status"`
	Plan             PlanSnapshotResponse `json:"plan"`
	PeriodStart      int64                `json:"period_start"`
	PeriodEnd        int64                `json:"period_end"`
	RemainingJobs    int                  `json:"remaining_jobs"`
	PaymentID        string               `json:"payment_id,omitempty"`
	EndReason        string               `json:"end_reason,omitempty"`
	PaymentReference string               `json:"payment_reference"`
}

func NewSubscriptionRowResponse(r models.SubscriptionRow) SubscriptionRowResponse {
	return SubscriptionRowResponse{ID: r.ID, ContractorID: r.ContractorID, Kind: r.Kind, Status: r.Status, Plan: planSnapshot(r.Plan),
		PeriodStart: r.PeriodStart, PeriodEnd: r.PeriodEnd, RemainingJobs: r.RemainingJobs, PaymentID: r.PaymentID,
		EndReason: r.EndReason, PaymentReference: r.PaymentReference}
}

func NewSubscriptionRowListResponse(rows []models.SubscriptionRow) []SubscriptionRowResponse {
	out := make([]SubscriptionRowResponse, 0, len(rows))
	for _, r := range rows {
		out = append(out, NewSubscriptionRowResponse(r))
	}
	return out
}

// SubscriptionResponse is what a contractor can use now: the monthly plan and the active add-on packs.
// remaining_jobs counts both; jobs_used = job_limit + addon_jobs - remaining_jobs.
type SubscriptionResponse struct {
	ID                 string                    `json:"id"`
	ContractorID       string                    `json:"contractor_id"`
	Status             string                    `json:"status"`
	Plan               PlanSnapshotResponse      `json:"plan"`
	CurrentPeriodStart int64                     `json:"current_period_start"`
	CurrentPeriodEnd   int64                     `json:"current_period_end"`
	JobLimit           int                       `json:"job_limit"`
	AddonJobs          int                       `json:"addon_jobs"`
	RemainingJobs      int                       `json:"remaining_jobs"`
	PlanRemainingJobs  int                       `json:"plan_remaining_jobs"`
	JobsUsed           int                       `json:"jobs_used"`
	PaymentReference   string                    `json:"payment_reference"`
	PrepaidPeriods     int                       `json:"prepaid_periods"`
	Addons             []SubscriptionRowResponse `json:"addons"`
}

func NewSubscriptionResponse(s models.ContractorSubscription) SubscriptionResponse {
	p := s.Plan
	return SubscriptionResponse{
		ID: p.ID, ContractorID: p.ContractorID, Status: p.Status, Plan: planSnapshot(p.Plan),
		CurrentPeriodStart: p.PeriodStart, CurrentPeriodEnd: p.PeriodEnd, JobLimit: p.Plan.JobLimit, AddonJobs: s.AddonJobs,
		RemainingJobs: s.RemainingJobs, PlanRemainingJobs: p.RemainingJobs, JobsUsed: p.Plan.JobLimit + s.AddonJobs - s.RemainingJobs,
		PaymentReference: s.PaymentReference, PrepaidPeriods: s.PrepaidPeriods, Addons: NewSubscriptionRowListResponse(s.Addons),
	}
}

// PaymentResponse is a recorded payment for a plan or add-on pack.
type PaymentResponse struct {
	ID           string            `json:"id"`
	ContractorID string            `json:"contractor_id"`
	Provider     string            `json:"provider"`
	ProviderRef  string            `json:"provider_ref"`
	Status       string            `json:"status"`
	PlanCode     string            `json:"plan_code"`
	Amount       float64           `json:"amount"`
	Currency     string            `json:"currency"`
	Attributes   map[string]string `json:"attributes"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    int64             `json:"created_at"`
}

func NewPaymentListResponse(pays []models.SubscriptionPayment) []PaymentResponse {
	out := make([]PaymentResponse, 0, len(pays))
	for _, p := range pays {
		attrs := p.Attributes
		if attrs == nil {
			attrs = map[string]string{}
		}
		out = append(out, PaymentResponse{ID: p.ID, ContractorID: p.ContractorID, Provider: p.Provider, ProviderRef: p.ProviderRef,
			Status: p.Status, PlanCode: p.PlanCode, Amount: p.Amount, Currency: p.Currency, Attributes: attrs,
			CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt})
	}
	return out
}

// BankTransferResponse is where and how a contractor pays by bank transfer.
type BankTransferResponse struct {
	BankName      string `json:"bank_name"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	Branch        string `json:"branch"`
	Reference     string `json:"reference"`
	Currency      string `json:"currency"`
}

func NewBankTransferResponse(d models.BankTransferDetails) BankTransferResponse {
	return BankTransferResponse{BankName: d.BankName, AccountName: d.AccountName, AccountNumber: d.AccountNumber,
		Branch: d.Branch, Reference: d.Reference, Currency: d.Currency}
}

// QuotaResponse is the result of taking or returning a process's job slot.
type QuotaResponse struct {
	Consumed  int `json:"consumed"`
	Remaining int `json:"remaining_jobs"`
}

func NewQuotaResponse(q models.QuotaResult) QuotaResponse {
	return QuotaResponse{Consumed: q.Consumed, Remaining: q.Remaining}
}

// ReconcileResponse summarises a daily sweep.
type ReconcileResponse struct {
	Checked int `json:"checked"`
	Changed int `json:"changed"`
	Failed  int `json:"failed"`
}

func NewReconcileResponse(r models.ReconcileResult) ReconcileResponse {
	return ReconcileResponse{Checked: r.Checked, Changed: r.Changed, Failed: r.Failed}
}
