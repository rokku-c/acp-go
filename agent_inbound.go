package acp

import (
	"context"
	"encoding/json"
	"strings"
)

type agentInboundHandler struct {
	agent Agent
}

func (h *agentInboundHandler) handleRequest(ctx context.Context, method string, params json.RawMessage) (any, Error, bool) {
	switch method {
	case AgentMethods.Initialize:
		req, err := decodeParams[InitializeRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.Initialize(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.Authenticate:
		req, err := decodeParams[AuthenticateRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.Authenticate(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.SessionNew:
		req, err := decodeParams[NewSessionRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		req.EnsureCWD()
		resp, callErr := h.agent.NewSession(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.SessionLoad:
		req, err := decodeParams[LoadSessionRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.LoadSession(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.SessionSetMode:
		req, err := decodeParams[SetSessionModeRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.SetSessionMode(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.SessionSetModel:
		req, err := decodeParams[SetSessionModelRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.SetSessionModel(ctx, req)
		if callErr != nil {
			return nil, IntoInternalError(callErr), true
		}
		return resp, Error{}, true
	case AgentMethods.SessionPrompt:
		req, err := decodeParams[PromptRequest](params)
		if err != nil {
			return nil, InvalidParams().WithData(err.Error()), true
		}
		resp, callErr := h.agent.Prompt(ctx, req)
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
			resp, callErr := h.agent.ExtMethod(ctx, req)
			if callErr != nil {
				return nil, IntoInternalError(callErr), true
			}
			return resp, Error{}, true
		}
		return nil, MethodNotFound(), true
	}
}

func (h *agentInboundHandler) handleNotification(ctx context.Context, method string, params json.RawMessage) Error {
	switch method {
	case AgentMethods.SessionCancel:
		req, err := decodeParams[CancelNotification](params)
		if err != nil {
			return InvalidParams().WithData(err.Error())
		}
		if callErr := h.agent.Cancel(ctx, req); callErr != nil {
			return IntoInternalError(callErr)
		}
		return Error{}
	default:
		if strings.HasPrefix(method, "_") {
			notification := ExtNotification{
				Method: strings.TrimPrefix(method, "_"),
				Params: params,
			}
			if callErr := h.agent.ExtNotification(ctx, notification); callErr != nil {
				return IntoInternalError(callErr)
			}
			return Error{}
		}
		return MethodNotFound()
	}
}
