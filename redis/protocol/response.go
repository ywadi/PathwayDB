package protocol

// ResponseType represents the type of Redis response
type ResponseType int

const (
	ResponseTypeString ResponseType = iota
	ResponseTypeInt
	ResponseTypeArray
	ResponseTypeNestedArray
	ResponseTypeBulk
	ResponseTypeNull
	ResponseTypeError
)

// Response represents a Redis protocol response
type Response struct {
	Type             ResponseType
	StringValue      string
	IntValue         int64
	ArrayValue       []string
	NestedArrayValue []interface{}
}

// NewStringResponse creates a simple string response
func NewStringResponse(value string) *Response {
	return &Response{
		Type:        ResponseTypeString,
		StringValue: value,
	}
}

// NewIntResponse creates an integer response
func NewIntResponse(value int64) *Response {
	return &Response{
		Type:     ResponseTypeInt,
		IntValue: value,
	}
}

// NewNestedArrayResponse creates a nested array response
func NewNestedArrayResponse(values []interface{}) *Response {
	return &Response{
		Type:             ResponseTypeNestedArray,
		NestedArrayValue: values,
	}
}

// NewArrayResponse creates an array response
func NewArrayResponse(values []string) *Response {
	return &Response{
		Type:       ResponseTypeArray,
		ArrayValue: values,
	}
}

// NewBulkResponse creates a bulk string response
func NewBulkResponse(value string) *Response {
	return &Response{
		Type:        ResponseTypeBulk,
		StringValue: value,
	}
}

// NewNullResponse creates a null response
func NewNullResponse() *Response {
	return &Response{
		Type: ResponseTypeNull,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(message string) *Response {
	return &Response{
		Type:        ResponseTypeError,
		StringValue: message,
	}
}

// OK returns a standard OK response
func OK() *Response {
	return NewStringResponse("OK")
}
