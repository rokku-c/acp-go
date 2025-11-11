package acp

import (
	"context"
	"encoding/json"
	"io"
)

// AgentSideConnection 为代理提供 Client 接口。
type AgentSideConnection struct {
	rpc *rpcConnection
}

// NewAgentSideConnection 创建代理侧连接。
func NewAgentSideConnection(
	ctx context.Context,
	agent Agent,
	outgoing io.Writer,
	incoming io.Reader,
) *AgentSideConnection {
	handler := &agentInboundHandler{agent: agent}
	return &AgentSideConnection{
		rpc: newRPCConnection(ctx, handler, outgoing, incoming),
	}
}

// Close 关闭连接。
func (a *AgentSideConnection) Close() {
	a.rpc.Close(io.EOF)
}

// Subscribe 订阅流。
func (a *AgentSideConnection) Subscribe() StreamReceiver {
	return a.rpc.subscribe()
}

// RequestPermission 调用 session/request_permission。
func (a *AgentSideConnection) RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
	var resp RequestPermissionResponse
	raw, err := a.rpc.request(ctx, ClientMethods.SessionRequestPermission, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// WriteTextFile 调用 fs/write_text_file。
func (a *AgentSideConnection) WriteTextFile(ctx context.Context, req WriteTextFileRequest) (WriteTextFileResponse, error) {
	var resp WriteTextFileResponse
	raw, err := a.rpc.request(ctx, ClientMethods.FSWriteTextFile, req)
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

// ReadTextFile 调用 fs/read_text_file。
func (a *AgentSideConnection) ReadTextFile(ctx context.Context, req ReadTextFileRequest) (ReadTextFileResponse, error) {
	var resp ReadTextFileResponse
	raw, err := a.rpc.request(ctx, ClientMethods.FSReadTextFile, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// CreateTerminal 调用 terminal/create。
func (a *AgentSideConnection) CreateTerminal(ctx context.Context, req CreateTerminalRequest) (CreateTerminalResponse, error) {
	var resp CreateTerminalResponse
	raw, err := a.rpc.request(ctx, ClientMethods.TerminalCreate, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// TerminalOutput 调用 terminal/output。
func (a *AgentSideConnection) TerminalOutput(ctx context.Context, req TerminalOutputRequest) (TerminalOutputResponse, error) {
	var resp TerminalOutputResponse
	raw, err := a.rpc.request(ctx, ClientMethods.TerminalOutput, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// ReleaseTerminal 调用 terminal/release。
func (a *AgentSideConnection) ReleaseTerminal(ctx context.Context, req ReleaseTerminalRequest) (ReleaseTerminalResponse, error) {
	var resp ReleaseTerminalResponse
	raw, err := a.rpc.request(ctx, ClientMethods.TerminalRelease, req)
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

// WaitForTerminalExit 调用 terminal/wait_for_exit。
func (a *AgentSideConnection) WaitForTerminalExit(ctx context.Context, req WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error) {
	var resp WaitForTerminalExitResponse
	raw, err := a.rpc.request(ctx, ClientMethods.TerminalWaitForExit, req)
	if err != nil {
		return resp, err
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// KillTerminalCommand 调用 terminal/kill。
func (a *AgentSideConnection) KillTerminalCommand(ctx context.Context, req KillTerminalCommandRequest) (KillTerminalCommandResponse, error) {
	var resp KillTerminalCommandResponse
	raw, err := a.rpc.request(ctx, ClientMethods.TerminalKill, req)
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

// SessionNotification 发送 session/update。
func (a *AgentSideConnection) SessionNotification(ctx context.Context, note SessionNotification) error {
	return a.rpc.notify(ctx, ClientMethods.SessionUpdate, note)
}

// ExtMethod 调用扩展方法。
func (a *AgentSideConnection) ExtMethod(ctx context.Context, method string, params json.RawMessage) (ExtResponse, error) {
	var resp ExtResponse
	raw, err := a.rpc.request(ctx, "_"+method, ExtRequest{Method: method, Params: params})
	if err != nil {
		return resp, err
	}
	resp = ExtResponse(raw)
	return resp, nil
}

// ExtNotification 发送扩展通知。
func (a *AgentSideConnection) ExtNotification(ctx context.Context, method string, params json.RawMessage) error {
	return a.rpc.notify(ctx, "_"+method, ExtNotification{Method: method, Params: params})
}
