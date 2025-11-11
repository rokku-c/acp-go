package acp

import (
	"context"
	"encoding/json"
	"io"
	"sync"
)

// StreamMessageDirection 用于标识消息方向。
type StreamMessageDirection string

const (
	StreamIncoming StreamMessageDirection = "incoming"
	StreamOutgoing StreamMessageDirection = "outgoing"
)

// StreamMessageContent 描述消息内容。
type StreamMessageContent struct {
	Type   string          `json:"type"`
	ID     string          `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

// StreamMessage 表示一条流消息。
type StreamMessage struct {
	Direction StreamMessageDirection `json:"direction"`
	Content   StreamMessageContent   `json:"content"`
}

// StreamReceiver 订阅者。
type StreamReceiver struct {
	ch <-chan StreamMessage
}

// Recv 从流中读取一条消息。
func (r StreamReceiver) Recv(ctx context.Context) (StreamMessage, error) {
	select {
	case <-ctx.Done():
		return StreamMessage{}, ctx.Err()
	case msg, ok := <-r.ch:
		if !ok {
			return StreamMessage{}, io.EOF
		}
		return msg, nil
	}
}

type streamBroadcast struct {
	mu   sync.Mutex
	subs map[chan StreamMessage]struct{}
}

func newStreamBroadcast() *streamBroadcast {
	return &streamBroadcast{
		subs: make(map[chan StreamMessage]struct{}),
	}
}

func (b *streamBroadcast) subscribe() StreamReceiver {
	ch := make(chan StreamMessage, 32)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()
	return StreamReceiver{ch: ch}
}

func (b *streamBroadcast) send(msg StreamMessage) {
	b.mu.Lock()
	for ch := range b.subs {
		select {
		case ch <- msg:
		default:
		}
	}
	b.mu.Unlock()
}

func (b *streamBroadcast) outgoingRequest(id, method string, params json.RawMessage) {
	b.send(StreamMessage{
		Direction: StreamOutgoing,
		Content: StreamMessageContent{
			Type:   "request",
			ID:     id,
			Method: method,
			Params: params,
		},
	})
}

func (b *streamBroadcast) outgoingResponse(id string, result *json.RawMessage, err *Error) {
	content := StreamMessageContent{
		Type: "response",
		ID:   id,
	}
	if result != nil {
		content.Result = *result
	}
	if err != nil {
		content.Error = err
	}
	b.send(StreamMessage{
		Direction: StreamOutgoing,
		Content:   content,
	})
}

func (b *streamBroadcast) outgoingNotification(method string, params json.RawMessage) {
	b.send(StreamMessage{
		Direction: StreamOutgoing,
		Content: StreamMessageContent{
			Type:   "notification",
			Method: method,
			Params: params,
		},
	})
}

func (b *streamBroadcast) incomingRequest(id, method string, params json.RawMessage) {
	b.send(StreamMessage{
		Direction: StreamIncoming,
		Content: StreamMessageContent{
			Type:   "request",
			ID:     id,
			Method: method,
			Params: params,
		},
	})
}

func (b *streamBroadcast) incomingResponse(id string, result *json.RawMessage, err *Error) {
	content := StreamMessageContent{
		Type: "response",
		ID:   id,
	}
	if result != nil {
		content.Result = *result
	}
	if err != nil {
		content.Error = err
	}
	b.send(StreamMessage{
		Direction: StreamIncoming,
		Content:   content,
	})
}

func (b *streamBroadcast) incomingNotification(method string, params json.RawMessage) {
	b.send(StreamMessage{
		Direction: StreamIncoming,
		Content: StreamMessageContent{
			Type:   "notification",
			Method: method,
			Params: params,
		},
	})
}
