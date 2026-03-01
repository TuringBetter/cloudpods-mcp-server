package tools

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

/*
 * 负责管理所有的工具函数，当server启动之前，就会将所有工具函数先注册到Registry中，再由Registry统一注册到mcp-server中。
 * 具体就是调用mcp-go的AddTool将tool和对应的handler建立映射关系，并且一同交由MCPServer管理
 */

type Registry struct {
	server *server.MCPServer
}

func NewRegistry(s *server.MCPServer) *Registry {
	return &Registry{
		server: s,
	}
}

func (r *Registry) RegisterTool(tool mcp.Tool, handlerFunc server.ToolHandlerFunc) {
	r.server.AddTool(tool, handlerFunc)
}
