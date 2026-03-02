package cloudpods

import (
	"bytes"
	"cloudpods-mcp-server/pkg/config"
	"cloudpods-mcp-server/pkg/tools"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
	Domain   string
	Project  string

	toolsReg *tools.Registry

	Token      string
	TokenType  string
	ExpireTime time.Time
	HTTPClient *http.Client
}

func NewClient(config *config.Config, toolsReg *tools.Registry) (*Client, error) {
	client := &Client{
		BaseURL:  config.CloudpodsAPI,
		Username: config.Username,
		Password: config.Password,
		Domain:   config.Domain,
		Project:  config.Project,

		toolsReg: toolsReg,

		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	err := client.Authenticate()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) RegisterAllTools() {
	c.RegisterListRegionsTool()
	c.RegisterGetRegionTool()
	c.RegisterListVPCsTool()
	c.RegisterListServersTool()
	c.RegisterGetServerTool()
	c.RegisterStartServerTool()
	c.RegisterStopServerTool()
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// 检查认证信息存在
	if err := c.ensureAuthenticated(ctx); err != nil {
		return fmt.Errorf("fail to authenticate: %w", err)
	}

	// 拼接完成http请求url
	url := fmt.Sprintf("%s%s%s", c.BaseURL, "/api", path)

	// 如果有请求体，就先将请求体转成json，然后加载到内存中
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("fail to serial the body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("fail to create the http request: %w", err)
	}

	// 将yunionauth令牌加到Header的Cookie段
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("yunionauth=%s", c.Token))

	// 发送请求
	zap.L().Debug("send http request to cloudpods",
		zap.String("method", method),
		zap.String("URL", url),
		zap.Any("body", body))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to send http request: %w", err)
	}
	defer resp.Body.Close()

	// 检查请求是否成功
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("request error (status %d): %s", resp.StatusCode, errResp.Error)
		}
		return fmt.Errorf("fail to request, status: %d", resp.StatusCode)
	}

	// 解析响应体
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("fail to decode the response body : %w", err)
		}
	}
	return nil
}
