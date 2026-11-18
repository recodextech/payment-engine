package services

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"payment-engine/pkg/errors"
	"regexp"
	"strings"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"

	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

// SubscriptionService serves plans, contractor subscriptions, bank transfer details and admin actions.
type SubscriptionService struct {
	core     *SubscriptionCore
	accounts repositories.AccountRepository
	log      adaptors.Logger
}

func (s *SubscriptionService) Init(c container.Container) error {
	s.core = c.Resolve(application.ModuleSubscriptionCore).(*SubscriptionCore)
	s.accounts = c.Resolve(application.ModuleAccountRepo).(repositories.AccountRepository)
	s.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`services.subscription`))
	return nil
}

func planModel(p events.SubscriptionPlanPayload, currency string) models.Plan {
	return models.Plan{ID: p.ID.String(), Code: p.Code, Kind: p.Kind, Status: p.Status, Name: p.Name, Price: p.Price,
		Currency: currency, JobLimit: p.JobLimit}
}

func rowModel(r events.SubscriptionEvent, reference string) models.SubscriptionRow {
	p := r.Payload
	m := models.SubscriptionRow{ID: p.ID.String(), ContractorID: p.ContractorID.String(), Kind: p.Kind, Status: p.Status,
		Plan: p.Plan, PeriodStart: p.PeriodStart, PeriodEnd: p.PeriodEnd, RemainingJobs: p.RemainingJobs,
		EndReason: p.EndReason, PaymentReference: reference}
	if p.PaymentID != nil {
		m.PaymentID = p.PaymentID.String()
	}
	return m
}

func (s *SubscriptionService) subscriptionModel(ctx context.Context, acc models.ContractorAccount, c current) (models.ContractorSubscription, error) {
	ref := s.core.reference(acc)
	view := models.ContractorSubscription{Plan: rowModel(c.Plan, ref), RemainingJobs: c.remaining(), PaymentReference: ref}
	for _, a := range c.Addons {
		view.Addons = append(view.Addons, rowModel(a, ref))
		view.AddonJobs += a.Payload.Plan.JobLimit
	}
	paid, err := s.core.subs.ListPaidPayments(ctx, acc.ContractorID)
	view.PrepaidPeriods = len(paid)
	if err != nil {
		return view, errors.Wrap(err, "failed to list paid payments")
	}
	return view, nil
}

func paymentModel(p events.SubscriptionPaymentEvent) models.SubscriptionPayment {
	return models.SubscriptionPayment{ID: p.Payload.ID.String(), ContractorID: p.Payload.ContractorID.String(), Provider: p.Payload.Provider,
		ProviderRef: p.Payload.ProviderRef, Status: p.Payload.Status, PlanCode: p.Payload.PlanCode, Amount: p.Payload.Amount,
		Currency: p.Payload.Currency, Attributes: p.Payload.Attributes, CreatedBy: p.EventMeta.CreatedBy.String(),
		CreatedAt: p.EventMeta.CreatedAt}
}

// ---- app ----

// ListPlans returns ACTIVE plans priced in the calling account's currency.
func (s *SubscriptionService) ListPlans(ctx context.Context, accountID, kind string) ([]models.Plan, error) {
	if kind != "" && kind != domain.PlanKindSubscription && kind != domain.PlanKindAddon {
		return nil, serviceError(http.StatusBadRequest, CodeInvalidRequest, "kind must be SUBSCRIPTION or ADDON")
	}
	currency, err := s.accounts.GetAccountCurrency(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account currency")
	}
	if currency == "" {
		currency = domain.DefaultCurrency
	}
	plans, err := s.core.plans.ListPlans(ctx, kind, domain.PlanStatusActive)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list plans")
	}
	views := make([]models.Plan, 0, len(plans))
	for _, p := range plans {
		views = append(views, planModel(p, strings.ToUpper(currency)))
	}
	return views, nil
}

