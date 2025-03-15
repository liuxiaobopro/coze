package coze

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestClient_CreateConversation(t *testing.T) {
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

	resp, err := cc.CreateConversation(&CreateConversationBody{
		Messages: []CreateConversationEnterMessage{
			{
				Role:    "user",
				Content: "你好",
			},
		},
	})
	if err != nil {
		log.Fatalf("Error creating conversation: %v", err)
	}

	fmt.Println(resp.Data.ID)

}

func TestClient_RetrieveConversation(t *testing.T) {
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

	resp, err := cc.RetrieveConversation(&RetrieveConversationQuery{
		ConversationId: "7481851577591054371",
	})
	if err != nil {
		log.Fatalf("Error retrieving conversation: %v", err)
	}
	fmt.Println(resp)
}
