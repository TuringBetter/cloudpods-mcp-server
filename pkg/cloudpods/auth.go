package cloudpods

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

/*
 * 如果要想通过http调用cloudpods的后端接口，就要先使用用户名和密码，通过接口/v1/auth/login获得yunnionauth token，
 * 这一步在auth模块中完成，并且将yunionauth token缓存起来，在之后的每一次请求中都会把这个token添加到Header的Cookie段中
 */

// tokenExpiryBuffer 是 token 即将过期的缓冲时间，提前刷新以避免在请求中途 token 过期
const tokenExpiryBuffer = 5 * time.Minute

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
		return fmt.Errorf("failed to authenticate, status:%d, response body: %s", resp.StatusCode, string(body))
	}

	// 从响应的 Header 中获取 yunionauth token 令牌及其过期时间
	var yunionAuthToken string
	var expireTime time.Time

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "yunionauth" {
			yunionAuthToken = cookie.Value
			// 优先使用 cookie 中的 Expires 字段
			if !cookie.Expires.IsZero() {
				expireTime = cookie.Expires
			} else {
				// 若 cookie 没有 Expires，则用 MaxAge（秒数）推算
				if cookie.MaxAge > 0 {
					expireTime = time.Now().Add(time.Duration(cookie.MaxAge) * time.Second)
				} else {
					// 兜底：默认 1 小时后过期
					expireTime = time.Now().Add(1 * time.Hour)
				}
			}
			break
		}
	}

	if yunionAuthToken == "" {
		return fmt.Errorf("failed to authenticate, no yunionauth token found")
	}

	c.Token = yunionAuthToken
	c.ExpireTime = expireTime
	zap.L().Info("authenticated successfully", zap.Time("expire_at", expireTime))
	return nil
}

// ensureAuthenticated 检查 token 是否存在，以及是否即将过期（5 分钟内），
// 如果满足任一条件则主动重新认证。
func (c *Client) ensureAuthenticated(ctx context.Context) error {
	// token 为空，必须认证
	if c.Token == "" {
		zap.L().Info("no yunionauth token found, authenticating...")
		return c.Authenticate()
	}

	// token 即将在 tokenExpiryBuffer 内过期，提前刷新
	if !c.ExpireTime.IsZero() && time.Now().Add(tokenExpiryBuffer).After(c.ExpireTime) {
		zap.L().Info("yunionauth token is about to expire, refreshing...",
			zap.Time("expire_at", c.ExpireTime))
		return c.Authenticate()
	}

	return nil
}
