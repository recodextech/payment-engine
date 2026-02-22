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
