package vmapi

import "errors"

type ErrorCode string

const (
	CodeInvalidArgument ErrorCode = "invalid_argument"
	CodeNotFound        ErrorCode = "not_found"
	CodeConflict        ErrorCode = "conflict"
	CodeUnsupported     ErrorCode = "unsupported"
	CodeUnavailable     ErrorCode = "unavailable"
	CodeUnknownOutcome  ErrorCode = "unknown_outcome"
	CodeCursorExpired   ErrorCode = "cursor_expired"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrUnsupported     = errors.New("unsupported")
	ErrUnavailable     = errors.New("unavailable")
	ErrUnknownOutcome  = errors.New("unknown outcome")
	ErrCursorExpired   = errors.New("cursor expired")
)

// Error holds safe correlation details for errors.As and supports errors.Is
// through wrapping. Message must already be lowercase, sanitized and bounded by
// the backend. No request, credential or raw backend cause is retained or echoed.
type Error struct {
	Code      ErrorCode
	RequestID RequestID
	Message   string
}

func (e *Error) Error() string {
	if e == nil {
		return "vmapi error"
	}
	if e.Message != "" {
		return e.Message
	}
	if s := e.sentinel(); s != nil {
		return s.Error()
	}
	return "vmapi error"
}
func (e *Error) Is(target error) bool {
	return e != nil && e.sentinel() != nil && target == e.sentinel()
}
func (e *Error) sentinel() error {
	switch e.Code {
	case CodeInvalidArgument:
		return ErrInvalidArgument
	case CodeNotFound:
		return ErrNotFound
	case CodeConflict:
		return ErrConflict
	case CodeUnsupported:
		return ErrUnsupported
	case CodeUnavailable:
		return ErrUnavailable
	case CodeUnknownOutcome:
		return ErrUnknownOutcome
	case CodeCursorExpired:
		return ErrCursorExpired
	}
	return nil
}
