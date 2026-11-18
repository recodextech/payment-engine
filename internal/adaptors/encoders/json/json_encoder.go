package json

import (
	"encoding/json"
	"payment-engine/pkg/errors"
)

type Encoder struct{}

func (e Encoder) Encode(in interface{}) (out []byte, err error) {
	return json.Marshal(in)
}

func (e Encoder) Decode(data []byte) (interface{}, error) {
	out := Encoder{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode json")
	}
	return out, nil
}
