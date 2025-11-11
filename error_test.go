package acp

import (
	"errors"
	"testing"
)

func TestErrorWithData(t *testing.T) {
	err := MethodNotFound().WithData(map[string]string{"reason": "missing"})
	if err.Code != ErrorCodeMethodNotFound.Code {
		t.Fatalf("unexpected code: %d", err.Code)
	}
	if err.Message != ErrorCodeMethodNotFound.Message {
		t.Fatalf("unexpected message: %s", err.Message)
	}
	if len(err.Data) == 0 {
		t.Fatalf("expected data to be present")
	}
}

func TestIntoInternalError(t *testing.T) {
	src := errors.New("boom")
	err := IntoInternalError(src)
	if err.Code != ErrorCodeInternalError.Code {
		t.Fatalf("unexpected code: %d", err.Code)
	}
	if err.Data == nil {
		t.Fatalf("expected data payload")
	}
}

func TestResourceNotFound(t *testing.T) {
	uri := "file:///tmp/test"
	err := ResourceNotFound(&uri)
	if err.Code != ErrorCodeResourceNotFound.Code {
		t.Fatalf("unexpected code: %d", err.Code)
	}
	if err.Data == nil {
		t.Fatalf("expected data")
	}
}
