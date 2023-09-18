package errors

import "fmt"

func (err NonOkError) Error() string {
	return fmt.Sprintf(`fetching '%s' returned non ok status code '%d'
URL: %s
Method: %s
BODY: %s`,
		err.Operation,
		err.Code,
		err.Response.Request.Method,
		err.Response.Request.URL,
		err.Response.Body())
}

func (err UnknownPlatformError) Error() string {
	return fmt.Sprintf("unknown platform, either insighter doesnot support the platform currently or platform doesnot exists: '%s'", err.Platform)
}
