package request

import "errors"

var (
	ErrBadReq           = errors.New("errors: bad requset")
	ErrValidateDataFail = errors.New("errors: validate data failed")
)
