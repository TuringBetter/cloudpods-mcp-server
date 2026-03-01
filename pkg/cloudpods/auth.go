package cloudpods

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

/*
 * 如果要想通过http调用cloudpods的后端接口，就要先使用用户名和密码，通过接口/v1/auth/login获得yunnionauth token，
 * 这一步在auth模块中完成，并且将yunionauth token缓存起来，在之后的每一次请求中都会把这个token添加到Header的Cookie段中
 */

func (c *Client) Authenticate() error {
	authReq := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
		"domain":   c.Domain,
	}

	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/login", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to authenticate, status:%s, response body: %s", resp.StatusCode, string(body))
	}

	// 从响应地Header中获取yunionauth token令牌
	var yunionAuthToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "yunionauth" {
			yunionAuthToken = cookie.Value
			break
		}
	}

	if yunionAuthToken == "" {
		return fmt.Errorf("failed to authenticate, no yunionauth token found")
	}

	c.Token = yunionAuthToken
	return nil
}

func (c *Client) ensureAuthenticated(ctx context.Context) error {
	if c.Token == "" {
		zap.L().Error("no yunionauth token found, re-authenticating")
		return c.Authenticate()
	}
	return nil
}
