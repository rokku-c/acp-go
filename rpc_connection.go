package acp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type rpcConnection struct {
	outgoing  chan jsonrpcEnvelope
	encoder   *json.Encoder
	reader    *bufio.Reader
	handler   inboundHandler
	nextID    atomic.Int64
	pendingMu sync.Mutex
	pending   map[string]*pendingRequest
	broadcast *streamBroadcast
	closeOnce sync.Once
	closeErr  error
	closeCh   chan struct{}
}

type inboundHandler interface {
	handleRequest(context.Context, string, json.RawMessage) (any, Error, bool)
	handleNotification(context.Context, string, json.RawMessage) Error
}

type pendingRequest struct {
	result chan rpcResult
}

type rpcResult struct {
	result json.RawMessage
	err    *Error
}

type jsonrpcEnvelope struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  *json.RawMessage `json:"params,omitempty"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

func newRPCConnection(
	ctx context.Context,
	handler inboundHandler,
	outgoingWriter io.Writer,
	incomingReader io.Reader,
) *rpcConnection {
	conn := &rpcConnection{
		outgoing:  make(chan jsonrpcEnvelope, 32),
		encoder:   json.NewEncoder(outgoingWriter),
		reader:    bufio.NewReader(incomingReader),
		handler:   handler,
		pending:   make(map[string]*pendingRequest),
		broadcast: newStreamBroadcast(),
		closeCh:   make(chan struct{}),
	}

	go conn.writeLoop()
	go conn.readLoop(ctx)

	return conn
}

func (c *rpcConnection) Close(err error) {
	c.closeOnce.Do(func() {
		c.closeErr = err
		close(c.closeCh)
		close(c.outgoing)
		c.pendingMu.Lock()
		for _, p := range c.pending {
			p.result <- rpcResult{err: &Error{Code: ErrorCodeInternalError.Code, Message: err.Error()}}
		}
		c.pending = map[string]*pendingRequest{}
		c.pendingMu.Unlock()
	})
}

func (c *rpcConnection) notify(ctx context.Context, method string, params any) error {
	raw, err := marshalRaw(params)
	if err != nil {
		return err
	}

	envelope := jsonrpcEnvelope{
		JSONRPC: "2.0",
		Method:  method,
	}
	if raw != nil {
		envelope.Params = &raw
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.outgoing <- envelope:
		c.broadcast.outgoingNotification(method, raw)
		return nil
	case <-c.closeCh:
		if c.closeErr != nil {
			return c.closeErr
		}
		return fmt.Errorf("connection closed")
	}
}

func (c *rpcConnection) request(ctx context.Context, method string, params any) (json.RawMessage, error) {
	raw, err := marshalRaw(params)
	if err != nil {
		return nil, err
	}

	id := c.nextID.Add(1)
	idRaw := json.RawMessage([]byte(fmt.Sprintf("%d", id)))
	envelope := jsonrpcEnvelope{
		JSONRPC: "2.0",
		ID:      &idRaw,
		Method:  method,
	}
	if raw != nil {
		envelope.Params = &raw
	}

	pending := &pendingRequest{
		result: make(chan rpcResult, 1),
	}
	key := string(idRaw)

	c.pendingMu.Lock()
	c.pending[key] = pending
	c.pendingMu.Unlock()

	select {
	case <-ctx.Done():
		c.removePending(key)
		return nil, ctx.Err()
	case c.outgoing <- envelope:
		c.broadcast.outgoingRequest(string(idRaw), method, raw)
	case <-c.closeCh:
		c.removePending(key)
		if c.closeErr != nil {
			return nil, c.closeErr
		}
		return nil, fmt.Errorf("connection closed")
	}

	select {
	case <-ctx.Done():
		c.removePending(key)
		return nil, ctx.Err()
	case <-c.closeCh:
		c.removePending(key)
		if c.closeErr != nil {
			return nil, c.closeErr
		}
		return nil, fmt.Errorf("connection closed")
	case res := <-pending.result:
		if res.err != nil {
			return nil, res.err
		}
		return res.result, nil
	}
}

func (c *rpcConnection) removePending(key string) {
	c.pendingMu.Lock()
	delete(c.pending, key)
	c.pendingMu.Unlock()
}

func (c *rpcConnection) writeLoop() {
	for envelope := range c.outgoing {
		if err := c.encoder.Encode(envelope); err != nil {
			c.Close(err)
			return
		}
	}
}

func (c *rpcConnection) readLoop(ctx context.Context) {
	scanner := bufio.NewScanner(c.reader)
	for scanner.Scan() {
		line := scanner.Bytes()
		var envelope jsonrpcEnvelope
		if err := json.Unmarshal(line, &envelope); err != nil {
			continue
		}
		c.handleIncoming(ctx, envelope)
	}
	if err := scanner.Err(); err != nil {
		c.Close(err)
	} else {
		c.Close(io.EOF)
	}
}

func (c *rpcConnection) handleIncoming(ctx context.Context, envelope jsonrpcEnvelope) {
	hasID := envelope.ID != nil
	hasMethod := envelope.Method != ""
	switch {
	case hasID && hasMethod:
		params := json.RawMessage(nil)
		if envelope.Params != nil {
			params = *envelope.Params
		}
		c.broadcast.incomingRequest(string(*envelope.ID), envelope.Method, params)
		go c.dispatchRequest(ctx, envelope.Method, params, envelope.ID)
	case hasID:
		key := string(*envelope.ID)
		c.pendingMu.Lock()
		pending := c.pending[key]
		if pending != nil {
			delete(c.pending, key)
		}
		c.pendingMu.Unlock()
		if pending == nil {
			return
		}

		if envelope.Error != nil {
			errCopy := *envelope.Error
			c.broadcast.incomingResponse(key, nil, envelope.Error)
			pending.result <- rpcResult{err: &errCopy}
			return
		}
		if envelope.Result != nil {
			c.broadcast.incomingResponse(key, envelope.Result, nil)
			pending.result <- rpcResult{result: *envelope.Result}
			return
		}
		c.broadcast.incomingResponse(key, nil, nil)
		pending.result <- rpcResult{}
	case hasMethod:
		params := json.RawMessage(nil)
		if envelope.Params != nil {
			params = *envelope.Params
		}
		c.broadcast.incomingNotification(envelope.Method, params)
		if err := c.handler.handleNotification(ctx, envelope.Method, params); err.Code != 0 && err.Message != "" {
			// notifications do not send response; log could be added
		}
	default:
		// ignore invalid message
	}
}

func (c *rpcConnection) dispatchRequest(ctx context.Context, method string, params json.RawMessage, id *json.RawMessage) {
	res, err, ok := c.handler.handleRequest(ctx, method, params)
	if !ok {
		return
	}
	envelope := jsonrpcEnvelope{
		JSONRPC: "2.0",
	}
	if id != nil {
		idCopy := append(json.RawMessage(nil), (*id)...)
		envelope.ID = &idCopy
	}
	if err.Code != 0 || err.Message != "" {
		errCopy := err
		envelope.Error = &errCopy
		c.broadcast.outgoingResponse(string(*envelope.ID), nil, &errCopy)
	} else {
		raw, marshalErr := marshalRaw(res)
		if marshalErr != nil {
			errCopy := IntoInternalError(marshalErr)
			envelope.Error = &errCopy
			c.broadcast.outgoingResponse(string(*envelope.ID), nil, &errCopy)
		} else {
			if raw != nil {
				envelope.Result = &raw
			} else {
				nullRaw := json.RawMessage("null")
				envelope.Result = &nullRaw
			}
			c.broadcast.outgoingResponse(string(*envelope.ID), envelope.Result, nil)
		}
	}
	select {
	case c.outgoing <- envelope:
	case <-ctx.Done():
	case <-c.closeCh:
	}
}

func marshalRaw(value any) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case json.RawMessage:
		return append(json.RawMessage(nil), v...), nil
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
}

func (c *rpcConnection) subscribe() StreamReceiver {
	return c.broadcast.subscribe()
}
