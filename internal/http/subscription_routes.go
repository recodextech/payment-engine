package http

import (
	"fmt"
	"net/http"

	"payment-engine/internal/http/handlers"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/http/responses/writers"
	"payment-engine/internal/http/validator"

	"github.com/recodextech/container"

	"github.com/recodextech/krouter"
)

// initSubscriptionRoutes registers the subscription routes:
//   - app routes need account-id, and contractor routes check the contractor belongs to it
//   - admin routes are only exposed through the gateway's admin-portal realm
//   - /internal routes are called by demand-handler and must never be exposed on the gateway
func initSubscriptionRoutes(r *Router, c container.Container) {
	subscription := c.Resolve(handlers.ModuleSubscription).(*handlers.SubscriptionHandler)

	// ---- app ----

	// List active plans priced in the account's currency
	r.router.Handle("/subscription/plans",
		r.krouter.NewHandler(
			"subscription.plans",
			request.JSONResponse{},
			subscription.ListPlans,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// Get the contractor's current subscription
	r.router.Handle(fmt.Sprintf("/contractor/{%s}/subscription", request.PathParamContractorID.String()),
		r.krouter.NewHandler(
			"subscription.get",
			request.JSONResponse{},
			subscription.GetSubscription,
			krouter.HandlerWithParameter(request.PathParamContractorID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleContractorOwner).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// Bank transfer details for paying the contractor's subscription
	r.router.Handle(fmt.Sprintf("/contractor/{%s}/subscription/bank-transfer", request.PathParamContractorID.String()),
		r.krouter.NewHandler(
			"subscription.bank-transfer",
			request.JSONResponse{},
			subscription.BankTransfer,
			krouter.HandlerWithParameter(request.PathParamContractorID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleAccountIDValidator).(krouter.Validator)),
			krouter.HandlerWithValidator(c.Resolve(validator.ModuleContractorOwner).(krouter.Validator)),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// ---- admin ----

	// List all plans
	r.router.Handle("/admin/subscription/plans",
		r.krouter.NewHandler(
			"admin.plans.list",
			request.JSONResponse{},
			subscription.AdminListPlans,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// Create a plan
	r.router.Handle("/admin/subscription/plans",
		r.krouter.NewHandler(
			"admin.plans.create",
			request.PlanRequest{},
			subscription.AdminCreatePlan,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Update a plan
	r.router.Handle(fmt.Sprintf("/admin/subscription/plans/{%s}", request.PathParamPlanID.String()),
		r.krouter.NewHandler(
			"admin.plans.update",
			request.PlanRequest{},
			subscription.AdminUpdatePlan,
			krouter.HandlerWithParameter(request.PathParamPlanID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPut)

	// List subscriptions
	r.router.Handle("/admin/subscriptions",
		r.krouter.NewHandler(
			"admin.subscriptions.list",
			request.JSONResponse{},
			subscription.AdminListSubscriptions,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// Grant a subscription to a contractor
	r.router.Handle(fmt.Sprintf("/admin/contractor/{%s}/subscription/grant", request.PathParamContractorID.String()),
		r.krouter.NewHandler(
			"admin.subscription.grant",
			request.GrantRequest{},
			subscription.AdminGrant,
			krouter.HandlerWithParameter(request.PathParamContractorID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// List a contractor's subscription payments
	r.router.Handle(fmt.Sprintf("/admin/contractor/{%s}/subscription/payments", request.PathParamContractorID.String()),
		r.krouter.NewHandler(
			"admin.subscription.payments.list",
			request.JSONResponse{},
			subscription.AdminListPayments,
			krouter.HandlerWithParameter(request.PathParamContractorID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)

	// Record a subscription payment
	r.router.Handle("/admin/subscription/payments",
		r.krouter.NewHandler(
			"admin.subscription.payment",
			request.PaymentRequest{},
			subscription.AdminRecordPayment,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// ---- internal ----

	// Consume one job from the contractor's quota
	r.router.Handle("/internal/quota/consume",
		r.krouter.NewHandler(
			"internal.quota.consume",
			request.QuotaConsumeRequest{},
			subscription.ConsumeQuota,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Release a previously consumed job back to the quota
	r.router.Handle("/internal/quota/release",
		r.krouter.NewHandler(
			"internal.quota.release",
			request.QuotaReleaseRequest{},
			subscription.ReleaseQuota,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)

	// Roll over due subscriptions
	r.router.Handle("/internal/subscriptions/reconcile",
		r.krouter.NewHandler(
			"internal.subscriptions.reconcile",
			request.JSONResponse{},
			subscription.Reconcile,
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJSONWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodPost)
}
