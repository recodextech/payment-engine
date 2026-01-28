package json

type ValueEncodingError struct {
	error
}

const (
	failedToEncodeValue = `failed to encode value`
	failedToDecodeValue = `failed to decode value`
)
