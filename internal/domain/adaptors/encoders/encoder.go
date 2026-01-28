package encoders

type Encoder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}

type EncoderJSON interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte, out interface{}) error
}

/*
Builder returns a function of return type Encoder
*/
type Builder interface {
	Encoder() func() Encoder
}
