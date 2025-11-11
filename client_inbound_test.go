package acp

import (
	"context"
	"encoding/json"
	"testing"
)

type testClient struct {
	UnimplementedClient
	permissionCalled bool
	lastPermission   RequestPermissionRequest
	sessionNotified  bool
	extMethodCalled  bool
	extNotifCalled   bool
}

func (c *testClient) RequestPermission(ctx context.Context, req RequestPermissionRequest) (RequestPermissionResponse, error) {
	c.permissionCalled = true
	c.lastPermission = req
	var optionPtr *PermissionOptionID
	if len(req.Options) > 0 {
		id := req.Options[0].ID
		optionPtr = &id
	}
	return RequestPermissionResponse{
		Outcome: RequestPermissionOutcome{
			Outcome: "selected",
			Option:  optionPtr,
		},
	}, nil
}

func (c *testClient) SessionNotification(ctx context.Context, note SessionNotification) error {
	c.sessionNotified = true
	return nil
}

func (c *testClient) ExtMethod(ctx context.Context, req ExtRequest) (ExtResponse, error) {
	c.extMethodCalled = true
	return ExtResponse(json.RawMessage(`"ok"`)), nil
}

func (c *testClient) ExtNotification(ctx context.Context, note ExtNotification) error {
	c.extNotifCalled = true
	return nil
}

func TestClientInboundHandleRequest(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}

	params := mustRawJSON(RequestPermissionRequest{
		SessionID: SessionID("s"),
		ToolCall: ToolCallUpdate{
			ID:     "tool",
			Status: "pending",
		},
		Options: []PermissionOption{
			{ID: PermissionOptionID("allow"), Name: "Allow", Kind: PermissionOptionKindAllowOnce},
		},
	})

	resp, errObj, ok := handler.handleRequest(context.Background(), ClientMethods.SessionRequestPermission, params)
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	result, ok := resp.(RequestPermissionResponse)
	if !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
	if result.Outcome.Option == nil || *result.Outcome.Option != PermissionOptionID("allow") {
		t.Fatalf("unexpected outcome: %+v", result)
	}
	if !client.permissionCalled {
		t.Fatalf("expected client permission to be called")
	}
}

func TestClientInboundHandleRequestInvalidParams(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}
	_, errObj, ok := handler.handleRequest(context.Background(), ClientMethods.SessionRequestPermission, json.RawMessage("[]"))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != ErrorCodeInvalidParams.Code {
		t.Fatalf("expected invalid params error, got %+v", errObj)
	}
}

func TestClientInboundHandleExtMethod(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}
	resp, errObj, ok := handler.handleRequest(context.Background(), "_custom", mustRawJSON(map[string]string{"key": "value"}))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !client.extMethodCalled {
		t.Fatalf("expected ext method to be called")
	}
	if _, ok := resp.(ExtResponse); !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
}

func TestClientInboundHandleUnknownMethod(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}
	_, errObj, ok := handler.handleRequest(context.Background(), "unknown/method", mustRawJSON(map[string]string{}))
	if !ok {
		t.Fatalf("expected ok")
	}
	if errObj.Code != ErrorCodeMethodNotFound.Code {
		t.Fatalf("expected method not found, got %+v", errObj)
	}
}

func TestClientInboundHandleNotification(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}
	errObj := handler.handleNotification(context.Background(), ClientMethods.SessionUpdate, mustRawJSON(SessionNotification{
		SessionID: SessionID("sess"),
		Update: SessionUpdate{
			Type: SessionUpdateTypeAgentMessageChunk,
		},
	}))
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !client.sessionNotified {
		t.Fatalf("expected session notification to be handled")
	}
}

func TestClientInboundHandleExtNotification(t *testing.T) {
	client := &testClient{}
	handler := &clientInboundHandler{client: client}
	errObj := handler.handleNotification(context.Background(), "_custom", mustRawJSON(map[string]string{"k": "v"}))
	if errObj.Code != 0 {
		t.Fatalf("unexpected error: %+v", errObj)
	}
	if !client.extNotifCalled {
		t.Fatalf("expected ext notification")
	}
}
