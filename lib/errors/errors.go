package errors

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	err     []error
	Type    string
	Message string
	Data    interface{}
}

func (e *Error) Error() string {
	msg := make([]string, 0)
	for _, value := range e.err {
		msg = append(msg, value.Error())
	}

	bytes, _ := json.Marshal(msg)
	return string(bytes)
}

func NewError(errMessage ...string) *Error {
	err := &Error{
		Type: TypeInternalServerError,
		err:  make([]error, 0),
	}
	if errMessage == nil {
		return err
	}

	for _, value := range errMessage {
		err.err = append(err.err, fmt.Errorf("%v", value))
	}

	return err
}

func NewWrapError(err error, errMessage string) *Error {
	errs, ok := err.(*Error)
	if ok {
		errs.err = append(errs.err, fmt.Errorf("%v", errMessage))
		return errs
	}

	errs = &Error{
		Type: TypeInternalServerError,
		err:  make([]error, 0),
	}
	errs.err = append(errs.err, fmt.Errorf("%v", err.Error()))
	errs.err = append(errs.err, fmt.Errorf("%v", errMessage))
	return errs
}
func (e *Error) SetType(errType string) *Error {
	e.Type = errType
	return e
}
