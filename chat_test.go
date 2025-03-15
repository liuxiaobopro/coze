package coze

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClient_Chat_Stream(t *testing.T) {
	st := time.Now()
	defer func() {
		fmt.Printf("耗时: %v\n", time.Since(st))
	}()

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

	var chatList []string
	var doneChan = make(chan bool, 1)

	go func() {
		_, _ = cc.Chat(&ChatBody{
			BotId:  "7480830640758308874",
			UserId: "123",
			AdditionalMessages: []ChatBodyAdditionalMessages{
				{
					Role:        "user",
					Type:        "question",
					Content:     "https://juejin.cn/post/7468323178931879972",
					ContentType: "text",
				},
			},
			Stream: true,
		}, nil, func(line string, err error) {
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				log.Fatalf("Error chatting: %v", err)
			}

			chatList = append(chatList, strings.TrimSpace(line))

			fmt.Println(line)

			if strings.Contains(line, "[DONE]") {
				doneChan <- true
			}
		})
	}()

	<-doneChan

	fmt.Println("--------------------------------")

	type resp struct {
		Type string `json:"type"`
	}

	for k, v := range chatList {
		if strings.Contains(v, "conversation.message.completed") {
			jsonStr := strings.TrimPrefix(chatList[k+1], "data:")

			var r resp
			if err := json.Unmarshal([]byte(jsonStr), &r); err != nil {
				log.Fatalf("Error unmarshalling: %v", err)
			}

			if r.Type == "answer" {
				fmt.Printf("第%d条消息[最终答案]:\n", k+1)
				fmt.Println(jsonStr)
			}
		}
	}
}

func TestClient_Chat_NO_Stream(t *testing.T) {
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

	resp, err := cc.Chat(&ChatBody{
		BotId:  "7480830640758308874",
		UserId: "123",
		AdditionalMessages: []ChatBodyAdditionalMessages{
			{
				Role:        "user",
				Type:        "question",
				Content:     "https://juejin.cn/post/7468323178931879972",
				ContentType: "text",
			},
		},
	}, nil, nil)
	if err != nil {
		log.Fatalf("Error chatting: %v", err)
	}

	fmt.Println(resp)
}
