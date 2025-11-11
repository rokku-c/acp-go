package acp

import "encoding/json"

// SessionNotification 对应 session/update 通知。
type SessionNotification struct {
	SessionID SessionID       `json:"sessionId"`
	Update    SessionUpdate   `json:"update"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// SessionUpdateType 更新类型。
type SessionUpdateType string

const (
	SessionUpdateTypeAgentMessageChunk SessionUpdateType = "agent_message_chunk"
	SessionUpdateTypeUserMessageChunk  SessionUpdateType = "user_message_chunk"
	SessionUpdateTypeAgentThoughtChunk SessionUpdateType = "agent_thought_chunk"
	SessionUpdateTypeToolCall          SessionUpdateType = "tool_call"
	SessionUpdateTypeToolCallUpdate    SessionUpdateType = "tool_call_update"
	SessionUpdateTypePlan              SessionUpdateType = "plan"
	SessionUpdateTypeAvailableCommands SessionUpdateType = "available_commands_update"
	SessionUpdateTypeCurrentMode       SessionUpdateType = "current_mode_update"
)

// SessionUpdate 描述会话更新。
type SessionUpdate struct {
	Type           SessionUpdateType `json:"sessionUpdate"`
	Content        *ContentBlock     `json:"content,omitempty"`
	ToolCall       *ToolCall         `json:"toolCall,omitempty"`
	ToolCallUpdate *ToolCallUpdate   `json:"toolCallUpdate,omitempty"`
	Meta           json.RawMessage   `json:"_meta,omitempty"`
}

// ToolCall 描述一次工具调用。
type ToolCall struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Input json.RawMessage `json:"input,omitempty"`
	Meta  json.RawMessage `json:"_meta,omitempty"`
}

// ToolCallUpdate 描述工具调用更新。
type ToolCallUpdate struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Output json.RawMessage `json:"output,omitempty"`
	Meta   json.RawMessage `json:"_meta,omitempty"`
}

// AvailableCommand 描述可用命令。
type AvailableCommand struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputHint   string          `json:"inputHint,omitempty"`
	Meta        json.RawMessage `json:"_meta,omitempty"`
}

// CurrentModeUpdate 会话模式更新。
type CurrentModeUpdate struct {
	CurrentModeID SessionModeID   `json:"currentModeId"`
	Meta          json.RawMessage `json:"_meta,omitempty"`
}

// AvailableCommandsUpdate 描述可用命令更新。
type AvailableCommandsUpdate struct {
	AvailableCommands []AvailableCommand `json:"availableCommands"`
	Meta              json.RawMessage    `json:"_meta,omitempty"`
}

// RequestPermissionRequest 请求权限。
type RequestPermissionRequest struct {
	SessionID SessionID          `json:"sessionId"`
	ToolCall  ToolCallUpdate     `json:"toolCall"`
	Options   []PermissionOption `json:"options"`
	Meta      json.RawMessage    `json:"_meta,omitempty"`
}

// PermissionOption 权限选项。
type PermissionOption struct {
	ID   PermissionOptionID   `json:"optionId"`
	Name string               `json:"name"`
	Kind PermissionOptionKind `json:"kind"`
	Meta json.RawMessage      `json:"_meta,omitempty"`
}

// PermissionOptionKind 选项类型。
type PermissionOptionKind string

const (
	PermissionOptionKindAllowOnce    PermissionOptionKind = "allow_once"
	PermissionOptionKindAllowAlways  PermissionOptionKind = "allow_always"
	PermissionOptionKindRejectOnce   PermissionOptionKind = "reject_once"
	PermissionOptionKindRejectAlways PermissionOptionKind = "reject_always"
)

// RequestPermissionOutcome 描述权限请求结果。
type RequestPermissionOutcome struct {
	Outcome string              `json:"outcome"`
	Option  *PermissionOptionID `json:"optionId,omitempty"`
}

// RequestPermissionResponse 对应响应。
type RequestPermissionResponse struct {
	Outcome RequestPermissionOutcome `json:"outcome"`
	Meta    json.RawMessage          `json:"_meta,omitempty"`
}

// WriteTextFileRequest 请求写文件。
type WriteTextFileRequest struct {
	SessionID SessionID       `json:"sessionId"`
	Path      string          `json:"path"`
	Content   string          `json:"content"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// WriteTextFileResponse 写文件响应。
type WriteTextFileResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// ReadTextFileRequest 读文件请求。
type ReadTextFileRequest struct {
	SessionID SessionID       `json:"sessionId"`
	Path      string          `json:"path"`
	Line      *uint32         `json:"line,omitempty"`
	Limit     *uint32         `json:"limit,omitempty"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// ReadTextFileResponse 读文件响应。
type ReadTextFileResponse struct {
	Content string          `json:"content"`
	Meta    json.RawMessage `json:"_meta,omitempty"`
}

// CreateTerminalRequest 创建终端。
type CreateTerminalRequest struct {
	SessionID SessionID       `json:"sessionId"`
	Command   string          `json:"command"`
	Args      []string        `json:"args,omitempty"`
	Env       []EnvVariable   `json:"env,omitempty"`
	CWD       string          `json:"cwd,omitempty"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

// CreateTerminalResponse 创建终端响应。
type CreateTerminalResponse struct {
	TerminalID TerminalID      `json:"terminalId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// TerminalOutputRequest 终端输出请求。
type TerminalOutputRequest struct {
	SessionID  SessionID       `json:"sessionId"`
	TerminalID TerminalID      `json:"terminalId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// TerminalExitStatus 终端退出状态。
type TerminalExitStatus struct {
	ExitCode *uint32         `json:"exitCode,omitempty"`
	Signal   *string         `json:"signal,omitempty"`
	Meta     json.RawMessage `json:"_meta,omitempty"`
}

// TerminalOutputResponse 终端输出响应。
type TerminalOutputResponse struct {
	Output     string              `json:"output"`
	Truncated  bool                `json:"truncated"`
	ExitStatus *TerminalExitStatus `json:"exitStatus,omitempty"`
	Meta       json.RawMessage     `json:"_meta,omitempty"`
}

// ReleaseTerminalRequest 释放终端请求。
type ReleaseTerminalRequest struct {
	SessionID  SessionID       `json:"sessionId"`
	TerminalID TerminalID      `json:"terminalId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// ReleaseTerminalResponse 释放终端响应。
type ReleaseTerminalResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// KillTerminalCommandRequest 终止终端命令请求。
type KillTerminalCommandRequest struct {
	SessionID  SessionID       `json:"sessionId"`
	TerminalID TerminalID      `json:"terminalId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// KillTerminalCommandResponse 终止终端命令响应。
type KillTerminalCommandResponse struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

// WaitForTerminalExitRequest 等待终端退出请求。
type WaitForTerminalExitRequest struct {
	SessionID  SessionID       `json:"sessionId"`
	TerminalID TerminalID      `json:"terminalId"`
	Meta       json.RawMessage `json:"_meta,omitempty"`
}

// WaitForTerminalExitResponse 等待终端退出响应。
type WaitForTerminalExitResponse struct {
	ExitStatus TerminalExitStatus `json:"exitStatus"`
	Meta       json.RawMessage    `json:"_meta,omitempty"`
}

// ClientCapabilities 描述客户端能力。
type ClientCapabilities struct {
	FS       FileSystemCapability `json:"fs"`
	Terminal bool                 `json:"terminal,omitempty"`
	Meta     json.RawMessage      `json:"_meta,omitempty"`
}

// FileSystemCapability 描述文件系统能力。
type FileSystemCapability struct {
	ReadTextFile  bool            `json:"readTextFile,omitempty"`
	WriteTextFile bool            `json:"writeTextFile,omitempty"`
	Meta          json.RawMessage `json:"_meta,omitempty"`
}
