package api

import (
	"cloudpods-mcp-server/pkg/cloudpods"
	"cloudpods-mcp-server/pkg/config"
	"cloudpods-mcp-server/pkg/tools"

	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	server   *server.MCPServer // 指向 MCPServer 实例（第三方库）
	cpClient *cloudpods.Client // 指向 cloudpods 客户端实例
	toolsReg *tools.Registry   // 指向工具注册表实例
	config   *config.Config
}

func NewMCPServer(config *config.Config) (*MCPServer, error) {
	s := server.NewMCPServer(
		"cloudpods",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery())

	toolsReg := tools.NewRegistry(s)

	client, err := cloudpods.NewClient(config, toolsReg)
	if err != nil {
		return nil, err
	}
	return &MCPServer{
		cpClient: client,
		config:   config,
		server:   s,
		toolsReg: toolsReg,
	}, nil
}

func (s *MCPServer) Start() error {
	s.registerAllTools()
	return server.ServeStdio(s.server)
}

func (s *MCPServer) registerAllTools() {
	s.cpClient.RegisterAllTools()
}