// GetCurrent returns the contractor's monthly plan and add-on packs, starting FREE or ending months that are over.
func (s *SubscriptionService) GetCurrent(ctx context.Context, contractorID string) (models.ContractorSubscription, error) {
	acc, err := s.core.contractorAccount(ctx, contractorID)
	if err != nil {
		return models.ContractorSubscription{}, errors.Wrap(err, "failed to get contractor account")
	}
	var view models.ContractorSubscription
	err = s.core.withContractor(ctx, acc, func(ctx context.Context, fx *effects) error {
		cur, err := s.core.resolveCurrent(ctx, acc, fx)
		if err != nil {
			return errors.Wrap(err, "failed to resolve current subscription")
		}
		view, err = s.subscriptionModel(ctx, acc, cur)
		if err != nil {
			return errors.Wrap(err, "failed to build subscription view")
		}
		return nil
	})
	if err != nil {
		return view, errors.Wrap(err, "failed to get current subscription")
	}
	return view, nil
}

// BankTransferDetails is the account a contractor pays plans and add-ons into, with the reference to use.
func (s *SubscriptionService) BankTransferDetails(ctx context.Context, contractorID string) (models.BankTransferDetails, error) {
	acc, err := s.core.contractorAccount(ctx, contractorID)
	if err != nil {
		return models.BankTransferDetails{}, errors.Wrap(err, "failed to get contractor account")
	}
	conf := s.core.conf
	return models.BankTransferDetails{
		BankName: conf.BankName, AccountName: conf.BankAccountName, AccountNumber: conf.BankAccountNumber,
		Branch: conf.BankBranch, Reference: s.core.reference(acc), Currency: currencyOf(acc),
	}, nil
}

// ---- admin: plans ----

var planCodePattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]{1,49}$`)

func validatePlan(in models.PlanInput) error {
	switch {
	case !planCodePattern.MatchString(in.Code):
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "code must be 2-50 uppercase letters, digits or _")
	case in.Kind != domain.PlanKindSubscription && in.Kind != domain.PlanKindAddon:
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "kind must be SUBSCRIPTION or ADDON")
	case in.Status != domain.PlanStatusActive && in.Status != domain.PlanStatusArchived:
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "status must be ACTIVE or ARCHIVED")
	case strings.TrimSpace(in.Name) == "":
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "name is required")
	case in.Price < 0:
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "price must not be negative")
	case in.JobLimit < 1:
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "job_limit must be at least 1")
	case in.Kind == domain.PlanKindAddon && in.Price == 0:
		return serviceError(http.StatusBadRequest, CodeInvalidRequest, "an add-on pack must have a price")
	}
	return nil
}

// AdminListPlans returns plans in every status.
func (s *SubscriptionService) AdminListPlans(ctx context.Context) ([]models.Plan, error) {
	plans, err := s.core.plans.ListPlans(ctx, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list plans")
	}
	views := make([]models.Plan, 0, len(plans))
	for _, p := range plans {
		views = append(views, planModel(p, ""))
	}
	return views, nil
}

func (s *SubscriptionService) AdminCreatePlan(ctx context.Context, in models.PlanInput) (models.Plan, error) {
	in.Code = strings.ToUpper(strings.TrimSpace(in.Code))
	if in.Status == "" {
		in.Status = domain.PlanStatusActive
	}
	if err := validatePlan(in); err != nil {
		return models.Plan{}, errors.Wrap(err, "invalid plan")
	}
	if _, err := s.core.plans.GetPlanByCode(ctx, in.Code); err == nil {
		return models.Plan{}, serviceError(http.StatusConflict, CodeInvalidPlan, "plan code "+in.Code+" already exists")
	} else if !isNotExist(err) {
		return models.Plan{}, errors.Wrap(err, "failed to check plan code")
	}
	plan := events.SubscriptionPlanPayload{ID: uuid.New(), Code: in.Code, Kind: in.Kind, Status: in.Status,
		Name: strings.TrimSpace(in.Name), Price: in.Price, JobLimit: in.JobLimit}
	if err := s.core.plans.CreatePlan(ctx, plan); err != nil {
		return models.Plan{}, errors.Wrap(err, "failed to create plan")
	}
	return planModel(plan, ""), nil
}

// AdminUpdatePlan edits a plan. The code can't change, and FREE must stay an active free monthly plan.
// A new job_limit applies to months and packs that start after the change.
func (s *SubscriptionService) AdminUpdatePlan(ctx context.Context, planID string, in models.PlanInput) (models.Plan, error) {
	plan, err := s.core.plans.GetPlanByID(ctx, planID)
	if err != nil {
		if isNotExist(err) {
			return models.Plan{}, serviceError(http.StatusNotFound, CodeInvalidPlan, "plan not found")
		}
		return models.Plan{}, errors.Wrap(err, "failed to get plan")
	}
	in.Code = plan.Code
	if in.Kind == "" {
		in.Kind = plan.Kind
	}
	if in.Status == "" {
		in.Status = plan.Status
	}
	if err := validatePlan(in); err != nil {
		return models.Plan{}, errors.Wrap(err, "invalid plan")
	}
	if plan.Code == domain.FreePlanCode && (in.Status != domain.PlanStatusActive || in.Kind != domain.PlanKindSubscription || in.Price != 0) {
		return models.Plan{}, serviceError(http.StatusBadRequest, CodeInvalidPlan, "the FREE plan must stay an active subscription with price 0")
	}
	plan.Kind, plan.Status, plan.Name, plan.Price, plan.JobLimit = in.Kind, in.Status, strings.TrimSpace(in.Name), in.Price, in.JobLimit
	if err := s.core.plans.UpdatePlan(ctx, plan); err != nil {
		return models.Plan{}, errors.Wrap(err, "failed to update plan")
	}
	return planModel(plan, ""), nil
}

// ---- admin: subscriptions and payments ----

// AdminListSubscriptions lists subscription rows (newest first), filtered by status and kind.
func (s *SubscriptionService) AdminListSubscriptions(ctx context.Context, status, kind string) ([]models.SubscriptionRow, error) {
	rows, err := s.core.subs.ListRows(ctx, status, kind, 500)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list subscriptions")
	}
	views := make([]models.SubscriptionRow, 0, len(rows))
	references := map[string]string{}
	for _, row := range rows {
		cid := row.Payload.ContractorID.String()
		ref, ok := references[cid]
		if !ok {
			if acc, err := s.core.accounts.GetContractorAccount(ctx, cid); err == nil {
				ref = s.core.reference(acc)
			}
			references[cid] = ref
		}
		views = append(views, rowModel(row, ref))
	}
	return views, nil
}

// AdminListPayments lists a contractor's payments, newest first.
func (s *SubscriptionService) AdminListPayments(ctx context.Context, contractorID string) ([]models.SubscriptionPayment, error) {
	if _, err := s.core.contractorAccount(ctx, contractorID); err != nil {
		return nil, errors.Wrap(err, "failed to get contractor account")
	}
	pays, err := s.core.subs.ListPayments(ctx, contractorID, 200)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list payments")
	}
	views := make([]models.SubscriptionPayment, 0, len(pays))
	for _, p := range pays {
		views = append(views, paymentModel(p))
	}
	return views, nil
}

func (s *SubscriptionService) planByCode(ctx context.Context, code string) (events.SubscriptionPlanPayload, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	plan, err := s.core.plans.GetPlanByCode(ctx, code)
	if err != nil && isNotExist(err) {
		return plan, serviceError(http.StatusBadRequest, CodeInvalidPlan, "unknown plan "+code)
	}
	if err != nil {
		return plan, errors.Wrap(err, "failed to get plan by code")
	}
	return plan, nil
}

// AdminGrant gives a contractor a monthly plan or an add-on pack without payment, from today for one
// month. It is recorded as an ADMIN payment of 0.
func (s *SubscriptionService) AdminGrant(ctx context.Context, contractorID, planCode, reason string) (models.ContractorSubscription, error) {
	acc, err := s.core.contractorAccount(ctx, contractorID)
	if err != nil {
		return models.ContractorSubscription{}, errors.Wrap(err, "failed to get contractor account")
	}
	plan, err := s.planByCode(ctx, planCode)
	if err != nil {
		return models.ContractorSubscription{}, errors.Wrap(err, "failed to get plan")
	}
	if plan.Code == domain.FreePlanCode {
		return models.ContractorSubscription{}, serviceError(http.StatusBadRequest, CodeInvalidPlan, "the FREE plan can't be granted")
	}
	attrs := map[string]string{}
	if r := strings.TrimSpace(reason); r != "" {
		attrs[domain.AttrNote] = r
	}
	return s.recordPayment(ctx, acc, plan, domain.ProviderAdmin, uuid.New().String(), 0, attrs)
}

// AdminRecordPayment records a bank transfer for a contractor picked in the admin portal, or found by
// the payment reference written on the transfer (AA + account number).
func (s *SubscriptionService) AdminRecordPayment(ctx context.Context, in models.PaymentInput) (models.ContractorSubscription, error) {
	acc, err := s.paymentContractor(ctx, in)
	if err != nil {
		return models.ContractorSubscription{}, errors.Wrap(err, "failed to resolve payment contractor")
	}
	plan, err := s.planByCode(ctx, in.PlanCode)
	if err != nil {
		return models.ContractorSubscription{}, errors.Wrap(err, "failed to get plan")
	}
	if plan.Price <= 0 {
		return models.ContractorSubscription{}, serviceError(http.StatusBadRequest, CodeInvalidPlan, "plan "+plan.Code+" is free")
	}
	amount := in.Amount
	if amount == 0 {
		amount = plan.Price
	}
	if math.Abs(amount-plan.Price) >= 0.005 {
		return models.ContractorSubscription{}, serviceError(http.StatusBadRequest, CodeAmountMismatch,
			fmt.Sprintf("amount %.2f doesn't match the %s price %.2f", amount, plan.Code, plan.Price))
	}
	attrs := map[string]string{domain.AttrReference: s.core.reference(acc)}
	ref := strings.TrimSpace(in.BankReference)
	if ref != "" {
		attrs[domain.AttrBankReference] = ref
	} else {
		ref = uuid.New().String()
	}
	if note := strings.TrimSpace(in.Note); note != "" {
		attrs[domain.AttrNote] = note
	}
	return s.recordPayment(ctx, acc, plan, domain.ProviderBankTransfer, ref, amount, attrs)
}

func (s *SubscriptionService) paymentContractor(ctx context.Context, in models.PaymentInput) (models.ContractorAccount, error) {
	if in.ContractorID != "" {
		return s.core.contractorAccount(ctx, in.ContractorID)
	}
	conf := s.core.conf
	no, ok := conf.ParseReference(in.Reference)
	if !ok {
		return models.ContractorAccount{}, serviceError(http.StatusBadRequest, CodeInvalidReference,
			fmt.Sprintf("%q is not a payment reference like %s", in.Reference, conf.Reference(1111111)))
	}
	acc, err := s.core.accounts.GetContractorAccountByNumber(ctx, no)
	if err != nil && isNotExist(err) {
		return acc, serviceError(http.StatusNotFound, CodeContractorNotFound, "no contractor has payment reference "+conf.Reference(no))
	}
	if err != nil {
		return acc, errors.Wrap(err, "failed to get contractor by account number")
	}
	return acc, nil
}

// recordPayment stores the payment and uses it:
//   - an add-on pack starts a new add-on row from now
//   - the monthly plan the contractor is already on is paid in advance: it stays PAID and starts the
//     next month when the current one ends
//   - any other monthly plan replaces the current one from now
func (s *SubscriptionService) recordPayment(ctx context.Context, acc models.ContractorAccount, plan events.SubscriptionPlanPayload,
	provider, providerRef string, amount float64, attrs map[string]string) (models.ContractorSubscription, error) {
	pay := events.SubscriptionPaymentEvent{EventMeta: newMeta(ctx, acc.AccountID)}
	pay.Payload = events.SubscriptionPaymentPayload{ID: uuid.New(), ContractorID: parseID(acc.ContractorID), Provider: provider,
		ProviderRef: providerRef, Status: domain.SubPaymentPaid, PlanCode: plan.Code, Amount: amount, Currency: currencyOf(acc),
		Attributes: attrs}

	var view models.ContractorSubscription
	err := s.core.withContractor(ctx, acc, func(ctx context.Context, fx *effects) error {
		cur, err := s.core.resolveCurrent(ctx, acc, fx)
		if err != nil {
			return errors.Wrap(err, "failed to resolve current subscription")
		}
		if err := s.core.subs.CreatePayment(ctx, pay); err != nil {
			if isDuplicate(err) {
				return serviceError(http.StatusConflict, CodeDuplicatePayment, "bank transaction "+providerRef+" is already recorded")
			}
			return errors.Wrap(err, "failed to record payment")
		}
		now := s.core.now()
		id := pay.Payload.ID
		switch {
		case plan.Kind == domain.PlanKindAddon:
			row, err := s.core.startRow(ctx, acc, plan, now, amount, &id)
			if err != nil {
				return errors.Wrap(err, "failed to start add-on")
			}
			cur.Addons = append(cur.Addons, row)
			fx.notify(addonNotice(row, cur.remaining()))
		case cur.Plan.Payload.PlanID == plan.ID:
			fx.notify(paymentReceivedNotice(cur.Plan))
			view, err = s.subscriptionModel(ctx, acc, cur)
			if err != nil {
				return errors.Wrap(err, "failed to build subscription view") // stays PAID until the current month ends
			}
			return nil // stays PAID until the current month ends
		default:
			if err := s.core.endRow(ctx, cur.Plan, domain.EndPlanChanged); err != nil {
				return errors.Wrap(err, "failed to end current plan")
			}
			cur.Plan, err = s.core.startRow(ctx, acc, plan, now, amount, &id)
			if err != nil {
				return errors.Wrap(err, "failed to start plan")
			}
			fx.notify(activatedNotice(cur.Plan))
		}
		if err := s.core.markApplied(ctx, &pay); err != nil {
			return errors.Wrap(err, "failed to mark payment applied")
		}
		view, err = s.subscriptionModel(ctx, acc, cur)
		if err != nil {
			return errors.Wrap(err, "failed to build subscription view")
		}
		return nil
	})
	if err != nil {
		return view, errors.Wrap(err, "failed to apply payment")
	}
	return view, nil
}

// ---- daily sweep ----

// Reconcile ends the rows whose month is over for every contractor, so contractors who don't open the
// app still fall back to FREE or start a month they paid in advance.
func (s *SubscriptionService) Reconcile(ctx context.Context) (models.ReconcileResult, error) {
	due, err := s.core.subs.ListDueContractors(ctx, s.core.now().UnixMilli(), 500)
	if err != nil {
		return models.ReconcileResult{}, errors.Wrap(err, "failed to list due contractors")
	}
	var result models.ReconcileResult
	for _, contractorID := range due {
		result.Checked++
		acc, err := s.core.contractorAccount(ctx, contractorID)
		if err == nil {
			err = s.core.withContractor(ctx, acc, func(ctx context.Context, fx *effects) error {
				_, err := s.core.resolveCurrent(ctx, acc, fx)
				if err != nil {
					return errors.Wrap(err, "failed to resolve current subscription")
				}
				return nil
			})
		}
		if err != nil {
			result.Failed++
			s.log.ErrorContext(ctx, fmt.Sprintf("reconcile subscriptions of contractor %s failed: %v", contractorID, err))
			continue
		}
		result.Changed++
	}
	return result, nil
}
