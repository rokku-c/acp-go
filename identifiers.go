package acp

import "encoding/json"

// SessionID 会话唯一标识。
type SessionID string

// MarshalJSON ensures Arc semantics? In Rust they just string. We'll just string.

// AuthMethodID 认证方式标识。
type AuthMethodID string

// PermissionOptionID 权限选项标识。
type PermissionOptionID string

// TerminalID 终端标识。
type TerminalID string

// RequestID 支持字符串、数字与 null。
type RequestID struct {
	raw any
}

// NewRequestIDNumber 创建数字 ID。
func NewRequestIDNumber(v int64) RequestID {
	return RequestID{raw: v}
}

// NewRequestIDString 创建字符串 ID。
func NewRequestIDString(v string) RequestID {
	return RequestID{raw: v}
}

// RequestIDNull 表示 null。
var RequestIDNull = RequestID{raw: nil}

// MarshalJSON 实现 json.Marshaler。
func (r RequestID) MarshalJSON() ([]byte, error) {
	switch v := r.raw.(type) {
	case nil:
		return []byte("null"), nil
	case int64:
		return json.Marshal(v)
	case string:
		return json.Marshal(v)
	default:
		return json.Marshal(v)
	}
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (r *RequestID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		r.raw = nil
		return nil
	}
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		r.raw = num
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		r.raw = str
		return nil
	}
	return json.Unmarshal(data, &r.raw)
}

// IsNull 判断是否为 null。
func (r RequestID) IsNull() bool {
	return r.raw == nil
}

// AsNumber 返回数字形式以及是否成功。
func (r RequestID) AsNumber() (int64, bool) {
	v, ok := r.raw.(int64)
	return v, ok
}

// AsString 返回字符串形式以及是否成功。
func (r RequestID) AsString() (string, bool) {
	v, ok := r.raw.(string)
	return v, ok
}
