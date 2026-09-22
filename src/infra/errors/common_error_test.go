package errors

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/stretchr/testify/assert"
)

func TestOnlyNewClientErrorMessage(t *testing.T) {
	errMsg := NewError(0, nil)
	errMsg.SetClientMessage("This is a new client message")

	if assert.NotNil(t, errMsg.ClientMessage) {
		assert.Equal(t, "This is a new client message", errMsg.ClientMessage, "Client message should be a new client message")
	}

	if assert.NotNil(t, errMsg.SystemMessage) {
		assert.Equal(t, errMsg.SystemMessage, "Unknown error.")
	}
}

func TestOnlyNewSystemErrorMessage(t *testing.T) {
	errMsg := NewError(0, nil)
	errMsg.SetSystemMessage("This is a new system message")

	if assert.NotNil(t, errMsg.ClientMessage) {
		assert.Equal(t, errMsg.ClientMessage, "Unknown error.")
	}

	if assert.NotNil(t, errMsg.SystemMessage) {
		assert.Equal(t, "This is a new system message", errMsg.SystemMessage, "System message should be a new system message")
	}
}

func TestNewErrorMessage(t *testing.T) {
	errMsg := NewError(0, nil)
	errMsg.SetClientMessage("This is a new client message")
	errMsg.SetSystemMessage("This is a new system message")

	if assert.NotNil(t, errMsg.ClientMessage) {
		assert.Equal(t, "This is a new client message", errMsg.ClientMessage, "Client message should be a new client message")
	}

	if assert.NotNil(t, errMsg.SystemMessage) {
		assert.Equal(t, "This is a new system message", errMsg.SystemMessage, "System message should be a new system message")
	}
}

func TestCommonError(t *testing.T) {
	errMsg := NewError(0, nil)

	if assert.NotNil(t, errMsg.ClientMessage) {
		assert.Equal(t, "Unknown error.", errMsg.ClientMessage, "Client message should be a new client message")
	}

	if assert.NotNil(t, errMsg.SystemMessage) {
		assert.Equal(t, "Unknown error.", errMsg.SystemMessage, "System message should be a new system message")
	}
}

func TestThrownError(t *testing.T) {
	errMsg := NewError(0, errors.New("this is another error"))

	if assert.NotNil(t, errMsg.ClientMessage) {
		assert.Equal(t, "Unknown error.", errMsg.ClientMessage, "Client message should be a new client message")
	}

	if assert.NotNil(t, errMsg.SystemMessage) {
		assert.Equal(t, "Unknown error.", errMsg.SystemMessage, "System message should be a new system message")
	}

	if assert.NotNil(t, errMsg.ErrorMessage) {
		assert.Equal(t, "this is another error", *errMsg.ErrorMessage, "System message should be a new system message")
	}
}

func TestReThrowCommonError(t *testing.T) {
	errorCodes[1] = &CommonError{
		ClientMessage: "This is my client message",
		SystemMessage: "This is my system message",
		ErrorCode:     1,
	}
	errorCodes[2] = &CommonError{
		ClientMessage: "This is my second client message",
		SystemMessage: "This is my second system message",
		ErrorCode:     2,
	}

	errMsg := NewError(1, errors.New("this is an error"))
	errMsg2 := NewError(2, errMsg)

	if assert.NotNil(t, errMsg2.ClientMessage) {
		assert.Equal(t, "This is my client message", errMsg2.ClientMessage, "Client message should be \"This is my client message\"")
	}

	if assert.NotNil(t, errMsg2.SystemMessage) {
		assert.Equal(t, "This is my system message", errMsg2.SystemMessage, "System message should be \"This is my system message\"")
	}
}

func TestValidationErrors(t *testing.T) {
	type testStruct struct {
		name string
		age  int
	}

	value := &testStruct{name: "John", age: 20}
	errorValidation := validation.ValidateStruct(
		value,
		validation.Field(&value.name, validation.Required),
		validation.Field(&value.age, validation.Min(25)),
	)

	errMsg := NewError(0, errorValidation)
	errMsg.SetValidationMessage(errorValidation)

	assert.NotNil(t, errMsg.ValidationErrors["age"])
}

func TestErrorString(t *testing.T) {
	commonError := NewError(0, errors.New("this is internal error"))
	errMsg := commonError.Error()

	if assert.NotNil(t, errMsg) {
		assert.Contains(t, errMsg, "CommonError", "Trace")
	}
}

// NewError is called with a nil error to raise a registered status error that
// has no underlying cause. Before the fix this dereferenced the nil error
// message and panicked, so the panic is pinned here for both a code that is
// registered with a client message and one that is not.
func TestNewErrorWithNilErrorDoesNotPanic(t *testing.T) {
	for _, errCode := range []ErrorCode{DATA_INVALID, UNAUTHORIZED, UNKNOWN_ERROR} {
		errMsg := NewError(errCode, nil)

		if assert.NotNil(t, errMsg) {
			assert.Equal(t, errCode, errMsg.ErrorCode)
			assert.Nil(t, errMsg.ErrorMessage, "there is no underlying error to record")
		}
	}
}

// The registered client message is what reaches the client. The raw error is
// internal: ErrorMessage is tagged json:"-", and the handlers pass driver
// errors such as Postgres constraint violations straight in, so the raw text
// must not become the client-facing message.
func TestClientMessageIsTheRegisteredOneNotTheRawError(t *testing.T) {
	raw := `pq: duplicate key value violates unique constraint "users_email_key"`

	errMsg := NewError(DATA_INVALID, errors.New(raw))

	assert.Equal(t, "Invalid Data Request", errMsg.ClientMessage)
	assert.NotContains(t, errMsg.ClientMessage, "pq:", "the raw driver error must not be the client message")

	if assert.NotNil(t, errMsg.ErrorMessage) {
		assert.Equal(t, raw, *errMsg.ErrorMessage, "the raw error is still kept for logging")
	}
}

// Error() is what the HTTP layer serialises, so the leak is observable there.
func TestHttpErrorDoesNotLeakTheRawError(t *testing.T) {
	raw := `pq: password authentication failed for user "wallet"`

	httpErr := NewError(DATA_INVALID, errors.New(raw)).ToHttpError()

	assert.Equal(t, "Invalid Data Request", httpErr.Error())
	assert.NotContains(t, httpErr.Error(), "password authentication failed")
}
