package coze

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 文档链接: https://www.coze.cn/open/docs/developer_guides/chat_v3

type ChatQuery struct {
	ConversationID string `json:"conversation_id"`
}

func (q ChatQuery) Encode() string {
	return url.Values{"conversation_id": {q.ConversationID}}.Encode()
}

type ChatBodyAdditionalMessages struct {
	/*
		发送这条消息的实体。取值：

		user：代表该条消息内容是用户发送的。
		assistant：代表该条消息内容是智能体发送的。
	*/
	Role string `json:"role"`

	/*
		消息类型。默认为 question。

		question：用户输入内容。

		answer：智能体返回给用户的消息内容，支持增量返回。如果工作流绑定了消息节点，可能会存在多 answer 场景，此时可以用流式返回的结束标志来判断所有 answer 完成。

		function_call：智能体对话过程中调用函数（function call）的中间结果。

		tool_response：调用工具 （function call）后返回的结果。

		follow_up：如果在 智能体上配置打开了用户问题建议开关，则会返回推荐问题相关的回复内容。不支持在请求中作为入参。

		verbose：多 answer 场景下，服务端会返回一个 verbose 包，对应的 content 为 JSON 格式，content.msg_type =generate_answer_finish 代表全部 answer 回复完成。不支持在请求中作为入参。
	*/
	Type string `json:"type"`

	/*
		消息的内容，支持纯文本、多模态（文本、图片、文件混合输入）、卡片等多种类型的内容。

		content_type 为 object_string 时，content 为 object_string object 数组序列化之后的 JSON String，详细说明可参考 object_string object。

		当 content_type = text 时，content 为普通文本，例如 "content" :"Hello!"。
	*/
	Content string `json:"content"`

	/*
	   消息内容的类型，支持设置为：

	   text：文本。

	   object_string：多模态内容，即文本和文件的组合、文本和图片的组合。

	   card：卡片。此枚举值仅在接口响应中出现，不支持作为入参
	*/
	ContentType string `json:"content_type"`
}

type ChatBody struct {
	// 机器人ID(必填)
	BotId string `json:"bot_id"`
	// 用户ID(必填)
	UserId string `json:"user_id"`
	// 额外消息
	AdditionalMessages []ChatBodyAdditionalMessages `json:"additional_messages"`
	// 是否流式
	Stream bool `json:"stream"`
}

type ChatResp struct {
	Data struct {
		ID             string `json:"id"`
		ConversationID string `json:"conversation_id"`
		BotID          string `json:"bot_id"`
		CreatedAt      int    `json:"created_at"`
		LastErr        struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		} `json:"last_error"`
		Status string `json:"status"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c ChatResp) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c Client) Chat(body *ChatBody, query *ChatQuery, callback func(line string, err error)) (*ChatResp, error) {
	api := c.OAuthConfig.CozeAPIBase + "/v3/chat"
	if query != nil {
		api += "?" + query.Encode()
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

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

	if body.Stream {
		reader := bufio.NewReader(response.Body)
		for {
			callback(reader.ReadString('\n'))
		}
	} else {
		var resp ChatResp
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, fmt.Errorf("解析失败: %v", err)
		}

		return &resp, nil
	}
}
