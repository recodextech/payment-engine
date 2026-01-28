package writers

type ResponseWriterError struct {
	error
}

var (
	errorWritingResponse = "error populating information"
)
