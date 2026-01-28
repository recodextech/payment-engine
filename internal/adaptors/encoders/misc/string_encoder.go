package misc

import "payment-engine/pkg/errors"

type StringEncoder struct{}

// Encode to encode strings
func (s StringEncoder) Encode(in interface{}) (out []byte, err error) {
	str, ok := in.(string)
	if !ok {
		return nil, errors.New("key encoder")
	}

	return []byte(str), nil
}

// Decode to decode strings
func (s StringEncoder) Decode(in []byte) (out interface{}, err error) {
	if in == nil {
		return nil, errors.New("key encoder")
	}

	return string(in), nil
}
