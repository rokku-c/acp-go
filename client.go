package acp

import (
	"context"
	"encoding/json"
)

// Client 定义 ACP 客户端需要实现的接口。
type Client interface {
	RequestPermission(context.Context, RequestPermissionRequest) (RequestPermissionResponse, error)
	SessionNotification(context.Context, SessionNotification) error
	WriteTextFile(context.Context, WriteTextFileRequest) (WriteTextFileResponse, error)
	ReadTextFile(context.Context, ReadTextFileRequest) (ReadTextFileResponse, error)
	CreateTerminal(context.Context, CreateTerminalRequest) (CreateTerminalResponse, error)
	TerminalOutput(context.Context, TerminalOutputRequest) (TerminalOutputResponse, error)
	ReleaseTerminal(context.Context, ReleaseTerminalRequest) (ReleaseTerminalResponse, error)
	WaitForTerminalExit(context.Context, WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error)
	KillTerminalCommand(context.Context, KillTerminalCommandRequest) (KillTerminalCommandResponse, error)
	ExtMethod(context.Context, ExtRequest) (ExtResponse, error)
	ExtNotification(context.Context, ExtNotification) error
}

// UnimplementedClient 为可选方法提供默认实现。
type UnimplementedClient struct{}

func (UnimplementedClient) RequestPermission(context.Context, RequestPermissionRequest) (RequestPermissionResponse, error) {
	return RequestPermissionResponse{}, MethodNotFound()
}

func (UnimplementedClient) SessionNotification(context.Context, SessionNotification) error {
	return nil
}

func (UnimplementedClient) WriteTextFile(context.Context, WriteTextFileRequest) (WriteTextFileResponse, error) {
	return WriteTextFileResponse{}, MethodNotFound()
}

func (UnimplementedClient) ReadTextFile(context.Context, ReadTextFileRequest) (ReadTextFileResponse, error) {
	return ReadTextFileResponse{}, MethodNotFound()
}

func (UnimplementedClient) CreateTerminal(context.Context, CreateTerminalRequest) (CreateTerminalResponse, error) {
	return CreateTerminalResponse{}, MethodNotFound()
}

func (UnimplementedClient) TerminalOutput(context.Context, TerminalOutputRequest) (TerminalOutputResponse, error) {
	return TerminalOutputResponse{}, MethodNotFound()
}

func (UnimplementedClient) ReleaseTerminal(context.Context, ReleaseTerminalRequest) (ReleaseTerminalResponse, error) {
	return ReleaseTerminalResponse{}, MethodNotFound()
}

func (UnimplementedClient) WaitForTerminalExit(context.Context, WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error) {
	return WaitForTerminalExitResponse{}, MethodNotFound()
}

func (UnimplementedClient) KillTerminalCommand(context.Context, KillTerminalCommandRequest) (KillTerminalCommandResponse, error) {
	return KillTerminalCommandResponse{}, MethodNotFound()
}

func (UnimplementedClient) ExtMethod(context.Context, ExtRequest) (ExtResponse, error) {
	return ExtResponse(json.RawMessage("null")), nil
}

func (UnimplementedClient) ExtNotification(context.Context, ExtNotification) error {
	return nil
}
