package commons

import "errors"

var (
	ErrorMismatchConfiguration = errors.New("mismatch_configuration_error")
	ErrorDriverNotFound        = errors.New("driver_not_found_error")
)
