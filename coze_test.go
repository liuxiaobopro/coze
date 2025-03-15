package coze

import (
	"encoding/json"
	"fmt"
	"log"
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

	cc, err := NewClient(&oc)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	resp, err := cc.get("/v1/workspaces", nil)
	if err != nil {
		log.Fatalf("Error getting workspaces: %v", err)
	}

	fmt.Println(string(resp))
}
