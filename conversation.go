package coze

import (
	"encoding/json"
	"net/url"
)

type CreateConversationEnterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreateConversationBody struct {
	// 会话中的消息内容。详细说明可参考 EnterMessage object
	Messages []CreateConversationEnterMessage `json:"messages,omitempty"`

	/*
		创建会话时的附加消息，查看会话时也会返回此附加消息。

		自定义键值对，应指定为 Map 对象格式。长度为 16 对键值对，其中键（key）的长度范围为 1～64 个字符，值（value）的长度范围为 1～512 个字符
	*/
	MetaData map[string]any `json:"meta_data,omitempty"`
}

type CreateConversationResp struct {
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

// CreateConversation 创建会话
func (c Client) CreateConversation(body *CreateConversationBody) (*CreateConversationResp, error) {
	api := "/v1/conversation/create"

	b, err := c.post(api, body, nil)
	if err != nil {
		return nil, err
	}

	var resp CreateConversationResp
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type RetrieveConversationQuery struct {
	ConversationId string `json:"conversation_id"` // 会话ID
}

func (r RetrieveConversationQuery) Encode() string {
	query := url.Values{}
	query.Add("conversation_id", r.ConversationId)
	return query.Encode()
}

type RetrieveConversationResp struct {
	Code int `json:"code"`
	Data struct {
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
}

// RetrieveConversation 查询会话信息
func (c Client) RetrieveConversation(query *RetrieveConversationQuery) (*RetrieveConversationResp, error) {
	api := "/v1/conversation/retrieve"

	b, err := c.get(api, query)
	if err != nil {
		return nil, err
	}

	var resp RetrieveConversationResp
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
