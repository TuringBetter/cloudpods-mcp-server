package cloudpods

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func (c *Client) RegisterListRegionsTool() {
	tool := mcp.NewTool("list_regions", mcp.WithDescription("List All Regions"))

	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			regions, err := c.ListRegions(ctx)
			if err != nil {
				zap.L().Error("fail to list regions", zap.Any("err", err))
				return nil, err
			}
			zap.L().Info("list regions", zap.Any("regions", regions))
			bytes, err := json.Marshal(regions)
			if err != nil {
				return nil, err
			}
			return mcp.NewToolResultText(string(bytes)), nil
		})
}

func (c *Client) RegisterGetRegionsTool() {
	tool := mcp.NewTool("get_regions", mcp.WithDescription("Get Region Detail"))
	c.toolsReg.RegisterTool(tool,
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_, err := c.GetRegion(ctx, "testId")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("fail to list regions", err),
					nil
			}
			return mcp.NewToolResultText("ok"), nil
		})
}
