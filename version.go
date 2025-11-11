package acp

import (
	"encoding/json"
	"fmt"
)

// ProtocolVersion 表示 ACP 协议版本号。
type ProtocolVersion struct {
	value uint16
}

// 常用版本常量。
var (
	ProtocolVersionV0      = ProtocolVersion{value: 0}
	ProtocolVersionV1      = ProtocolVersion{value: 1}
	ProtocolVersionCurrent = ProtocolVersionV1
)

// NewProtocolVersion 构造指定版本。
func NewProtocolVersion(v uint16) ProtocolVersion {
	return ProtocolVersion{value: v}
}

// Value 返回版本号。
func (v ProtocolVersion) Value() uint16 {
	return v.value
}

// MarshalJSON 实现 json.Marshaler。
func (v ProtocolVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// UnmarshalJSON 实现 json.Unmarshaler。
func (v *ProtocolVersion) UnmarshalJSON(data []byte) error {
	switch {
	case len(data) == 0:
		return fmt.Errorf("protocol version: empty input")
	case data[0] == '"':
		// 旧协议使用字符串形式，视为 0
		*v = ProtocolVersionV0
		return nil
	default:
		var val uint16
		if err := json.Unmarshal(data, &val); err != nil {
			return fmt.Errorf("protocol version: %w", err)
		}
		*v = ProtocolVersion{value: val}
		return nil
	}
}
