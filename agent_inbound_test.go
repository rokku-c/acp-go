package acp

import (
	"context"
	"encoding/json"
	"testing"
)

type testAgent struct {
	UnimplementedAgent
	initializeCalled bool
	cancelCalled     bool
	extMethodCalled  bool
	extNotifCalled   bool
	setModelCalled   bool
}

func (a *testAgent) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	a.initializeCalled = true
	return InitializeResponse{
		ProtocolVersion: req.ProtocolVersion,
	}, nil
}

func (a *testAgent) Cancel(ctx context.Context, note CancelNotification) error {
	a.cancelCalled = true
	return nil
}

func (a *testAgent) ExtMethod(ctx context.Context, req ExtRequest) (ExtResponse, error) {
	a.extMethodCalled = true
	return ExtResponse(json.RawMessage(`{"ok":true}`)), nil
}

func (a *testAgent) ExtNotification(ctx context.Context, note ExtNotification) error {
	a.extNotifCalled = true
	return nil
}

func (a *testAgent) SetSessionModel(ctx context.Context, req SetSessionModelRequest) (SetSessionModelResponse, error) {
	a.setModelCalled = true
	if req.ModelID == "" {
		return SetSessionModelResponse{}, MethodNotFound()
	}
	return SetSessionModelResponse{}, nil
}

func TestAgentInboundHandleRequest(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}

	resp, errObj, ok := handler.handleRequest(context.Background(), AgentMethods.Initialize, mustRawJSON(InitializeRequest{
		ProtocolVersion: ProtocolVersionV1,
	}))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	result, ok := resp.(InitializeResponse)
	if !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
	if !agent.initializeCalled {
		t.Fatalf("expected initialize to be called")
	}
	if result.ProtocolVersion.Value() != ProtocolVersionV1.Value() {
		t.Fatalf("unexpected protocol version: %d", result.ProtocolVersion.Value())
	}
}

func TestAgentInboundHandleInvalidParams(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}
	_, errObj, ok := handler.handleRequest(context.Background(), AgentMethods.Initialize, json.RawMessage("[]"))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != ErrorCodeInvalidParams.Code {
		t.Fatalf("expected invalid params error, got %+v", errObj)
	}
}

func TestAgentInboundHandleExtMethod(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}
	resp, errObj, ok := handler.handleRequest(context.Background(), "_ext", mustRawJSON(map[string]string{"a": "b"}))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !agent.extMethodCalled {
		t.Fatalf("expected ext method to be called")
	}
	if _, ok := resp.(ExtResponse); !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
}

func TestAgentInboundHandleNotification(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}
	errObj := handler.handleNotification(context.Background(), AgentMethods.SessionCancel, mustRawJSON(CancelNotification{
		SessionID: SessionID("sess"),
	}))
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !agent.cancelCalled {
		t.Fatalf("expected cancel to be handled")
	}
}

func TestAgentInboundHandleExtNotification(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}
	errObj := handler.handleNotification(context.Background(), "_ext", mustRawJSON(map[string]string{"n": "v"}))
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !agent.extNotifCalled {
		t.Fatalf("expected ext notification to be handled")
	}
}

func TestAgentInboundHandleSetModel(t *testing.T) {
	agent := &testAgent{}
	handler := &agentInboundHandler{agent: agent}
	resp, errObj, ok := handler.handleRequest(context.Background(), AgentMethods.SessionSetModel, mustRawJSON(SetSessionModelRequest{
		SessionID: SessionID("sess"),
		ModelID:   ModelID("gpt-5/high"),
	}))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !agent.setModelCalled {
		t.Fatalf("expected set model to be called")
	}
	if _, ok := resp.(SetSessionModelResponse); !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
}
