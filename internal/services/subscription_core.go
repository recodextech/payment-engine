package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"payment-engine/internal/adaptors/billing"
	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"
	"payment-engine/internal/domain/models"
	"payment-engine/pkg/errors"

	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

// Service error codes for subscriptions
const (
	CodeQuotaExceeded      = 402001
	CodeInvalidPlan        = 400101
	CodeInvalidReference   = 400102
	CodeAmountMismatch     = 400103
	CodeContractorNotFound = 404101
	CodeNotContractorOwner = 403101
	CodeInvalidRequest     = 400105
	CodeDuplicatePayment   = 409101

	// ErrorQuotaExceeded is the machine readable error the app and demand-handler check for.
	ErrorQuotaExceeded = "SUBSCRIPTION_QUOTA_EXCEEDED"
)

func serviceError(status, code int, msg string) ServiceError {
	return NewServiceError(status, code, msg)
}

// quotaExceeded is returned when a contractor has no job left in any active row.
func quotaExceeded() ServiceError {
	e := serviceError(http.StatusPaymentRequired, CodeQuotaExceeded, "subscription job quota exceeded")
	e.Details = map[string]any{"error": ErrorQuotaExceeded, "remaining": 0}
	return e
}

// notice is a push notification queued during a transaction and sent after it commits.
type notice struct {
	contractorID string
	typ          string
	title        string
	body         string
	data         map[string]string
}

// effects collects what must happen after a transaction commits.
type effects struct {
	notices []notice
}

func (e *effects) notify(n notice) {
	e.notices = append(e.notices, n)
}

// current is what a contractor can use now: the active monthly plan and the active add-on packs.
type current struct {
	Plan   events.SubscriptionEvent
	Addons []events.SubscriptionEvent
}

// remaining is every job the contractor can still post.
func (c current) remaining() int {
	n := c.Plan.Payload.RemainingJobs
	for _, a := range c.Addons {
		n += a.Payload.RemainingJobs
	}
	return n
}

// rows are the active rows, soonest ending first: the order jobs are taken in.
func (c current) rows() []events.SubscriptionEvent {
	rows := append([]events.SubscriptionEvent{c.Plan}, c.Addons...)
	for i := 1; i < len(rows); i++ {
		for j := i; j > 0 && rows[j].Payload.PeriodEnd < rows[j-1].Payload.PeriodEnd; j-- {
			rows[j], rows[j-1] = rows[j-1], rows[j]
		}
	}
	return rows
}

// SubscriptionCore holds the subscription rules shared by the subscription and quota services.
//
// Each monthly plan or add-on pack is a fix.subscriptions row that starts on the day it is paid for,
// lasts one month and has its own job count. When a monthly row ends, the next month starts from a
// payment made in advance, or the contractor falls back to FREE. Ended add-on rows just end.
type SubscriptionCore struct {
	subs     repositories.SubscriptionRepository
	plans    repositories.PlanRepository
	accounts repositories.ContractorAccountRepository
	conf     *billing.Config
	notifier *SubscriptionNotifier
	log      adaptors.Logger
	now      func() time.Time
}

func (c *SubscriptionCore) Init(di container.Container) error {
	c.subs = di.Resolve(application.ModuleSubscriptionRepo).(repositories.SubscriptionRepository)
	c.plans = di.Resolve(application.ModulePlanRepo).(repositories.PlanRepository)
	c.accounts = di.Resolve(application.ModuleContractorAccountRepo).(repositories.ContractorAccountRepository)
	c.conf = di.GetGlobalConfig(application.ModuleBillingConfig).(*billing.Config)
	c.notifier = di.Resolve(application.ModuleSubscriptionNotifier).(*SubscriptionNotifier)
	c.log = di.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`services.subscription-core`))
	c.now = func() time.Time { return time.Now().UTC() }
	return nil
}

// contractorAccount loads the contractor's account or answers 404.
func (c *SubscriptionCore) contractorAccount(ctx context.Context, contractorID string) (models.ContractorAccount, error) {
	acc, err := c.accounts.GetContractorAccount(ctx, contractorID)
	if err != nil {
		if isNotExist(err) {
			return acc, serviceError(http.StatusNotFound, CodeContractorNotFound, "contractor not found")
		}
		return acc, errors.Wrap(err, "failed to get contractor account")
	}
	return acc, nil
}

