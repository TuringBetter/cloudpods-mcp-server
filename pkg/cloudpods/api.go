package cloudpods

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListRegions 获取区域列表
func (c *Client) ListRegions(ctx context.Context) ([]Region, error) {
	var result struct {
		Data  []Region `json:"data"`
		Limit int      `json:"limit"`
		Total int      `json:"total"`
	}
	err := c.doRequest(ctx, http.MethodGet, "/v2/cloudregions", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("获取区域列表失败: %w", err)
	}
	return result.Data, nil
}

// GetRegion 获取区域详情
func (c *Client) GetRegion(ctx context.Context, regionID string) (*Region, error) {
	var result struct {
		Region Region `json:"region"`
	}
	path := fmt.Sprintf("/v1/regions/%s", regionID)
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("获取区域详情失败: %w", err)
	}
	return &result.Region, nil
}

// ListVPCs 获取VPC列表
func (c *Client) ListVPCs(ctx context.Context, regionID string) ([]VPC, error) {
	var result struct {
		VPCs []VPC `json:"vpcs"`
	}
	query := url.Values{}
	if regionID != "" {
		query.Set("region_id", regionID)
	}
	path := fmt.Sprintf("/v1/vpcs?%s", query.Encode())
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("获取VPC列表失败: %w", err)
	}
	return result.VPCs, nil
}

// ListServers 获取虚拟机列表
func (c *Client) ListServers(ctx context.Context, options *ListServerOptions) ([]Server, error) {
	var result struct {
		Servers []Server `json:"servers"`
	}
	query := url.Values{}
	if options.RegionID != "" {
		query.Set("region_id", options.RegionID)
	}
	if options.Status != "" {
		query.Set("status", options.Status)
	}
	if options.Limit > 0 {
		query.Set("limit", strconv.Itoa(options.Limit))
	}
	if options.Offset > 0 {
		query.Set("offset", strconv.Itoa(options.Offset))
	}
	path := fmt.Sprintf("/api/v2/servers?%s", query.Encode())
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("获取虚拟机列表失败: %w", err)
	}
	return result.Servers, nil
}

// GetServer 获取虚拟机详情
func (c *Client) GetServer(ctx context.Context, serverID string) (*Server, error) {
	var result struct {
		Server Server `json:"server"`
	}
	path := fmt.Sprintf("/v1/servers/%s", serverID)
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("获取虚拟机详情失败: %w", err)
	}
	return &result.Server, nil
}

// StartServer 启动虚拟机
func (c *Client) StartServer(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("/v1/servers/%s/start", serverID)
	err := c.doRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("启动虚拟机失败: %w", err)
	}

	return nil
}

// StopServer 停止虚拟机
func (c *Client) StopServer(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("/v1/servers/%s/stop", serverID)
	err := c.doRequest(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return fmt.Errorf("停止虚拟机失败: %w", err)
	}
	return nil
}
