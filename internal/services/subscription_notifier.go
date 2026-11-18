package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"payment-engine/internal/adaptors/clients"
	"payment-engine/internal/domain"
	"payment-engine/internal/domain/adaptors"
	"payment-engine/internal/domain/adaptors/repositories"
	"payment-engine/internal/domain/application"

	"github.com/recodextech/api-definitions/events"

	"github.com/recodextech/container"
)

const subscriptionScreen = "subscription"

// PushSender sends a push notification to a user's devices.
type PushSender interface {
	SendToUser(ctx context.Context, actorUserID, actorAccountID string, req clients.SendNotificationRequest) error
}

// SubscriptionNotifier sends subscription push notifications through notification-service (Firebase).
// Sending is best effort: failures are logged and never undo the subscription change.
type SubscriptionNotifier struct {
	push     PushSender
	accounts repositories.ContractorAccountRepository
	log      adaptors.Logger
}

func (n *SubscriptionNotifier) Init(c container.Container) error {
	n.push = c.Resolve(application.ModuleNotificationClient).(PushSender)
	n.accounts = c.Resolve(application.ModuleContractorAccountRepo).(repositories.ContractorAccountRepository)
	n.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(`services.subscription-notifier`))
	return nil
}

func (n *SubscriptionNotifier) Send(ctx context.Context, notices []notice) {
	for _, nt := range notices {
		acc, err := n.accounts.GetContractorAccount(ctx, nt.contractorID)
		if err != nil {
			n.log.ErrorContext(ctx, fmt.Sprintf("subscription notification %s: contractor %s: %v", nt.typ, nt.contractorID, err))
			continue
		}
		data := map[string]string{"screen": subscriptionScreen, "contractor_id": nt.contractorID}
		for k, v := range nt.data {
			data[k] = v
		}
		err = n.push.SendToUser(ctx, acc.UserID, acc.AccountID, clients.SendNotificationRequest{
			UserID: acc.UserID, Type: nt.typ, Title: nt.title, Description: nt.body, Data: data,
		})
		if err != nil {
			n.log.ErrorContext(ctx, fmt.Sprintf("subscription notification %s to user %s failed: %v", nt.typ, acc.UserID, err))
		}
	}
}

func formatDate(millis int64) string {
	return time.UnixMilli(millis).In(colombo).Format("2 Jan 2006")
}

var colombo = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		return time.UTC
	}
	return loc
}()

func rowData(r events.SubscriptionEvent) map[string]string {
	return map[string]string{
		"subscription_id": r.Payload.ID.String(),
		"plan_code":       r.Payload.Plan.Code,
		"remaining":       strconv.Itoa(r.Payload.RemainingJobs),
		"period_end":      strconv.FormatInt(r.Payload.PeriodEnd, 10),
	}
}

func activatedNotice(r events.SubscriptionEvent) notice {
	return notice{contractorID: r.Payload.ContractorID.String(), typ: domain.NotifySubscriptionActivated,
		title: r.Payload.Plan.Code + " plan active",
		body:  fmt.Sprintf("You can post %d jobs until %s.", r.Payload.RemainingJobs, formatDate(r.Payload.PeriodEnd)),
		data:  rowData(r)}
}

func renewedNotice(r events.SubscriptionEvent) notice {
	return notice{contractorID: r.Payload.ContractorID.String(), typ: domain.NotifySubscriptionRenewed,
		title: "Plan renewed",
		body:  fmt.Sprintf("Your %s plan renewed with %d jobs until %s.", r.Payload.Plan.Code, r.Payload.RemainingJobs, formatDate(r.Payload.PeriodEnd)),
		data:  rowData(r)}
}

func paymentReceivedNotice(plan events.SubscriptionEvent) notice {
	return notice{contractorID: plan.Payload.ContractorID.String(), typ: domain.NotifySubscriptionPaymentRecv,
		title: "Payment received",
		body:  fmt.Sprintf("Thanks! Your %s plan continues from %s.", plan.Payload.Plan.Code, formatDate(plan.Payload.PeriodEnd)),
		data:  rowData(plan)}
}

func expiredNotice(contractorID, planCode string) notice {
	return notice{contractorID: contractorID, typ: domain.NotifySubscriptionExpired,
		title: "Your plan has ended",
		body:  fmt.Sprintf("Your %s plan has ended and you're now on the Free plan.", planCode),
		data:  map[string]string{"plan_code": planCode}}
}

func addonNotice(r events.SubscriptionEvent, remaining int) notice {
	data := rowData(r)
	data["remaining"] = strconv.Itoa(remaining)
	return notice{contractorID: r.Payload.ContractorID.String(), typ: domain.NotifySubscriptionAddon,
		title: fmt.Sprintf("%d extra jobs added", r.Payload.Plan.JobLimit),
		body:  fmt.Sprintf("You can use them until %s. You have %d jobs left.", formatDate(r.Payload.PeriodEnd), remaining),
		data:  data}
}

func quotaData(c current, remaining int) map[string]string {
	data := rowData(c.Plan)
	data["remaining"] = strconv.Itoa(remaining)
	return data
}

func quotaLowNotice(contractorID string, c current) notice {
	return notice{contractorID: contractorID, typ: domain.NotifyQuotaLow,
		title: "1 job left",
		body:  fmt.Sprintf("Your plan renews on %s. Upgrade or buy more jobs to post more now.", formatDate(c.Plan.Payload.PeriodEnd)),
		data:  quotaData(c, 1)}
}

func quotaExhaustedNotice(contractorID string, c current) notice {
	return notice{contractorID: contractorID, typ: domain.NotifyQuotaExhausted,
		title: "Job limit reached",
		body:  fmt.Sprintf("You've used all your jobs until %s. Upgrade or buy more jobs to keep posting.", formatDate(c.Plan.Payload.PeriodEnd)),
		data:  quotaData(c, 0)}
}
