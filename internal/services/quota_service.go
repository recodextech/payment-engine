package services

import (
	"context"
	"net/http"
	"payment-engine/pkg/errors"

	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/application"

	"payment-engine/internal/domain/models"

	"github.com/recodextech/api-definitions/events"
	"github.com/recodextech/api-definitions/pkg/uuid"
	"github.com/recodextech/container"
)

// QuotaService takes and returns job slots. Creating a process takes one job, however many jobs the
// process has, from the contractor's active row that ends soonest.
type QuotaService struct {
	core *SubscriptionCore
	log  adaptors.Logger
}

func (s *QuotaService) Init(c container.Container) error {
	s.core = c.Resolve(application.ModuleSubscriptionCore).(*SubscriptionCore)
	s.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`services.quota`))
	return nil
}

// Consume takes the process's slot. A process that already holds one takes nothing, so retries and
// rescheduling a job of a process that still holds its slot are safe; a RELEASED slot is taken again.
// When no active row has a job left, nothing is taken and a 402 SUBSCRIPTION_QUOTA_EXCEEDED is returned.
func (s *QuotaService) Consume(ctx context.Context, contractorID, processID string) (models.QuotaResult, error) {
	process := parseID(processID)
	if process == uuid.Nil {
		return models.QuotaResult{}, serviceError(http.StatusBadRequest, CodeInvalidRequest, "process_id is required")
	}
	acc, err := s.core.contractorAccount(ctx, contractorID)
	if err != nil {
		return models.QuotaResult{}, errors.Wrap(err, "failed to get contractor account")
	}
	var result models.QuotaResult
	err = s.core.withContractor(ctx, acc, func(ctx context.Context, fx *effects) error {
		cur, err := s.core.resolveCurrent(ctx, acc, fx)
		if err != nil {
			return errors.Wrap(err, "failed to resolve current subscription")
		}
		entry, found, err := s.core.subs.GetLedgerByProcessForUpdate(ctx, processID)
		if err != nil {
			return errors.Wrap(err, "failed to get quota ledger")
		}
		if found && entry.Payload.Status == domain.QuotaConsumed {
			result = models.QuotaResult{Remaining: cur.remaining()}
			return nil
		}
		for _, row := range cur.rows() {
			remaining, ok, err := s.core.subs.TakeJob(ctx, row.Payload.ID.String())
			if err != nil {
				return errors.Wrap(err, "failed to take a job")
			}
			if !ok {
				continue
			}
			if !found {
				entry = events.QuotaLedgerEvent{EventMeta: newMeta(ctx, acc.AccountID)}
				entry.Payload.ID = uuid.New()
				entry.Payload.ProcessID = process
			} else {
				entry.EventMeta = touchMeta(ctx, entry.EventMeta)
			}
			entry.Payload.Status = domain.QuotaConsumed
			entry.Payload.SubscriptionID = row.Payload.ID
			entry.Payload.Reason = ""
			if err := s.core.subs.SaveLedger(ctx, entry); err != nil {
				return errors.Wrap(err, "failed to save quota ledger")
			}
			total := cur.remaining() - row.Payload.RemainingJobs + remaining
			s.quotaWarnings(acc.ContractorID, cur, total, fx)
			result = models.QuotaResult{Consumed: 1, Remaining: total}
			return nil
		}
		return quotaExceeded()
	})
	if err != nil {
		return models.QuotaResult{}, errors.Wrap(err, "failed to consume quota")
	}
	return result, nil
}

// quotaWarnings queues the push when the last job left is reached and when none are left.
func (s *QuotaService) quotaWarnings(contractorID string, cur current, total int, fx *effects) {
	switch total {
	case 0:
		fx.notify(quotaExhaustedNotice(contractorID, cur))
	case 1:
		fx.notify(quotaLowNotice(contractorID, cur))
	}
}

// Release gives a process's slot back when it never ran. The job returns to the row it came from while
// that row is still active; otherwise the release is recorded with no refund. Releasing twice, or a
// process that never held a slot, does nothing.
func (s *QuotaService) Release(ctx context.Context, processID, reason string) (models.QuotaResult, error) {
	err := s.core.subs.InTransaction(ctx, func(ctx context.Context) error {
		entry, ok, err := s.core.subs.GetLedgerByProcessForUpdate(ctx, processID)
		if err != nil {
			return errors.Wrap(err, "failed to get quota ledger")
		}
		if !ok || entry.Payload.Status != domain.QuotaConsumed {
			return nil
		}
		entry.Payload.Status = domain.QuotaReleased
		entry.Payload.Reason = reason
		entry.EventMeta = touchMeta(ctx, entry.EventMeta)
		if err := s.core.subs.SaveLedger(ctx, entry); err != nil {
			return errors.Wrap(err, "failed to save quota ledger")
		}
		_, err = s.core.subs.ReturnJob(ctx, entry.Payload.SubscriptionID.String(), s.core.now().UnixMilli())
		if err != nil {
			return errors.Wrap(err, "failed to return a job")
		}
		return nil
	})
	if err != nil {
		return models.QuotaResult{}, errors.Wrap(err, "failed to release quota")
	}
	return models.QuotaResult{}, nil
}
