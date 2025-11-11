package acp

// AgentMethodNames 定义客户端发送的 Agent 侧方法名。
type AgentMethodNames struct {
	Initialize      string
	Authenticate    string
	SessionNew      string
	SessionLoad     string
	SessionSetMode  string
	SessionSetModel string
	SessionPrompt   string
	SessionCancel   string
}

// ClientMethodNames 定义 Agent 发送给客户端的方法名。
type ClientMethodNames struct {
	SessionRequestPermission string
	SessionUpdate            string
	FSWriteTextFile          string
	FSReadTextFile           string
	TerminalCreate           string
	TerminalOutput           string
	TerminalRelease          string
	TerminalWaitForExit      string
	TerminalKill             string
}

// 协议约定的方法名常量。
var (
	AgentMethods = AgentMethodNames{
		Initialize:      "initialize",
		Authenticate:    "authenticate",
		SessionNew:      "session/new",
		SessionLoad:     "session/load",
		SessionSetMode:  "session/set_mode",
		SessionSetModel: "session/set_model",
		SessionPrompt:   "session/prompt",
		SessionCancel:   "session/cancel",
	}
	ClientMethods = ClientMethodNames{
		SessionRequestPermission: "session/request_permission",
		SessionUpdate:            "session/update",
		FSWriteTextFile:          "fs/write_text_file",
		FSReadTextFile:           "fs/read_text_file",
		TerminalCreate:           "terminal/create",
		TerminalOutput:           "terminal/output",
		TerminalRelease:          "terminal/release",
		TerminalWaitForExit:      "terminal/wait_for_exit",
		TerminalKill:             "terminal/kill",
	}
)
