package users

import "github.com/pkg/errors"

// coded allows any level of our system to expose an HTTP error code.
// Once the error bubbles up to our HTTP layer, we can extract out the code
// to determine which to send.

type coded interface {
	Code() int
}

// getCode is the function to use to resolve our code.
func getCode(err error) int {
	err = errors.Cause(err)

	if err == nil {
		return 200
	}

	var code = 500

	if c, ok := err.(coded); ok {
		code = c.Code()
	}

	return code
}
