package acp

import (
	"context"
	"testing"
	"time"
)

func TestStreamBroadcast(t *testing.T) {
	b := newStreamBroadcast()
	receiver := b.subscribe()

	// Send outgoing request
	b.outgoingRequest("1", "initialize", nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, err := receiver.Recv(ctx)
	if err != nil {
		t.Fatalf("Recv failed: %v", err)
	}
	if msg.Direction != StreamOutgoing {
		t.Fatalf("expected outgoing direction, got %s", msg.Direction)
	}
	if msg.Content.Type != "request" {
		t.Fatalf("expected request type, got %s", msg.Content.Type)
	}
	if msg.Content.ID != "1" {
		t.Fatalf("unexpected id %s", msg.Content.ID)
	}

	// Ensure EOF when channel closes
	b.mu.Lock()
	for ch := range b.subs {
		close(ch)
		delete(b.subs, ch)
	}
	b.mu.Unlock()

	if _, err := receiver.Recv(ctx); err == nil {
		t.Fatalf("expected error after channel closed")
	}
}
