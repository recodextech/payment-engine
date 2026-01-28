package json

import (
	"encoding/json"
)

type Encoder struct{}

func (e Encoder) Encode(in interface{}) (out []byte, err error) {
	return json.Marshal(in)
}

func (e Encoder) Decode(data []byte) (interface{}, error) {
	out := Encoder{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
