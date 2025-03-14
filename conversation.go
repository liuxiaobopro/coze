package coze

import "encoding/json"

type ConversationEnterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ConversationReq struct {
	// 会话中的消息内容。详细说明可参考 EnterMessage object
	Messages []ConversationEnterMessage `json:"messages,omitempty"`

	/*
		创建会话时的附加消息，查看会话时也会返回此附加消息。

		自定义键值对，应指定为 Map 对象格式。长度为 16 对键值对，其中键（key）的长度范围为 1～64 个字符，值（value）的长度范围为 1～512 个字符
	*/
	MetaData map[string]any `json:"meta_data,omitempty"`
}

type ConversationResp struct {
	Code int `json:"code"`
	Data struct {
		CreatedAt     int64  `json:"created_at"`
		ID            string `json:"id"`
		LastSectionID string `json:"last_section_id"`
		MetaData      any    `json:"meta_data"`
	} `json:"data"`
	Detail struct {
		Logid string `json:"logid"`
	} `json:"detail"`
	Msg string `json:"msg"`
}

func (c Client) CreateConversation(body *ConversationReq) (*ConversationResp, error) {
	api := "/v1/conversation/create"

	b, err := c.post(api, body)
	if err != nil {
		return nil, err
	}

	var resp ConversationResp
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
