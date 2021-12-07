package errors

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"time"
)

type wlError struct {
	baseError
	error
}

type wlErrorC struct {
	wlError
	errCode ErrorCode
}

type baseError struct {
	time     time.Time
	location string
	msg      string
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func NewErrorf(err error, format string, a ...interface{}) *wlError {
	e := &wlError{}
	e.newf(err, format, a...)

	return e
}

func NewError(err error) *wlError {
	return NewErrorf(err, "")
}

func NewErrorCf(errCode ErrorCode, caused error, format string, a ...interface{}) *wlErrorC {
	e := &wlErrorC{}
	e.newf(errCode, caused, format, a...)

	return e
}

func NewErrorC(errCode ErrorCode, caused error) *wlErrorC {
	return NewErrorCf(errCode, caused, "")
}

func (e wlError) Error() string {
	var errStr string

	errStr = fmt.Sprintf("s - %v | %v", e.time.Format(time.RFC3339Nano), e.location)

	if e.error != nil {
		errStr = errStr + " | caused by:" + e.error.Error()
	}

	if e.baseError.msg != "" {
		errStr = errStr + " | detail:" + e.baseError.msg
	}

	return errStr
}

func (e wlError) Unwrap() error {
	return e.error
}

func (e wlError) Is(target error) bool {
	if e == target {
		return true
	}

	return errors.Is(e.error, target)
}

func (e wlError) As(target interface{}) bool {
	if sameTypeAssign(&e, target) {
		return true
	}

	return errors.As(e.error, target)
}

func (e wlErrorC) Error() string {
	var errStr string

	errStr = fmt.Sprintf("s - %v | %v | errcode:%v(%v)",
		e.time.Format(time.RFC3339Nano), e.location, e.errCode.Code, e.errCode.Message)

	if e.error != nil {
		errStr = errStr + " | caused by:" + e.error.Error()
	}

	if e.baseError.msg != "" {
		errStr = errStr + " | detail:" + e.baseError.msg
	}

	return errStr
}

func (e wlErrorC) Is(target error) bool {
	if e == target {
		return true
	}

	return errors.Is(e, target)
}

func (e wlErrorC) As(target interface{}) bool {
	if sameTypeAssign(&e, target) {
		return true
	}

	return errors.As(e.error, target)
}

func (e wlErrorC) Unwrap() error {
	return e.error
}

func (e *wlError) newf(caused error, format string, a ...interface{}) {
	e.baseError.newf(format, a...)
	e.error = caused
}

func (e *wlErrorC) newf(errCode ErrorCode, caused error, format string, a ...interface{}) {
	e.wlError.newf(caused, format, a...)
	e.errCode = errCode
}

func (be *baseError) newf(format string, a ...interface{}) {
	be.msg = fmt.Sprintf(format, a...)
	be.time = time.Now()
	pc, filePath, lineNo, ok := runtime.Caller(4)
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(filePath)
	be.location = fmt.Sprintf("%v:%v:%v", funcName, fileName, lineNo)
}

func Code(err error) ErrorCode {
	var e *wlErrorC

	if ok := errors.As(err, &e); ok {
		return e.errCode
	}

	errcode := ErrCdInterSys
	errcode.Message = err.Error()

	return errcode
}

func sameTypeAssign(source interface{}, target interface{}) bool {
	if reflect.TypeOf(target).String() == reflect.TypeOf(source).String() {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(source))

		return true
	}

	return false
}
