package coze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/coze-dev/coze-go"
)

type CozeOAuthConfig struct {
	ClientType  string `json:"client_type"`   // 客户端类型
	ClientID    string `json:"client_id"`     // 客户端ID
	CozeWWWBase string `json:"coze_www_base"` // 官网地址
	CozeAPIBase string `json:"coze_api_base"` // 接口地址
	PrivateKey  string `json:"private_key"`   // 私钥
	PublicKeyID string `json:"public_key_id"` // 公钥ID
}

func loadConfig(c *CozeOAuthConfig) (*coze.JWTOAuthClient, *coze.OAuthConfig, error) {
	configFile, err := json.Marshal(c)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal config: %v", err)
	}

	var oauthConfig coze.OAuthConfig
	if err := json.Unmarshal(configFile, &oauthConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	oauth, err := coze.LoadOAuthAppFromConfig(&oauthConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load OAuth config: %v", err)
	}

	jwtClient, ok := oauth.(*coze.JWTOAuthClient)
	if !ok {
		return nil, nil, fmt.Errorf("invalid OAuth client type: expected JWT client")
	}
	return jwtClient, &oauthConfig, nil
}

type Client struct {
	JWTOAuthClient *coze.JWTOAuthClient
	OAuthConfig    *coze.OAuthConfig // 配置(有coze_www_base和coze_api_base)
	OAuthToken     *coze.OAuthToken  // 令牌(有access_token, refresh_token, expires_in, token_type)
}

func NewClient(c *CozeOAuthConfig) (*Client, error) {
	oauth, oauthConfig, err := loadConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth config: %v", err)
	}

	ctx := context.Background()

	resp, err := oauth.GetAccessToken(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}

	return &Client{
		JWTOAuthClient: oauth,
		OAuthConfig:    oauthConfig,
		OAuthToken:     resp,
	}, nil
}

func (c Client) post(api string, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

	api = c.OAuthConfig.CozeAPIBase + api

	request, err := http.NewRequest(http.MethodPost, api, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.OAuthToken.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d, 响应: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil

}

func (c Client) get(api string) ([]byte, error) {
	api = c.OAuthConfig.CozeAPIBase + api

	request, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.OAuthToken.AccessToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d, 响应: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
