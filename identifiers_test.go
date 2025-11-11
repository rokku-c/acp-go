package acp

import (
	"encoding/json"
	"testing"
)

func TestRequestIDMarshalNumber(t *testing.T) {
	id := NewRequestIDNumber(42)
	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(data) != "42" {
		t.Fatalf("expected 42, got %s", string(data))
	}
}

func TestRequestIDMarshalString(t *testing.T) {
	id := NewRequestIDString("abc")
	data, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(data) != `"abc"` {
		t.Fatalf(`expected "abc", got %s`, string(data))
	}
}

func TestRequestIDUnmarshalVariants(t *testing.T) {
	var id RequestID
	if err := json.Unmarshal([]byte("null"), &id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !id.IsNull() {
		t.Fatalf("expected IsNull true")
	}

	if err := json.Unmarshal([]byte("123"), &id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n, ok := id.AsNumber(); !ok || n != 123 {
		t.Fatalf("expected number 123, got %d (ok=%v)", n, ok)
	}

	if err := json.Unmarshal([]byte(`"req-1"`), &id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s, ok := id.AsString(); !ok || s != "req-1" {
		t.Fatalf("expected string req-1, got %q (ok=%v)", s, ok)
	}
}
