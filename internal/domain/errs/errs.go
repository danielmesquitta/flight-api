package errs

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
)

type Err struct {
	Message    string
	StackTrace string
	Code       Code
	Errors     []ErrorItem
}

type ErrorItem struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type Code string

const (
	ErrCodeUnknown      Code = "unknown"
	ErrCodeNotFound     Code = "not_found"
	ErrCodeUnauthorized Code = "unauthorized"
	ErrCodeForbidden    Code = "forbidden"
	ErrCodeValidation   Code = "validation_error"
)

// NewErr creates a new Err instance from either an error or a string,
// and sets the Code flag to unknown. This is useful when you want to
// create an error that is not expected to happen, and you want to
// log it with stack tracing.
func New(err any, codes ...Code) *Err {
	if len(codes) > 0 {
		return newErr(err, codes[0])
	}

	return newErr(err, ErrCodeUnknown)
}

func newErr(err any, c Code) *Err {
	if err == nil {
		return nil
	}
	switch v := err.(type) {
	case *Err:
		return v
	case error:
		return &Err{
			Message:    v.Error(),
			StackTrace: string(debug.Stack()),
			Code:       c,
		}
	case string:
		return &Err{
			Message:    v,
			StackTrace: string(debug.Stack()),
			Code:       c,
		}
	case []byte:
		return &Err{
			Message:    string(v),
			StackTrace: string(debug.Stack()),
			Code:       c,
		}

	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return &Err{
				Message:    fmt.Sprintf("unsupported err type %T: %+v", v, err),
				StackTrace: string(debug.Stack()),
				Code:       c,
			}
		}
		return &Err{
			Message:    string(jsonData),
			StackTrace: string(debug.Stack()),
			Code:       c,
		}
	}
}

func (e *Err) Error() string {
	return e.Message
}

var _ error = (*Err)(nil)
