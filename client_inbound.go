package acp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type clientInboundHandler struct {
	client Client
}

func (h *clientInboundHandler) handleRequest(ctx context.Context, method string, params json.RawMessage) (any, Error, bool) {
	switch method {
	case ClientMethods.SessionRequestPermission:
		req, err := decodeParams[RequestPermissionRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.RequestPermission(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.FSWriteTextFile:
		req, err := decodeParams[WriteTextFileRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.WriteTextFile(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.FSReadTextFile:
		req, err := decodeParams[ReadTextFileRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.ReadTextFile(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.TerminalCreate:
		req, err := decodeParams[CreateTerminalRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.CreateTerminal(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.TerminalOutput:
		req, err := decodeParams[TerminalOutputRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.TerminalOutput(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.TerminalRelease:
		req, err := decodeParams[ReleaseTerminalRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.ReleaseTerminal(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.TerminalWaitForExit:
		req, err := decodeParams[WaitForTerminalExitRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.WaitForTerminalExit(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case ClientMethods.TerminalKill:
		req, err := decodeParams[KillTerminalCommandRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.client.KillTerminalCommand(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	default:
		if strings.HasPrefix(method, "_") {
			req := ExtRequest{
				Method: strings.TrimPrefix(method, "_"),
				Params: params,
			}
			resp, callErr := h.client.ExtMethod(ctx, req)
			if callErr != nil {
				return nil, IntoInternalError(callErr), true
			}
			return resp, Error{}, true
		}
		return nil, MethodNotFound(), true
	}
}

func (h *clientInboundHandler) handleNotification(ctx context.Context, method string, params json.RawMessage) Error {
	switch method {
	case ClientMethods.SessionUpdate:
		req, err := decodeParams[SessionNotification](params)
		if err != nil {
			return InvalidParams().WithData(err.Error())
		}
		if callErr := h.client.SessionNotification(ctx, req); callErr != nil {
			return IntoInternalError(callErr)
		}
		return Error{}
	default:
		if strings.HasPrefix(method, "_") {
			notification := ExtNotification{
				Method: strings.TrimPrefix(method, "_"),
				Params: params,
			}
			if callErr := h.client.ExtNotification(ctx, notification); callErr != nil {
				return IntoInternalError(callErr)
			}
			return Error{}
		}
		return MethodNotFound()
	}
}

func decodeParams[T any](params json.RawMessage) (T, error) {
	var zero T
	if len(params) == 0 {
		return zero, fmt.Errorf("params must not be empty")
	}
	if err := json.Unmarshal(params, &zero); err != nil {
		return zero, err
	}
	return zero, nil
}
