package main

import (
	"encoding/json"
	"fmt"

	"github.com/coze-dev/coze-go"
)

type CozeOAuthConfig struct {
	ClientType  string `json:"client_type"`
	ClientID    string `json:"client_id"`
	CozeWWWBase string `json:"coze_www_base"`
	CozeAPIBase string `json:"coze_api_base"`
	PrivateKey  string `json:"private_key"`
	PublicKeyID string `json:"public_key_id"`
}

func (c CozeOAuthConfig) LoadConfig() (*coze.JWTOAuthClient, *coze.OAuthConfig, error) {
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
