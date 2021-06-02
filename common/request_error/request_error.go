package requesterror

import "errors"

const (
	ErrorTypeNone     = iota
	ErrorTypeInternal = iota
	ErrorTypeClient   = iota
)

// RequestError is an error with an added type field to determine how it should be handled
type RequestError struct {
	error
	Type int
}

// NoError returns a RequestError with type ErrorTypeNone
func NoError() RequestError {
	return RequestError{
		error: nil,
		Type:  ErrorTypeNone,
	}
}

// InternalError returns a RequestError with type ErrorTypeInternal and an internal error message
func InternalError() RequestError {
	return RequestError{
		error: errors.New("an internal error occurred"),
		Type:  ErrorTypeInternal,
	}
}

// ClientError returns a RequestError with type ErrorTypeClient and the provided message
func ClientError(message string) RequestError {
	return RequestError{
		error: errors.New(message),
		Type:  ErrorTypeClient,
	}
}

// OAuthRequestError is a RequestError with an added error name field to determine the oauth error name
type OAuthRequestError struct {
	RequestError
	ErrorName string
}

// OAuthNoError returns an OAuthRequestError with type ErrorTypeNone
func OAuthNoError() OAuthRequestError {
	return OAuthRequestError{
		RequestError: NoError(),
		ErrorName:    "",
	}
}

// OAuthInternalError returns an OAuthRequestError with type ErrorTypeInternal and an internal error message
func OAuthInternalError() OAuthRequestError {
	return OAuthRequestError{
		RequestError: InternalError(),
		ErrorName:    "",
	}
}

// ClientError returns a RequestError with type ErrorTypeClient and the provided error name and message
func OAuthClientError(errorName string, message string) OAuthRequestError {
	return OAuthRequestError{
		RequestError: ClientError(message),
		ErrorName:    errorName,
	}
}
