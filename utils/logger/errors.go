package log

import (
	"net/http"
	"reflect"
	"runtime/debug"
)

var errorType = reflect.TypeOf((*LogError)(nil)).Elem()

// Error is the type which can be used as an implementation of error
type LogError struct {
	StatusCode  int         `json:"-"`
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Reason      string      `json:"reason,omitempty"`
	Category    string      `json:"category,omitempty"`
	SubCategory string      `json:"subCategory,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	trace       string
}

// wrap is the method to be used to wrap the LogError
type wrap interface {
	Wrap(err error) error
}

// unwrap is the method to be used to unwrap the LogError
type unwrap interface {
	Unwrap() error
}

// is used to check if it matches the LogError provided
type is interface {
	Is(err error) bool
}

// as is used to cast error to interface
type as interface {
	As(interface{}) bool
}

// WithStatusCode is used to create a new error with the status code changed
func (e LogError) WithStatusCode(statusCode int) *LogError {
	return &LogError{
		StatusCode:  statusCode,
		Code:        e.Code,
		Message:     e.Message,
		Reason:      e.Reason,
		Category:    e.Category,
		SubCategory: e.SubCategory,
		Details:     e.Details,
		trace:       e.GetTrace(),
	}
}

// WithMessage is used to create a new error with the message changed
func (e LogError) WithMessage(message string) *LogError {
	return &LogError{
		StatusCode:  e.StatusCode,
		Code:        e.Code,
		Message:     message,
		Reason:      e.Reason,
		Category:    e.Category,
		SubCategory: e.SubCategory,
		Details:     e.Details,
		trace:       e.GetTrace(),
	}
}

// WithDetails is used to attach the details to the LogError response
func (e LogError) WithDetails(details interface{}) *LogError {
	return &LogError{
		StatusCode:  e.StatusCode,
		Code:        e.Code,
		Message:     e.Message,
		Reason:      e.Reason,
		Category:    e.Category,
		SubCategory: e.SubCategory,
		Details:     details,
		trace:       e.GetTrace(),
	}
}

// Value is used to get the reference to the value
func (e LogError) Value() *LogError {
	return &LogError{
		StatusCode:  e.StatusCode,
		Code:        e.Code,
		Message:     e.Message,
		Reason:      e.Reason,
		Category:    e.Category,
		SubCategory: e.SubCategory,
		Details:     e.Details,
		trace:       e.GetTrace(),
	}
}

// GetTrace is used to get the current trace
func (e LogError) GetTrace() string {
	if e.trace == "" {
		return string(debug.Stack())
	}
	return e.trace
}

// GetReason is used to get the current reason
func (e *LogError) GetReason() string {
	if e.Reason == "" {
		return e.Message
	}
	return e.Reason
}

// GetSubCategory is used to get the current sub category
func (e *LogError) GetSubCategory() string {
	if e.SubCategory == "" {
		return e.Message
	}
	return e.SubCategory
}

// GetCategory is used to get the current category
func (e *LogError) GetCategory() string {
	if e.Category == "" {
		return e.Message
	}
	return e.Category
}

// Error is used to get the detail from the LogError
func (e *LogError) Error() string {
	if e.Code == "" {
		return e.Message
	}
	return e.Code
}

// Wrap is used to wrap the LogError
func (e *LogError) Wrap(err error) error {
	if err != nil {
		if d, ok := err.(*LogError); ok {
			e.Details = d
		} else {
			e.Details = New(err.Error())
		}
	}
	return e
}

// Unwrap is used to unwrap the LogError to the details
func (e *LogError) Unwrap() error {
	if e == nil {
		return nil
	}

	// first check if details exist
	if e.Details == nil {
		return nil
	}
	// next try to cast the details to error
	if u, ok := e.Details.(error); ok {
		return u
	}
	// otherwise nothing
	return nil
}

// Is used to check if it matches the LogError provided
func (e *LogError) Is(err error) bool {
	return e != nil && err != nil && e.Error() == err.Error()
}

// As is used to case the LogError type to the target
func (e *LogError) As(target interface{}) bool {
	if t, ok := target.(*LogError); ok {
		t.StatusCode = e.StatusCode
		t.Code = e.Code
		t.Message = e.Message
		t.Details = e.Details
		return true
	}
	return false
}

// New is used to create a new error
func New(message string) error {
	return &LogError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		trace:      string(debug.Stack()),
	}
}

// Wrap is used to wrap the LogError into another
func Wrap(parent, err error) error {
	if e, ok := parent.(wrap); ok {
		return e.Wrap(err)
	}
	return parent
}

// Unwrap is used to unwrap the LogError to the details
func Unwrap(err error) error {
	if u, ok := err.(unwrap); ok {
		return u.Unwrap()
	}
	return nil
}

// Is used to check if the LogErrors match or not
func Is(err, target error) bool {
	if target == nil {
		return err == target
	}
	isComparable := reflect.TypeOf(target).Comparable()
	for {
		if isComparable && err == target {
			return true
		}
		if x, ok := err.(is); ok && x.Is(target) {
			return true
		}
		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

// As finds the first error in error's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
func As(err error, target interface{}) bool {
	if target == nil {
		return false
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}
	if e := typ.Elem(); e.Kind() != reflect.Interface && e.Kind() != reflect.Struct && !e.Implements(errorType) {
		return false
	}
	targetType := typ.Elem()
	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			val.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if x, ok := err.(as); ok && x.As(target) {
			return true
		}
		err = Unwrap(err)
	}
	return false
}
