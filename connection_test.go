package acp

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"
)

type mockAgent struct {
	initializeFunc   func(context.Context, InitializeRequest) (InitializeResponse, error)
	authenticateFunc func(context.Context, AuthenticateRequest) (AuthenticateResponse, error)
	newSessionFunc   func(context.Context, NewSessionRequest) (NewSessionResponse, error)
	loadSessionFunc  func(context.Context, LoadSessionRequest) (LoadSessionResponse, error)
	setModeFunc      func(context.Context, SetSessionModeRequest) (SetSessionModeResponse, error)
	setModelFunc     func(context.Context, SetSessionModelRequest) (SetSessionModelResponse, error)
	promptFunc       func(context.Context, PromptRequest) (PromptResponse, error)
	cancelCh         chan CancelNotification
}

func (m *mockAgent) Initialize(ctx context.Context, req InitializeRequest) (InitializeResponse, error) {
	if m.initializeFunc != nil {
		return m.initializeFunc(ctx, req)
	}
	return InitializeResponse{}, nil
}

func (m *mockAgent) Authenticate(ctx context.Context, req AuthenticateRequest) (AuthenticateResponse, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, req)
	}
	return AuthenticateResponse{}, nil
}

func (m *mockAgent) NewSession(ctx context.Context, req NewSessionRequest) (NewSessionResponse, error) {
	if m.newSessionFunc != nil {
		return m.newSessionFunc(ctx, req)
	}
	return NewSessionResponse{}, nil
}

func (m *mockAgent) LoadSession(ctx context.Context, req LoadSessionRequest) (LoadSessionResponse, error) {
	if m.loadSessionFunc != nil {
		return m.loadSessionFunc(ctx, req)
	}
	return LoadSessionResponse{}, nil
}

func (m *mockAgent) SetSessionMode(ctx context.Context, req SetSessionModeRequest) (SetSessionModeResponse, error) {
	if m.setModeFunc != nil {
		return m.setModeFunc(ctx, req)
	}
	return SetSessionModeResponse{}, nil
}

func (m *mockAgent) SetSessionModel(ctx context.Context, req SetSessionModelRequest) (SetSessionModelResponse, error) {
	if m.setModelFunc != nil {
		return m.setModelFunc(ctx, req)
	}
	return SetSessionModelResponse{}, nil
}

func (m *mockAgent) Prompt(ctx context.Context, req PromptRequest) (PromptResponse, error) {
	if m.promptFunc != nil {
		return m.promptFunc(ctx, req)
	}
	return PromptResponse{}, nil
}

func (m *mockAgent) Cancel(ctx context.Context, note CancelNotification) error {
	select {
	case m.cancelCh <- note:
	default:
	}
	return nil
}

func (m *mockAgent) ExtMethod(context.Context, ExtRequest) (ExtResponse, error) {
	return ExtResponse(json.RawMessage(`"ok"`)), nil
}

func (m *mockAgent) ExtNotification(context.Context, ExtNotification) error {
	return nil
}

type mockClient struct {
	requestPermissionFunc func(context.Context, RequestPermissionRequest) (RequestPermissionResponse, error)
	sessionUpdateCh       chan SessionNotification
}

func (m *mockClient) RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
	if m.requestPermissionFunc != nil {
		return m.requestPermissionFunc(ctx, req)
	}
	return RequestPermissionResponse{}, nil
}

func (m *mockClient) SessionNotification(ctx context.Context, note SessionNotification) error {
	select {
	case m.sessionUpdateCh <- note:
	default:
	}
	return nil
}

func (m *mockClient) WriteTextFile(context.Context, WriteTextFileRequest) (WriteTextFileResponse, error) {
	return WriteTextFileResponse{}, nil
}

func (m *mockClient) ReadTextFile(context.Context, ReadTextFileRequest) (ReadTextFileResponse, error) {
	return ReadTextFileResponse{Content: "example"}, nil
}

func (m *mockClient) CreateTerminal(context.Context, CreateTerminalRequest) (CreateTerminalResponse, error) {
	return CreateTerminalResponse{TerminalID: TerminalID("term-1")}, nil
}

func (m *mockClient) TerminalOutput(context.Context, TerminalOutputRequest) (TerminalOutputResponse, error) {
	return TerminalOutputResponse{Output: "out", Truncated: false}, nil
}

func (m *mockClient) ReleaseTerminal(context.Context, ReleaseTerminalRequest) (ReleaseTerminalResponse, error) {
	return ReleaseTerminalResponse{}, nil
}

func (m *mockClient) WaitForTerminalExit(context.Context, WaitForTerminalExitRequest) (WaitForTerminalExitResponse, error) {
	return WaitForTerminalExitResponse{
		ExitStatus: TerminalExitStatus{ExitCode: ptr(uint32(0))},
	}, nil
}

func (m *mockClient) KillTerminalCommand(context.Context, KillTerminalCommandRequest) (KillTerminalCommandResponse, error) {
	return KillTerminalCommandResponse{}, nil
}

func (m *mockClient) ExtMethod(context.Context, ExtRequest) (ExtResponse, error) {
	return ExtResponse(json.RawMessage(`"client"`)), nil
}

func (m *mockClient) ExtNotification(context.Context, ExtNotification) error {
	return nil
}

func ptr[T any](v T) *T {
	return &v
}

