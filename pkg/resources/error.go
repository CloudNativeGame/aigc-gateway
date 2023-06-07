package resources

import "fmt"

type ResourceErrorType string

const (
	ApiCallError   ResourceErrorType = "ApiCallError"
	InternalError  ResourceErrorType = "InternalError"
	ParameterError ResourceErrorType = "ParameterError"
	NotFoundError  ResourceErrorType = "NotFoundError"
)

type ErrorReasonType string

const (
	PauseReason            = "Instance paused"
	NotExistReason         = "Instance not existed"
	ResourceMetaNullReason = "ResourceMeta is Null"
	IdNullReason           = "ID is Null"
	IdNotIntegerReason     = "ID is not Integer"
	NameNullReason         = "Name is Null"
	NamespaceNullReason    = "Namespace is Null"
)

type ResourceError interface {
	Error() string

	Reason() ErrorReasonType

	Type() ResourceErrorType
}

type ResourceErrorImplErrorImpl struct {
	errorType ResourceErrorType
	msg       string
	reason    ErrorReasonType
}

func (r ResourceErrorImplErrorImpl) Error() string {
	return r.msg
}

func (r ResourceErrorImplErrorImpl) Reason() ErrorReasonType {
	return r.reason
}

func (r ResourceErrorImplErrorImpl) Type() ResourceErrorType {
	return r.errorType
}

func IsResourceError(err error) bool {
	_, is := err.(ResourceErrorImplErrorImpl)
	return is
}

func IsNotFoundError(err error) bool {
	re, is := err.(ResourceErrorImplErrorImpl)
	if !is {
		return false
	}
	if re.Type() == NotFoundError {
		return true
	}
	return false
}

func IsApiCallError(err error) bool {
	re, is := err.(ResourceErrorImplErrorImpl)
	if !is {
		return false
	}
	if re.Type() == ApiCallError {
		return true
	}
	return false
}

func IsInternalError(err error) bool {
	re, is := err.(ResourceErrorImplErrorImpl)
	if !is {
		return false
	}
	if re.Type() == InternalError {
		return true
	}
	return false
}

func IsParameterError(err error) bool {
	re, is := err.(ResourceErrorImplErrorImpl)
	if !is {
		return false
	}
	if re.Type() == ParameterError {
		return true
	}
	return false
}

func GetErrorType(err error) ResourceErrorType {
	return err.(ResourceErrorImplErrorImpl).Type()
}

func GetErrorReason(err error) ErrorReasonType {
	return err.(ResourceErrorImplErrorImpl).Reason()
}

func NewResourceError(errorType ResourceErrorType, errorReason ErrorReasonType, msg string, args ...interface{}) ResourceError {
	return ResourceErrorImplErrorImpl{
		errorType: errorType,
		reason:    errorReason,
		msg:       fmt.Sprintf(msg, args...),
	}
}
