package acp

import (
	"encoding/json"
	"path/filepath"
)

// Implementation 描述客户端或代理程序的版本信息。
type Implementation struct {
	Name    string  `json:"name"`
	Title   *string `json:"title,omitempty"`
	Version string  `json:"version"`
}

// InitializeRequest 对应 session 初始化请求。
type InitializeRequest struct {
	ProtocolVersion    ProtocolVersion    `json:"protocolVersion"`
	ClientCapabilities ClientCapabilities `json:"clientCapabilities"`
	ClientInfo         *Implementation    `json:"clientInfo,omitempty"`
	Meta               json.RawMessage    `json:"_meta,omitempty"`
}

// InitializeResponse 为初始化响应。
type InitializeResponse struct {
	ProtocolVersion   ProtocolVersion   `json:"protocolVersion"`
	AgentCapabilities AgentCapabilities `json:"agentCapabilities"`
	AuthMethods       []AuthMethod      `json:"authMethods,omitempty"`
	AgentInfo         *Implementation   `json:"agentInfo,omitempty"`
	Meta              json.RawMessage   `json:"_meta,omitempty"`
}

// AuthMethod 描述一种认证方式。
type AuthMethod struct {
	ID          AuthMethodID    `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}

// AuthenticateRequest 对应认证请求。
type AuthenticateRequest struct {
	MethodID AuthMethodID    `json:"methodId"`
	Meta     json.RawMessage `json:"_meta,omitempty"`
}

// AuthenticateResponse 对应认证响应。
type AuthenticateResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// AgentCapabilities 描述代理能力。
type AgentCapabilities struct {
	LoadSession        bool               `json:"loadSession,omitempty"`
	PromptCapabilities PromptCapabilities `json:"promptCapabilities"`
	McpCapabilities    McpCapabilities    `json:"mcpCapabilities"`
	Meta               json.RawMessage    `json:"_meta,omitempty"`
}

// PromptCapabilities 描述 prompt 能力。
type PromptCapabilities struct {
	Image           bool            `json:"image,omitempty"`
	Audio           bool            `json:"audio,omitempty"`
	EmbeddedContext bool            `json:"embeddedContext,omitempty"`
	Meta            json.RawMessage `json:"_meta,omitempty"`
}

// McpCapabilities 描述 MCP 能力。
type McpCapabilities struct {
	HTTP bool            `json:"http,omitempty"`
	SSE  bool            `json:"sse,omitempty"`
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// McpServerType 统一定义。
type McpServerType string

const (
	McpServerTypeHTTP  McpServerType = "http"
	McpServerTypeSSE   McpServerType = "sse"
	McpServerTypeStdio McpServerType = "stdio"
)

// HttpHeader HTTP 头。
type HttpHeader struct {
	Name  string          `json:"name"`
	Value string          `json:"value"`
	Meta  json.RawMessage `json:"_meta,omitempty"`
}

// EnvVariable 描述环境变量。
type EnvVariable struct {
	Name  string          `json:"name"`
	Value string          `json:"value"`
	Meta  json.RawMessage `json:"_meta,omitempty"`
}

// McpServer MCP 服务配置。
type McpServer struct {
	Type    McpServerType `json:"type"`
	Name    string        `json:"name"`
	URL     string        `json:"url,omitempty"`
	Headers []HttpHeader  `json:"headers,omitempty"`
	Command string        `json:"command,omitempty"`
	Args    []string      `json:"args,omitempty"`
	Env     []EnvVariable `json:"env,omitempty"`
}

// NewSessionRequest 对应 session/new 请求。
type NewSessionRequest struct {
	CWD        string          `json:"cwd"`
	McpServers []McpServer     `json:"mcpServers"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// NewSessionResponse 对应 session/new 响应。
type NewSessionResponse struct {
	SessionID SessionID          `json:"sessionId"`
	Modes     *SessionModeState  `json:"modes,omitempty"`
	Models    *SessionModelState `json:"models,omitempty"`
	Meta      json.RawMessage    `json:"_meta,omitempty"`
}

// LoadSessionRequest 对应 session/load 请求。
type LoadSessionRequest struct {
	McpServers []McpServer     `json:"mcpServers"`
	CWD        string          `json:"cwd"`
	SessionID  SessionID       `json:"sessionId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// LoadSessionResponse 。
type LoadSessionResponse struct {
	Modes  *SessionModeState  `json:"modes,omitempty"`
	Models *SessionModelState `json:"models,omitempty"`
	Meta   json.RawMessage    `json:"_meta,omitempty"`
}

// SessionModeState 描述模式状态。
type SessionModeState struct {
	CurrentModeID  SessionModeID   `json:"currentModeId"`
	AvailableModes []SessionMode   `json:"availableModes"`
	Meta           json.RawMessage `json:"_meta,omitempty"`
}

// SessionModeID 模式标识。
type SessionModeID string

// SessionMode 模式信息。
type SessionMode struct {
	ID          SessionModeID   `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}

// SetSessionModeRequest 对应 session/set_mode 请求。
type SetSessionModeRequest struct {
	SessionID SessionID       `json:"sessionId"`
	ModeID    SessionModeID   `json:"modeId"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// SetSessionModeResponse 对应响应。
type SetSessionModeResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// SessionModelState 描述模型状态。
type SessionModelState struct {
	CurrentModelID  ModelID         `json:"currentModelId"`
	AvailableModels []ModelInfo     `json:"availableModels"`
	Meta            json.RawMessage `json:"_meta,omitempty"`
}

// ModelID 模型标识。
type ModelID string

// ModelInfo 模型信息。
type ModelInfo struct {
	ModelID     ModelID         `json:"modelId"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}

// SetSessionModelRequest 对应 session/set_model 请求。
type SetSessionModelRequest struct {
	SessionID SessionID       `json:"sessionId"`
	ModelID   ModelID         `json:"modelId"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// SetSessionModelResponse 对应 session/set_model 响应。
type SetSessionModelResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// PromptRequest 对应 session/prompt 请求。
type PromptRequest struct {
	SessionID SessionID       `json:"sessionId"`
	Prompt    []ContentBlock  `json:"prompt"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// PromptResponse 对应 session/prompt 响应。
type PromptResponse struct {
	StopReason StopReason      `json:"stopReason"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// StopReason 结束原因。
type StopReason string

// StopReason 枚举。
const (
	StopReasonEndTurn         StopReason = "end_turn"
	StopReasonMaxTokens       StopReason = "max_tokens"
	StopReasonMaxTurnRequests StopReason = "max_turn_requests"
	StopReasonRefusal         StopReason = "refusal"
	StopReasonCancelled       StopReason = "cancelled"
)

// CancelNotification 对应 session/cancel 通知。
type CancelNotification struct {
	SessionID SessionID       `json:"sessionId"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// ExtRequest 自定义请求。
type ExtRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// ExtResponse 自定义响应。
type ExtResponse json.RawMessage

// ExtNotification 自定义通知。
type ExtNotification struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// EnsureCWD 规范化路径。
func (req *NewSessionRequest) EnsureCWD() {
	if req == nil || req.CWD == "" {
		return
	}
	req.CWD = filepath.Clean(req.CWD)
}