func TestClientAgentRoundtrip(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clientToAgentReader, clientToAgentWriter := io.Pipe()
	agentToClientReader, agentToClientWriter := io.Pipe()

	agent := &mockAgent{
		cancelCh: make(chan CancelNotification, 1),
	}
	client := &mockClient{
		sessionUpdateCh: make(chan SessionNotification, 1),
	}

	agent.newSessionFunc = func(_ context.Context, req NewSessionRequest) (NewSessionResponse, error) {
		if req.CWD == "" {
			t.Fatalf("expected cwd")
		}
		return NewSessionResponse{
			SessionID: SessionID("sess"),
			Models: &SessionModelState{
				CurrentModelID: ModelID("gpt-5/medium"),
				AvailableModels: []ModelInfo{
					{ModelID: ModelID("gpt-5/high"), Name: "gpt-5 (high)"},
				},
			},
		}, nil
	}

	modelSet := false
	agent.setModelFunc = func(_ context.Context, req SetSessionModelRequest) (SetSessionModelResponse, error) {
		if req.ModelID != ModelID("gpt-5/high") {
			t.Fatalf("unexpected model id: %s", req.ModelID)
		}
		modelSet = true
		return SetSessionModelResponse{}, nil
	}

	clientConn := NewClientSideConnection(ctx, client, clientToAgentWriter, agentToClientReader)
	agentConn := NewAgentSideConnection(ctx, agent, agentToClientWriter, clientToAgentReader)

	t.Cleanup(func() {
		clientConn.Close()
		agentConn.Close()
		clientToAgentWriter.Close()
		clientToAgentReader.Close()
		agentToClientWriter.Close()
		agentToClientReader.Close()
	})

	agent.initializeFunc = func(_ context.Context, req InitializeRequest) (InitializeResponse, error) {
		if req.ProtocolVersion.Value() != ProtocolVersionV1.Value() {
			t.Fatalf("unexpected protocol version: %d", req.ProtocolVersion.Value())
		}
		return InitializeResponse{
			ProtocolVersion: ProtocolVersionV1,
			AgentCapabilities: AgentCapabilities{
				LoadSession: true,
			},
		}, nil
	}

	initResp, err := clientConn.Initialize(ctx, InitializeRequest{
		ProtocolVersion: ProtocolVersionV1,
	})
	if err != nil {
		t.Fatalf("initialize failed: %v", err)
	}
	if initResp.ProtocolVersion.Value() != ProtocolVersionV1.Value() {
		t.Fatalf("unexpected init response version: %d", initResp.ProtocolVersion.Value())
	}

	sessionResp, err := clientConn.NewSession(ctx, NewSessionRequest{
		CWD: "/tmp",
	})
	if err != nil {
		t.Fatalf("new session failed: %v", err)
	}
	if sessionResp.SessionID != SessionID("sess") {
		t.Fatalf("unexpected session id %s", sessionResp.SessionID)
	}
	if sessionResp.Models == nil {
		t.Fatalf("expected models")
	}
	if _, err := clientConn.SetSessionModel(ctx, SetSessionModelRequest{
		SessionID: sessionResp.SessionID,
		ModelID:   sessionResp.Models.AvailableModels[0].ModelID,
	}); err != nil {
		t.Fatalf("set session model failed: %v", err)
	}
	if !modelSet {
		t.Fatalf("expected set model to be called")
	}

	client.requestPermissionFunc = func(_ context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
		if req.ToolCall.ID == "" {
			t.Fatalf("missing tool call id")
		}
		return RequestPermissionResponse{
			Outcome: RequestPermissionOutcome{
				Outcome: "selected",
				Option:  ptr(PermissionOptionID("allow")),
			},
		}, nil
	}

	resp, err := agentConn.RequestPermission(ctx, RequestPermissionRequest{
		SessionID: SessionID("sess"),
		ToolCall: ToolCallUpdate{
			ID:     "tool-1",
			Status: "pending",
		},
		Options: []PermissionOption{
			{ID: PermissionOptionID("allow"), Name: "Allow", Kind: PermissionOptionKindAllowOnce},
		},
	})
	if err != nil {
		t.Fatalf("request permission failed: %v", err)
	}
	if resp.Outcome.Option == nil || *resp.Outcome.Option != PermissionOptionID("allow") {
		t.Fatalf("unexpected permission outcome: %+v", resp)
	}

	note := SessionNotification{
		SessionID: SessionID("sess"),
		Update: SessionUpdate{
			Type:    SessionUpdateTypeAgentMessageChunk,
			Content: ptr(NewTextContentBlock("hi")),
		},
	}
	if err := agentConn.SessionNotification(ctx, note); err != nil {
		t.Fatalf("session notification failed: %v", err)
	}

	select {
	case got := <-client.sessionUpdateCh:
		if got.SessionID != note.SessionID {
			t.Fatalf("unexpected session id: %s", got.SessionID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for session notification")
	}

	if err := clientConn.Cancel(ctx, CancelNotification{SessionID: SessionID("sess")}); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}

	select {
	case got := <-agent.cancelCh:
		if got.SessionID != "sess" {
			t.Fatalf("unexpected cancel session id: %s", got.SessionID)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for cancel notification")
	}

	stream := agentConn.Subscribe()
	if err := agentConn.ExtNotification(ctx, "custom", json.RawMessage(`{"value":1}`)); err != nil {
		t.Fatalf("ext notification failed: %v", err)
	}

	streamCtx, streamCancel := context.WithTimeout(ctx, 2*time.Second)
	defer streamCancel()
	msg, err := stream.Recv(streamCtx)
	if err != nil {
		t.Fatalf("failed to receive stream message: %v", err)
	}
	if msg.Content.Type != "notification" {
		t.Fatalf("unexpected stream message type: %s", msg.Content.Type)
	}
}
