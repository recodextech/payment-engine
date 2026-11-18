package errors

const (
	CodeResConflict  = 400001
	CodeResNotFound  = 400002
	CodeMalformedReq = 400003
	CodeValEncode    = 400008
	CodeResWFailed   = 500001
	CodeKeyEncode    = 400010
)

var messages = map[int]string{
	CodeResConflict:  "resource conflict",
	CodeResNotFound:  "resource not found",
	CodeMalformedReq: "malformed request",
	CodeResWFailed:   "resource write failed",
	CodeKeyEncode:    "key encode failed",
	CodeValEncode:    "value encode failed",
}

func Msg(code int) string {
	return messages[code]
}

type RepositoryDataNotExistError struct {
	error
}

// NotFound returns a RepositoryDataNotExistError that is safe to wrap and print.
func NotFound(what string) error {
	return RepositoryDataNotExistError{error: New(what + " not found")}
}

// RepositoryDuplicateError is returned when a row with the same unique value already exists.
type RepositoryDuplicateError struct {
	error
}

// Duplicate returns a RepositoryDuplicateError that is safe to wrap and print.
func Duplicate(what string) error {
	return RepositoryDuplicateError{error: New(what + " already exists")}
}
