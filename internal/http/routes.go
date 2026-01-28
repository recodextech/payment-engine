package http

import (
	"fmt"
	"net/http"

	"payment-engine/internal/http/handlers"
	"payment-engine/internal/http/request"
	"payment-engine/internal/http/responses"
	"payment-engine/internal/http/responses/writers"

	"github.com/recodextech/container"

	"github.com/recodextech/krouter"
)

func initRoutes(r *Router, c container.Container) {
	r.router.Handle(fmt.Sprintf("/jobs/worker/{%s}", request.PathParamWorkerID.String()),
		r.krouter.NewHandler(
			handlers.ModuleGetJobsByWorker,
			request.NoOp{},
			c.Resolve(handlers.ModuleGetJobsByWorker).(*handlers.GetJobsByWorkerHandler).Handle,
			krouter.HandlerWithParameter(request.PathParamWorkerID.String(), request.ParamTypeAppUUID),
			krouter.HandlerWithHeader(request.HeaderAccountID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderUserID.String(), request.ParamTypeAppUUID, nil),
			krouter.HandlerWithHeader(request.HeaderTraceID.String(), request.ParamTypeAppUUID, traceIDGenerate),
			krouter.HandlerWithSuccessHandlerFunc(c.Resolve(writers.ModuleJobListWriter).(responses.GenerateResponse).Response),
		)).Methods(http.MethodGet)
}
