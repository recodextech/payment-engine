package request

import (
	"encoding/json"

	"payment-engine/pkg/errors"
)

type GetJobsByWorker struct {
	WorkerID  string `json:"worker_id" validate:"required,uuid"`
	AccountID string `json:"-"`
}

func (g GetJobsByWorker) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (g GetJobsByWorker) Decode(data []byte) (interface{}, error) {
	req := GetJobsByWorker{}
	err := json.Unmarshal(data, &req)
	if err != nil {
		return req, errors.Wrap(err, "failed to decode get jobs by worker request")
	}
	return req, nil
}
