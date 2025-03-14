package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	fileBytes, err := os.ReadFile("coze_oauth_config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var oc CozeOAuthConfig
	if err := json.Unmarshal(fileBytes, &oc); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	oauth, oauthConfig, err := oc.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx := context.Background()

	resp, err := oauth.GetAccessToken(ctx, nil)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	api := oauthConfig.CozeAPIBase + "/v1/workspaces"

	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+resp.AccessToken)

	client := http.DefaultClient

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Println(string(body))

}