func isNotExist(err error) bool {
	var target errors.RepositoryDataNotExistError
	return asError(err, &target)
}

func isDuplicate(err error) bool {
	var target errors.RepositoryDuplicateError
	return asError(err, &target)
}

func (c *SubscriptionCore) reference(acc models.ContractorAccount) string {
	return c.conf.Reference(acc.AccountNumber)
}

// currencyOf is the account preference, defaulting to LKR.
func currencyOf(acc models.ContractorAccount) string {
	if acc.Currency == "" {
		return domain.DefaultCurrency
	}
	return strings.ToUpper(acc.Currency)
}

// newMeta builds the row meta. Calls without a user (the daily sweep) are recorded as the system user.
func newMeta(ctx context.Context, accountID string) events.Meta {
	meta := events.NewMeta()
	meta.Version = 1
	if id, err := uuid.Parse(accountID); err == nil {
		meta.AccountID = id
	}
	meta.CreatedBy = actorID(ctx)
	switch v := ctx.Value(events.ContextKeyTraceID.String()).(type) {
	case string:
		if id, err := uuid.Parse(v); err == nil && id != uuid.Nil {
			meta.TraceID = id
		}
	case uuid.UUID:
		if v != uuid.Nil {
			meta.TraceID = v
		}
	}
	return meta
}

// actorID is the calling user, or the system user for calls without one.
func actorID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(events.ContextKeyUserID.String()).(uuid.UUID); ok && id != uuid.Nil {
		return id
	}
	return domain.SystemUserID
}

func touchMeta(ctx context.Context, meta events.Meta) events.Meta {
	fresh := newMeta(ctx, meta.AccountID.String())
	fresh.Version = meta.Version + 1
	fresh.AccountID = meta.AccountID
	return fresh
}

func addMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// newRow builds an ACTIVE row for a plan or add-on pack starting at start, with the full job limit.
func (c *SubscriptionCore) newRow(ctx context.Context, acc models.ContractorAccount, plan events.SubscriptionPlanPayload,
	start time.Time, price float64, paymentID *uuid.UUID) events.SubscriptionEvent {
	row := events.SubscriptionEvent{EventMeta: newMeta(ctx, acc.AccountID)}
	p := &row.Payload
	p.ID = uuid.New()
	p.ContractorID = parseID(acc.ContractorID)
	p.PlanID = plan.ID
	p.Kind = plan.Kind
	p.Status = domain.SubscriptionActive
	p.PeriodStart = start.UnixMilli()
	p.PeriodEnd = addMonths(start, domain.PeriodMonths).UnixMilli()
	p.RemainingJobs = plan.JobLimit
	p.Plan = events.PlanSnapshot{Code: plan.Code, Price: price, Currency: currencyOf(acc), JobLimit: plan.JobLimit}
	p.PaymentID = paymentID
	return row
}

func (c *SubscriptionCore) startRow(ctx context.Context, acc models.ContractorAccount, plan events.SubscriptionPlanPayload,
	start time.Time, price float64, paymentID *uuid.UUID) (events.SubscriptionEvent, error) {
	row := c.newRow(ctx, acc, plan, start, price, paymentID)
	return row, c.subs.CreateRow(ctx, row)
}

func (c *SubscriptionCore) endRow(ctx context.Context, row events.SubscriptionEvent, reason string) error {
	row.Payload.Status = domain.SubscriptionEnded
	row.Payload.EndReason = reason
	row.EventMeta = touchMeta(ctx, row.EventMeta)
	return c.subs.EndRow(ctx, row)
}

// startFree starts a FREE month at start, moved on by whole months until it contains now.
func (c *SubscriptionCore) startFree(ctx context.Context, acc models.ContractorAccount, start time.Time) (events.SubscriptionEvent, error) {
	plan, err := c.plans.GetPlanByCode(ctx, domain.FreePlanCode)
	if err != nil {
		return events.SubscriptionEvent{}, errors.Wrap(err, "free plan is not configured")
	}
	now := c.now()
	for !now.Before(addMonths(start, domain.PeriodMonths)) {
		start = addMonths(start, domain.PeriodMonths)
	}
	return c.startRow(ctx, acc, plan, start, 0, nil)
}

