package acp

import (
	"context"
	"encoding/json"
	"io"
)

// ClientSideConnection 为客户端提供 Agent 接口。
type ClientSideConnection struct {
	rpc *rpcConnection
}

// NewClientSideConnection 创建客户端侧连接。
func NewClientSideConnection(
	ctx context.Context,
	client Client,
	outgoing io.Writer,
	incoming io.Reader,
) *ClientSideConnection {
	handler := &clientInboundHandler{client: client}
	return &ClientSideConnection{
		rpc: newRPCConnection(ctx, handler, outgoing, incoming),
	}
}

// Close 关闭连接。
func (c *ClientSideConnection) Close() {
	c.rpc.Close(io.EOF)
}

// Subscribe 订阅消息流。
func (c *ClientSideConnection) Subscribe() StreamReceiver {
	return c.rpc.subscribe()
}

// Initialize 调用 initialize 方法。
func (c *ClientSideConnection) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	var resp InitializeResponse
	raw, err := c.rpc.request(ctx, AgentMethods.Initialize, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// Authenticate 调用 authenticate。
func (c *ClientSideConnection) Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error) {
	var resp AuthenticateResponse
	raw, err := c.rpc.request(ctx, AgentMethods.Authenticate, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// NewSession 调用 session/new。
func (c *ClientSideConnection) NewSession(ctx context.Context, req NewSessionRequest) (NewSessionResponse, error) {
	var resp NewSessionResponse
	raw, err := c.rpc.request(ctx, AgentMethods.SessionNew, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// LoadSession 调用 session/load。
func (c *ClientSideConnection) LoadSession(ctx context.Context, req LoadSessionRequest) (LoadSessionResponse, error) {
	var resp LoadSessionResponse
	raw, err := c.rpc.request(ctx, AgentMethods.SessionLoad, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// SetSessionMode 调用 session/set_mode。
func (c *ClientSideConnection) SetSessionMode(ctx context.Context, req SetSessionModeRequest) (SetSessionModeResponse, error) {
	var resp SetSessionModeResponse
	raw, err := c.rpc.request(ctx, AgentMethods.SessionSetMode, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// SetSessionModel 调用 session/set_model。
func (c *ClientSideConnection) SetSessionModel(ctx context.Context, req SetSessionModelRequest) (SetSessionModelResponse, error) {
	var resp SetSessionModelResponse
	raw, err := c.rpc.request(ctx, AgentMethods.SessionSetModel, req)
	if err != nil {
		return resp, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return resp, nil
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// Prompt 调用 session/prompt。
func (c *ClientSideConnection) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	var resp PromptResponse
	raw, err := c.rpc.request(ctx, AgentMethods.SessionPrompt, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// Cancel 发送取消通知。
func (c *ClientSideConnection) Cancel(ctx context.Context, note CancelNotification) error {
	return c.rpc.notify(ctx, AgentMethods.SessionCancel, note)
}

// ExtMethod 调用扩展方法。
func (c *ClientSideConnection) ExtMethod(ctx context.Context, method string, params json.RawMessage) (ExtResponse, error) {
	var resp ExtResponse
	raw, err := c.rpc.request(ctx, "_"+method, ExtRequest{Method: method, Params: params})
	if err != nil {
		return resp, err
	}
	resp = ExtResponse(raw)
	return resp, nil
}

// ExtNotification 发送扩展通知。
func (c *ClientSideConnection) ExtNotification(ctx context.Context, method string, params json.RawMessage) error {
	return c.rpc.notify(ctx, "_"+method, ExtNotification{Method: method, Params: params})
}
