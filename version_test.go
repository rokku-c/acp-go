package acp

import (
	"encoding/json"
	"testing"
)

func TestProtocolVersionUnmarshalNumber(t *testing.T) {
	var v ProtocolVersion
	if err := json.Unmarshal([]byte("1"), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Value() != 1 {
		t.Fatalf("expected value 1, got %d", v.Value())
	}
}

func TestProtocolVersionUnmarshalStringFallback(t *testing.T) {
	var v ProtocolVersion
	if err := json.Unmarshal([]byte(`"1.2.3"`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Value() != ProtocolVersionV0.Value() {
		t.Fatalf("expected fallback to 0, got %d", v.Value())
	}
}

func TestProtocolVersionMarshalRoundtrip(t *testing.T) {
	v := ProtocolVersionV1
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(out) != "1" {
		t.Fatalf("expected `1`, got %s", string(out))
	}
}

func TestProtocolVersionUnmarshalTooLarge(t *testing.T) {
	var v ProtocolVersion
	if err := json.Unmarshal([]byte("70000"), &v); err == nil {
		t.Fatalf("expected error for value > uint16, got nil")
	}
}
