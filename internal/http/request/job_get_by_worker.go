package request

import "encoding/json"

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
	return req, err
}
