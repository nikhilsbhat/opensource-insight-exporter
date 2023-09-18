package errors

import (
	"github.com/go-resty/resty/v2"
)

type NonOkError struct {
	Code      int
	Operation string
	Response  *resty.Response
}

type UnknownPlatformError struct {
	Platform string
}
