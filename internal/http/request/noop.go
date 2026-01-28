package request

// NoOp is a no operation request DTO
type NoOp struct{}

// Encode encodes the request data to JSON
func (n NoOp) Encode(data interface{}) ([]byte, error) {
	return []byte{}, nil
}

// Decode decodes the request data from JSON
func (n NoOp) Decode(data []byte) (interface{}, error) {
	return NoOp{}, nil
}
