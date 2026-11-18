package services

import stderrors "errors"

func asError(err error, target any) bool {
	return stderrors.As(err, target)
}
