package request

import (
	"encoding/json"

	"payment-engine/pkg/errors"
)

// Path params and context keys used by the subscription routes
const (
	PathParamPlanID ContextParam = "plan_id"
	QueryParams     ContextParam = "query-params"
)

func decodeJSON[T any](data []byte) (interface{}, error) {
	var req T
	if len(data) == 0 {
		return req, nil
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}
	return req, nil
}

func encodeJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// PaymentRequest records a contractor's bank transfer from the admin portal: for a contractor picked
// in the portal, or found by the payment reference written on the transfer.
type PaymentRequest struct {
	ContractorID  string  `json:"contractor_id"`
	Reference     string  `json:"reference"` // contractor payment reference on the transfer, e.g. AA1111111
	PlanCode      string  `json:"plan_code"`
	Amount        float64 `json:"amount"` // optional, checked against the plan price
	BankReference string  `json:"bank_reference"`
	Note          string  `json:"note"`
}

func (PaymentRequest) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (PaymentRequest) Decode(data []byte) (interface{}, error) {
	return decodeJSON[PaymentRequest](data)
}

// PlanRequest creates or updates a plan from the admin portal.
type PlanRequest struct {
	Code     string  `json:"code"`
	Kind     string  `json:"kind"`
	Status   string  `json:"status"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	JobLimit int     `json:"job_limit"`
}

func (PlanRequest) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (PlanRequest) Decode(data []byte) (interface{}, error) { return decodeJSON[PlanRequest](data) }

// GrantRequest gives a contractor a monthly plan or an add-on pack without payment.
type GrantRequest struct {
	PlanCode string `json:"plan_code"`
	Reason   string `json:"reason"`
}

func (GrantRequest) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (GrantRequest) Decode(data []byte) (interface{}, error) { return decodeJSON[GrantRequest](data) }

// QuotaConsumeRequest takes one job slot for a process (internal, from demand-handler).
type QuotaConsumeRequest struct {
	ContractorID string `json:"contractor_id"`
	ProcessID    string `json:"process_id"`
}

func (QuotaConsumeRequest) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (QuotaConsumeRequest) Decode(data []byte) (interface{}, error) {
	return decodeJSON[QuotaConsumeRequest](data)
}

// QuotaReleaseRequest returns the slot of a process that never ran (internal, from demand-handler).
type QuotaReleaseRequest struct {
	ProcessID string `json:"process_id"`
	Reason    string `json:"reason"`
}

func (QuotaReleaseRequest) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (QuotaReleaseRequest) Decode(data []byte) (interface{}, error) {
	return decodeJSON[QuotaReleaseRequest](data)
}

// JSONResponse is a pass-through encoder for GET routes, so the response is logged as JSON.
type JSONResponse struct{}

func (JSONResponse) Encode(data interface{}) ([]byte, error) { return encodeJSON(data) }
func (JSONResponse) Decode(_ []byte) (interface{}, error)    { return JSONResponse{}, nil }
