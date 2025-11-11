package acp

import (
	"context"
	"encoding/json"
)

// Agent 定义 ACP 代理需要实现的接口。
type Agent interface {
	Initialize(context.Context, InitializeRequest) (InitializeResponse, error)
	Authenticate(context.Context, AuthenticateRequest) (AuthenticateResponse, error)
	NewSession(context.Context, NewSessionRequest) (NewSessionResponse, error)
	LoadSession(context.Context, LoadSessionRequest) (LoadSessionResponse, error)
	SetSessionMode(context.Context, SetSessionModeRequest) (SetSessionModeResponse, error)
	SetSessionModel(context.Context, SetSessionModelRequest) (SetSessionModelResponse, error)
	Prompt(context.Context, PromptRequest) (PromptResponse, error)
	Cancel(context.Context, CancelNotification) error
	ExtMethod(context.Context, ExtRequest) (ExtResponse, error)
	ExtNotification(context.Context, ExtNotification) error
}

// UnimplementedAgent 为可选方法提供默认实现。
type UnimplementedAgent struct{}

func (UnimplementedAgent) Initialize(context.Context, InitializeRequest) (InitializeResponse, error) {
	return InitializeResponse{}, MethodNotFound()
}

func (UnimplementedAgent) Authenticate(context.Context, AuthenticateRequest) (AuthenticateResponse, error) {
	return AuthenticateResponse{}, MethodNotFound()
}

func (UnimplementedAgent) NewSession(context.Context, NewSessionRequest) (NewSessionResponse, error) {
	return NewSessionResponse{}, MethodNotFound()
}

func (UnimplementedAgent) LoadSession(context.Context, LoadSessionRequest) (LoadSessionResponse, error) {
	return LoadSessionResponse{}, MethodNotFound()
}

func (UnimplementedAgent) SetSessionMode(context.Context, SetSessionModeRequest) (SetSessionModeResponse, error) {
	return SetSessionModeResponse{}, MethodNotFound()
}

func (UnimplementedAgent) SetSessionModel(context.Context, SetSessionModelRequest) (SetSessionModelResponse, error) {
	return SetSessionModelResponse{}, MethodNotFound()
}

func (UnimplementedAgent) Prompt(context.Context, PromptRequest) (PromptResponse, error) {
	return PromptResponse{}, MethodNotFound()
}

func (UnimplementedAgent) Cancel(context.Context, CancelNotification) error {
	return MethodNotFound()
}

func (UnimplementedAgent) ExtMethod(context.Context, ExtRequest) (ExtResponse, error) {
	return ExtResponse(json.RawMessage("null")), nil
}

func (UnimplementedAgent) ExtNotification(context.Context, ExtNotification) error {
	return nil
}
