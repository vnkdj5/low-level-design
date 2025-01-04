package errors

import "fmt"

type BaseError struct {
	Code           string      `json:"code,omitempty"`
	Description    string      `json:"description,omitempty"`
	AdditionalInfo interface{} `json:"additionalInfo,omitempty"`
	StatusCode     int         `json:"status_code,omitempty"`
	MethodCode     string      `json:"-"`
	Err            error       `json:"-"`
}

func (e1 *BaseError) Error() string {
	return e1.Description

}

func (e1 *BaseError) WithMethodCode(methodCode string) *BaseError {
	e1.MethodCode = methodCode
	return e1
}

func (e1 *BaseError) WithErrorCode(errorCode string) *BaseError {
	e1.Code = errorCode
	return e1
}

func (e1 *BaseError) ErrorCode() string {
	errorCode := e1.Code
	if len(errorCode) != 4 {
		errorCode = UnknownErrorCode
	}
	methodCode := e1.MethodCode
	if len(methodCode) != 4 {
		methodCode = UnknownMethodErrorCode
	}
	return fmt.Sprintf("%s%s%s", ServiceCode, methodCode, errorCode)

}

func (e *BaseError) UnWrapError() string {

	// Check if the error is of type BaseError
	if nestedErr, ok := e.Err.(*BaseError); ok {
		// Recursively unwrap and concatenate descriptions
		return e.Description + " : " + nestedErr.UnWrapError()
	}
	// If e.Err is not a BaseError, concatenate its error message
	if e.Err != nil {
		return e.Description + " : " + e.Err.Error()
	}

	return e.Description

}

// Get Full 12 digit Error code based on methodCode and errorCode
func GetErrorCode(methodCode string, errorCode string) string {
	if len(errorCode) != 4 {
		errorCode = UnknownErrorCode
	}
	if len(methodCode) != 4 {
		methodCode = UnknownMethodErrorCode
	}
	return fmt.Sprintf("%s%s%s", ServiceCode, methodCode, errorCode)
}
