package cloudpods

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// RegisterListRegionsTool 注册"获取区域列表"工具
func (c *Client) RegisterListRegionsTool() {
	tool := mcp.NewTool("list_regions",
		mcp.WithDescription("List all cloud regions"),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			regions, err := c.ListRegions(ctx)
			if err != nil {
				zap.L().Error("fail to list regions", zap.Any("err", err))
				return mcp.NewToolResultErrorFromErr("fail to list regions", err), nil
			}
			zap.L().Info("list regions", zap.Any("regions", regions))
			bytes, err := json.Marshal(regions)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to marshal regions", err), nil
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

// RegisterGetRegionTool 注册"获取区域详情"工具
func (c *Client) RegisterGetRegionTool() {
	tool := mcp.NewTool("get_region",
		mcp.WithDescription("Get the detail of a cloud region by ID"),
		mcp.WithString("region_id",
			mcp.Required(),
			mcp.Description("The ID of the region to retrieve"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			regionID, err := request.RequireString("region_id")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("missing required argument: region_id", err), nil
			}
			region, err := c.GetRegion(ctx, regionID)
			if err != nil {
				zap.L().Error("fail to get region", zap.String("region_id", regionID), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to get region", err), nil
			}
			bytes, err := json.Marshal(region)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to marshal region", err), nil
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

// RegisterListVPCsTool 注册"获取VPC列表"工具
func (c *Client) RegisterListVPCsTool() {
	tool := mcp.NewTool("list_vpcs",
		mcp.WithDescription("List VPCs, optionally filtered by region"),
		mcp.WithString("region_id",
			mcp.Description("Optional region ID to filter VPCs"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			regionID := request.GetString("region_id", "")
			vpcs, err := c.ListVPCs(ctx, regionID)
			if err != nil {
				zap.L().Error("fail to list VPCs", zap.String("region_id", regionID), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to list VPCs", err), nil
			}
			zap.L().Info("list VPCs", zap.Int("count", len(vpcs)))
			bytes, err := json.Marshal(vpcs)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to marshal VPCs", err), nil
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

// RegisterListServersTool 注册"获取虚拟机列表"工具
func (c *Client) RegisterListServersTool() {
	tool := mcp.NewTool("list_servers",
		mcp.WithDescription("List virtual machines (servers), with optional filters"),
		mcp.WithString("region_id",
			mcp.Description("Optional region ID to filter servers"),
		),
		mcp.WithString("status",
			mcp.Description("Optional server status filter (e.g. running, stopped)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of servers to return (default 0 means no limit)"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Number of servers to skip for pagination"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			options := &ListServerOptions{
				RegionID: request.GetString("region_id", ""),
				Status:   request.GetString("status", ""),
				Limit:    request.GetInt("limit", 0),
				Offset:   request.GetInt("offset", 0),
			}
			servers, err := c.ListServers(ctx, options)
			if err != nil {
				zap.L().Error("fail to list servers", zap.Any("options", options), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to list servers", err), nil
			}
			zap.L().Info("list servers", zap.Int("count", len(servers)))
			bytes, err := json.Marshal(servers)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to marshal servers", err), nil
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

// RegisterGetServerTool 注册"获取虚拟机详情"工具
func (c *Client) RegisterGetServerTool() {
	tool := mcp.NewTool("get_server",
		mcp.WithDescription("Get the detail of a virtual machine by ID"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("The ID of the server to retrieve"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serverID, err := request.RequireString("server_id")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("missing required argument: server_id", err), nil
			}
			server, err := c.GetServer(ctx, serverID)
			if err != nil {
				zap.L().Error("fail to get server", zap.String("server_id", serverID), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to get server", err), nil
			}
			bytes, err := json.Marshal(server)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to marshal server", err), nil
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

// RegisterStartServerTool 注册"启动虚拟机"工具
func (c *Client) RegisterStartServerTool() {
	tool := mcp.NewTool("start_server",
		mcp.WithDescription("Start a virtual machine by ID"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("The ID of the server to start"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serverID, err := request.RequireString("server_id")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("missing required argument: server_id", err), nil
			}
			if err := c.StartServer(ctx, serverID); err != nil {
				zap.L().Error("fail to start server", zap.String("server_id", serverID), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to start server", err), nil
			}
			zap.L().Info("server started", zap.String("server_id", serverID))
			return mcp.NewToolResultText("server started successfully"), nil
		})
}

// RegisterStopServerTool 注册"停止虚拟机"工具
func (c *Client) RegisterStopServerTool() {
	tool := mcp.NewTool("stop_server",
		mcp.WithDescription("Stop a virtual machine by ID"),
		mcp.WithString("server_id",
			mcp.Required(),
			mcp.Description("The ID of the server to stop"),
		),
	)

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			serverID, err := request.RequireString("server_id")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("missing required argument: server_id", err), nil
			}
			if err := c.StopServer(ctx, serverID); err != nil {
				zap.L().Error("fail to stop server", zap.String("server_id", serverID), zap.Error(err))
				return mcp.NewToolResultErrorFromErr("fail to stop server", err), nil
			}
			zap.L().Info("server stopped", zap.String("server_id", serverID))
			return mcp.NewToolResultText("server stopped successfully"), nil
		})
}