func (c *SubscriptionCore) markApplied(ctx context.Context, pay *events.SubscriptionPaymentEvent) error {
	pay.Payload.Status = domain.SubPaymentApplied
	pay.EventMeta = touchMeta(ctx, pay.EventMeta)
	return c.subs.SetPaymentStatus(ctx, *pay)
}

// nextPaidMonth opens the month after the one ending at end from the oldest payment made in advance.
func (c *SubscriptionCore) nextPaidMonth(ctx context.Context, acc models.ContractorAccount, end time.Time) (events.SubscriptionEvent, bool, error) {
	paid, err := c.subs.ListPaidPayments(ctx, acc.ContractorID)
	if err != nil {
		return events.SubscriptionEvent{}, false, errors.Wrap(err, "failed to list paid payments")
	}
	for _, pay := range paid {
		plan, err := c.plans.GetPlanByCode(ctx, pay.Payload.PlanCode)
		if err != nil {
			if isNotExist(err) {
				continue
			}
			return events.SubscriptionEvent{}, false, errors.Wrap(err, "failed to get plan")
		}
		if plan.Kind != domain.PlanKindSubscription {
			continue
		}
		id := pay.Payload.ID
		row, err := c.startRow(ctx, acc, plan, end, pay.Payload.Amount, &id)
		if err != nil {
			return row, false, errors.Wrap(err, "failed to start paid month")
		}
		return row, true, c.markApplied(ctx, &pay)
	}
	return events.SubscriptionEvent{}, false, nil
}

// resolveCurrent returns what the contractor can use now, ending rows whose month is over:
//   - an ended monthly row is followed by a month paid in advance, or by FREE
//   - an ended add-on row just ends
//
// It must run inside a transaction after LockContractor.
func (c *SubscriptionCore) resolveCurrent(ctx context.Context, acc models.ContractorAccount, fx *effects) (current, error) {
	now := c.now()
	for {
		rows, err := c.subs.ListActiveRows(ctx, acc.ContractorID)
		if err != nil {
			return current{}, errors.Wrap(err, "failed to list active subscriptions")
		}
		var cur current
		hasPlan, changed := false, false
		for _, row := range rows {
			end := time.UnixMilli(row.Payload.PeriodEnd).UTC()
			if end.After(now) {
				if row.Payload.Kind == domain.PlanKindSubscription {
					cur.Plan, hasPlan = row, true
				} else {
					cur.Addons = append(cur.Addons, row)
				}
				continue
			}
			changed = true
			if err := c.endRow(ctx, row, domain.EndExpired); err != nil {
				return current{}, errors.Wrap(err, "failed to end expired subscription")
			}
			if row.Payload.Kind != domain.PlanKindSubscription {
				continue
			}
			next, ok, err := c.nextPaidMonth(ctx, acc, end)
			if err != nil {
				return current{}, errors.Wrap(err, "failed to open next paid month")
			}
			switch {
			case ok:
				fx.notify(renewedNotice(next))
			case row.Payload.Plan.Code != domain.FreePlanCode:
				fx.notify(expiredNotice(acc.ContractorID, row.Payload.Plan.Code))
				fallthrough
			default:
				if _, err := c.startFree(ctx, acc, end); err != nil {
					return current{}, errors.Wrap(err, "failed to start free plan")
				}
			}
		}
		if changed {
			continue // a month that started from an old end date may itself be over already
		}
		if !hasPlan {
			plan, err := c.startFree(ctx, acc, now)
			if err != nil {
				return current{}, errors.Wrap(err, "failed to start free plan")
			}
			cur.Plan = plan
		}
		return cur, nil
	}
}

// withContractor runs fn in a transaction holding the contractor lock, then sends queued notifications
// once the transaction has committed.
func (c *SubscriptionCore) withContractor(ctx context.Context, acc models.ContractorAccount,
	fn func(ctx context.Context, fx *effects) error) error {
	fx := &effects{}
	err := c.subs.InTransaction(ctx, func(ctx context.Context) error {
		if err := c.subs.LockContractor(ctx, acc.ContractorID); err != nil {
			return errors.Wrap(err, "failed to lock contractor")
		}
		return fn(ctx, fx)
	})
	if err != nil {
		return errors.Wrap(err, "contractor transaction failed")
	}
	c.notifier.Send(ctx, fx.notices)
	return nil
}

func parseID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
