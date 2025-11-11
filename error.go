package acp

import (
	"encoding/json"
	"fmt"
)

// Result 是 ACP 调用的通用返回类型。
type Result[T any] struct {
	Value T
	Err   error
}

// Error 表示符合 JSON-RPC 2.0 规范的错误对象。
type Error struct {
	Code    int32           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Error 实现 error 接口。
func (e Error) Error() string {
	if len(e.Message) == 0 {
		return fmt.Sprintf("%d", e.Code)
	}
	if len(e.Data) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, string(e.Data))
}

// WithData 返回包含额外数据的错误拷贝。
func (e Error) WithData(data any) Error {
	switch v := data.(type) {
	case nil:
		e.Data = nil
	case json.RawMessage:
		e.Data = append(json.RawMessage(nil), v...)
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			raw = json.RawMessage(fmt.Sprintf("%q", err.Error()))
		}
		e.Data = raw
	}
	return e
}

// ErrorCode 预定义错误码。
type ErrorCode struct {
	Code    int32
	Message string
}

// ACP 与 JSON-RPC 使用的一些标准错误码。
var (
	ErrorCodeParseError       = ErrorCode{Code: -32700, Message: "Parse error"}
	ErrorCodeInvalidRequest   = ErrorCode{Code: -32600, Message: "Invalid Request"}
	ErrorCodeMethodNotFound   = ErrorCode{Code: -32601, Message: "Method not found"}
	ErrorCodeInvalidParams    = ErrorCode{Code: -32602, Message: "Invalid params"}
	ErrorCodeInternalError    = ErrorCode{Code: -32603, Message: "Internal error"}
	ErrorCodeAuthRequired     = ErrorCode{Code: -32000, Message: "Authentication required"}
	ErrorCodeResourceNotFound = ErrorCode{Code: -32002, Message: "Resource not found"}
)

// NewError 根据错误码创建 Error。
func NewError(code ErrorCode) Error {
	return Error{
		Code:    code.Code,
		Message: code.Message,
	}
}

// ParseError 返回 JSON 解析错误。
func ParseError() Error { return NewError(ErrorCodeParseError) }

// InvalidRequest 返回无效请求错误。
func InvalidRequest() Error { return NewError(ErrorCodeInvalidRequest) }

// MethodNotFound 返回方法不存在错误。
func MethodNotFound() Error { return NewError(ErrorCodeMethodNotFound) }

// InvalidParams 返回参数错误。
func InvalidParams() Error { return NewError(ErrorCodeInvalidParams) }

// InternalError 返回内部错误。
func InternalError() Error { return NewError(ErrorCodeInternalError) }

// AuthRequired 返回需要认证错误。
func AuthRequired() Error { return NewError(ErrorCodeAuthRequired) }

// ResourceNotFound 返回资源不存在错误。
func ResourceNotFound(uri *string) Error {
	err := NewError(ErrorCodeResourceNotFound)
	if uri != nil {
		err = err.WithData(map[string]string{"uri": *uri})
	}
	return err
}

// IntoInternalError 将普通 error 包装为内部错误。
func IntoInternalError(err error) Error {
	if err == nil {
		return InternalError()
	}
	return InternalError().WithData(err.Error())
}
