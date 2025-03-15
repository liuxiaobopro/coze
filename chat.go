package coze

// 文档链接: https://www.coze.cn/open/docs/developer_guides/chat_v3

import "encoding/json"

type ChatQuery struct {
	ConversationID string `json:"conversation_id"`
}

type ChatBody struct {
}

type ChatResp struct {
}

func (c Client) Chat(query *ChatQuery, body *ChatBody) (*ChatResp, error) {
	api := "/v3/chat"

	b, err := c.post(api, body, query)
	if err != nil {
		return nil, err
	}

	var resp ChatResp
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
