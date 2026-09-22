package errors

import (
	"fmt"
	"strings"

	ozzo_validation "github.com/go-ozzo/ozzo-validation"
)

type ErrorCode uint

type ValidationErrors map[string]string

type CommonError struct {
	ClientMessage    string           `json:"message"`
	SystemMessage    interface{}      `json:"errorMessage"`
	ValidationErrors ValidationErrors `json:"validationErrors,omitempty"`
	ErrorCode        ErrorCode        `json:"code"`
	ErrorMessage     *string          `json:"-"`
	ErrorTrace       *string          `json:"-"`
}

func (err CommonError) Error() string {
	return fmt.Sprintf("CommonError: %+v. Trace: %+v", err.ErrorMessage, err.ErrorTrace)
}

func buildValidationError(err error) ValidationErrors {
	var errors ValidationErrors = map[string]string{}

	errValidate := strings.Split(err.Error(), ";")
	for _, err := range errValidate {
		errPerField := strings.Split(err, ":")
		if len(errPerField[0]) <= 1 {
			errors["error"] = errPerField[0]
		} else {
			errors[strings.TrimSpace(errPerField[0])] = strings.TrimSpace(errPerField[1])
		}
	}

	return errors
}

func NewError(errCode ErrorCode, err error) *CommonError {
	if _err, ok := err.(*CommonError); ok {

		return _err
	}

	var errMsg *string
	var errTrace *string
	var clientMessage string = "Unknown error."
	var systemMessage interface{} = "Unknown error."
	var commonError = errorCodes[errCode]

	if err != nil {

		s := err.Error()
		errMsg = &s

		ss := fmt.Sprintf("%+v", err)
		errTrace = &ss

		if errCode == UNKNOWN_ERROR {
			systemMessage = ss
		}
	}

	if commonError == nil {

		return &CommonError{
			ClientMessage: clientMessage,
			SystemMessage: systemMessage,
			ErrorCode:     errCode,
			ErrorTrace:    errTrace,
			ErrorMessage:  errMsg,
		}
	}

	// The registered message is the client-facing one, and Error() returns
	// ClientMessage, so that is what callers and HTTP responses surface. The
	// raw underlying error goes to ErrorMessage, which is tagged json:"-" and
	// therefore never leaves the process. Handlers pass driver errors (e.g. a
	// Postgres constraint violation) straight in here, so using the raw text
	// as the client message would leak internals and drop the curated text.
	//
	// A few codes are registered without a ClientMessage (UNAUTHORIZED), so
	// there is no curated text to send; fall back to the raw error rather
	// than returning an empty message to the client.
	//
	// errMsg is nil when err is nil, which is a legitimate call meaning "no
	// underlying cause". It is kept as a pointer so that nil survives instead
	// of being dereferenced.
	registeredMessage := commonError.ClientMessage
	if registeredMessage == "" && errMsg != nil {
		registeredMessage = *errMsg
	}

	return &CommonError{
		ClientMessage: registeredMessage,
		SystemMessage: commonError.SystemMessage,
		ErrorCode:     errCode,
		ErrorTrace:    errTrace,
		ErrorMessage:  errMsg,
	}
}

func (err *CommonError) SetClientMessage(message string) {
	err.ClientMessage = message
}

func (err *CommonError) SetSystemMessage(message interface{}) {
	err.SystemMessage = message
}

func (err *CommonError) SetValidationMessage(message interface{}) {
	if _err, ok := message.(ozzo_validation.Errors); ok {
		err.ValidationErrors = buildValidationError(_err)
	}
}
