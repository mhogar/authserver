package requesterror

import "errors"

const (
	ErrorTypeNone     = iota
	ErrorTypeInternal = iota
	ErrorTypeClient   = iota
)

type RequestError struct {
	error
	Type int
}

func NoError() RequestError {
	return RequestError{
		error: nil,
		Type:  ErrorTypeNone,
	}
}

func InternalError() RequestError {
	return RequestError{
		error: errors.New("an internal error occurred"),
		Type:  ErrorTypeInternal,
	}
}

func ClientError(message string) RequestError {
	return RequestError{
		error: errors.New(message),
		Type:  ErrorTypeClient,
	}
}

type OAuthRequestError struct {
	RequestError
	ErrorName string
}

func OAuthNoError() OAuthRequestError {
	return OAuthRequestError{
		RequestError: NoError(),
		ErrorName:    "",
	}
}

func OAuthInternalError() OAuthRequestError {
	return OAuthRequestError{
		RequestError: InternalError(),
		ErrorName:    "",
	}
}

func OAuthClientError(errorName string, message string) OAuthRequestError {
	return OAuthRequestError{
		RequestError: ClientError(message),
		ErrorName:    errorName,
	}
}
