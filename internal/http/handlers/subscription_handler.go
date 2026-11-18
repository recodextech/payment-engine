package handlers

import (
	"context"
	"net/http"
	"net/url"
	"payment-engine/pkg/errors"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/services"

	"github.com/recodextech/container"
	"github.com/recodextech/krouter"
)

// SubscriptionHandler serves the app, admin and internal subscription routes.
type SubscriptionHandler struct {
	subscriptions *services.SubscriptionService
	quota         *services.QuotaService
}

func (h *SubscriptionHandler) Init(c container.Container) error {
	h.subscriptions = c.Resolve(application.ModuleSubscriptionService).(*services.SubscriptionService)
	h.quota = c.Resolve(application.ModuleQuotaService).(*services.QuotaService)
	return nil
}

func queryValue(ctx context.Context, key string) string {
	if q, ok := ctx.Value(request.QueryParams.String()).(url.Values); ok {
		return q.Get(key)
	}
	return ""
}

func badRequest(msg string) error {
	return services.NewServiceError(http.StatusBadRequest, services.CodeInvalidRequest, msg)
}

func contractorID(p krouter.HttpPayload) string {
	return p.Param(request.PathParamContractorID.String())
}

// respond maps a service result to its HTTP response, wrapping errors with request context.
func respond[M, R any](m M, err error, toResponse func(M) R) (interface{}, error) {
	if err != nil {
		return nil, errors.Wrap(err, "subscription request failed")
	}
	return toResponse(m), nil
}

// ---- app ----

// ListPlans: GET /subscription/plans?kind=SUBSCRIPTION|ADDON
func (h *SubscriptionHandler) ListPlans(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	plans, err := h.subscriptions.ListPlans(ctx, p.Header(request.HeaderAccountID.String()), queryValue(ctx, "kind"))
	return respond(plans, err, responses.NewPlanListResponse)
}

// GetSubscription: GET /contractor/{contractor_id}/subscription
func (h *SubscriptionHandler) GetSubscription(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	sub, err := h.subscriptions.GetCurrent(ctx, contractorID(p))
	return respond(sub, err, responses.NewSubscriptionResponse)
}

// BankTransfer: GET /contractor/{contractor_id}/subscription/bank-transfer
func (h *SubscriptionHandler) BankTransfer(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	details, err := h.subscriptions.BankTransferDetails(ctx, contractorID(p))
	return respond(details, err, responses.NewBankTransferResponse)
}

// ---- admin ----

func planInput(b request.PlanRequest) models.PlanInput {
	return models.PlanInput{Code: b.Code, Kind: b.Kind, Status: b.Status, Name: b.Name, Price: b.Price, JobLimit: b.JobLimit}
}

func (h *SubscriptionHandler) AdminListPlans(ctx context.Context, _ krouter.HttpPayload) (interface{}, error) {
	plans, err := h.subscriptions.AdminListPlans(ctx)
	return respond(plans, err, responses.NewPlanListResponse)
}

func (h *SubscriptionHandler) AdminCreatePlan(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.PlanRequest)
	plan, err := h.subscriptions.AdminCreatePlan(ctx, planInput(body))
	return respond(plan, err, responses.NewPlanResponse)
}

func (h *SubscriptionHandler) AdminUpdatePlan(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.PlanRequest)
	plan, err := h.subscriptions.AdminUpdatePlan(ctx, p.Param(request.PathParamPlanID.String()), planInput(body))
	return respond(plan, err, responses.NewPlanResponse)
}

func (h *SubscriptionHandler) AdminListSubscriptions(ctx context.Context, _ krouter.HttpPayload) (interface{}, error) {
	status, kind := queryValue(ctx, "status"), queryValue(ctx, "kind")
	switch status {
	case "", domain.SubscriptionActive, domain.SubscriptionEnded:
	default:
		return nil, badRequest("status must be ACTIVE or ENDED")
	}
	switch kind {
	case "", domain.PlanKindSubscription, domain.PlanKindAddon:
	default:
		return nil, badRequest("kind must be SUBSCRIPTION or ADDON")
	}
	rows, err := h.subscriptions.AdminListSubscriptions(ctx, status, kind)
	return respond(rows, err, responses.NewSubscriptionRowListResponse)
}

// AdminListPayments: GET /admin/contractor/{contractor_id}/subscription/payments
func (h *SubscriptionHandler) AdminListPayments(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	pays, err := h.subscriptions.AdminListPayments(ctx, contractorID(p))
	return respond(pays, err, responses.NewPaymentListResponse)
}

func (h *SubscriptionHandler) AdminGrant(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.GrantRequest)
	if body.PlanCode == "" {
		return nil, badRequest("plan_code is required")
	}
	sub, err := h.subscriptions.AdminGrant(ctx, contractorID(p), body.PlanCode, body.Reason)
	return respond(sub, err, responses.NewSubscriptionResponse)
}

// AdminRecordPayment: POST /admin/subscription/payments
func (h *SubscriptionHandler) AdminRecordPayment(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.PaymentRequest)
	if (body.ContractorID == "" && body.Reference == "") || body.PlanCode == "" {
		return nil, badRequest("contractor_id or reference, and plan_code are required")
	}
	if body.Amount < 0 {
		return nil, badRequest("amount must not be negative")
	}
	sub, err := h.subscriptions.AdminRecordPayment(ctx, models.PaymentInput{ContractorID: body.ContractorID, Reference: body.Reference, PlanCode: body.PlanCode,
		Amount: body.Amount, BankReference: body.BankReference, Note: body.Note})
	return respond(sub, err, responses.NewSubscriptionResponse)
}

// ---- internal (demand-handler) ----

func (h *SubscriptionHandler) ConsumeQuota(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.QuotaConsumeRequest)
	if body.ContractorID == "" {
		return nil, badRequest("contractor_id is required")
	}
	quota, err := h.quota.Consume(ctx, body.ContractorID, body.ProcessID)
	return respond(quota, err, responses.NewQuotaResponse)
}

func (h *SubscriptionHandler) ReleaseQuota(ctx context.Context, p krouter.HttpPayload) (interface{}, error) {
	body, _ := p.Body.(request.QuotaReleaseRequest)
	if body.ProcessID == "" {
		return nil, badRequest("process_id is required")
	}
	switch body.Reason {
	case domain.ReleaseJobCancelled, domain.ReleaseJobExpired, domain.ReleaseProcessCreateError, domain.ReleaseRescheduleError:
	default:
		return nil, badRequest("reason must be JOB_CANCELLED, JOB_EXPIRED, PROCESS_CREATE_FAILED or RESCHEDULE_FAILED")
	}
	quota, err := h.quota.Release(ctx, body.ProcessID, body.Reason)
	return respond(quota, err, responses.NewQuotaResponse)
}

func (h *SubscriptionHandler) Reconcile(ctx context.Context, _ krouter.HttpPayload) (interface{}, error) {
	result, err := h.subscriptions.Reconcile(ctx)
	return respond(result, err, responses.NewReconcileResponse)
}
