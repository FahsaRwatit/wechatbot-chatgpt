package sse

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/launchdarkly/eventsource"
)

/*
在请求头和响应头设置text/event-stream是实现SSE的关键。
SSE （ Server-sent Events ）是 WebSocket 的一种轻量代替方案，使用 HTTP 协议。
*/
type Client struct {
	URL          string
	EventChannel chan string
	Headers      map[string]string
}

func Init(url string) Client {
	return Client{
		URL:          url,
		EventChannel: make(chan string),
	}
}

// 请求 https://chat.openai.com/backend-api/conversation
func (c *Client) Connect(message string, conversationId string, parentMessageId string) error {
	messages, err := json.Marshal([]string{message})
	if err != nil {
		return errors.New(fmt.Sprintf("failed to encode message: %v", err))
	}
	if parentMessageId == "" {
		parentMessageId = uuid.NewString()
	}
	var conversationIdString string
	if conversationId != "" {
		conversationIdString = fmt.Sprintf(`, "conversation_id": "%s"`, conversationId)
	}
	// if conversation id is empty, don't send it
	body := fmt.Sprintf(`{
        "action": "next",
        "messages": [
            {
                "id": "%s",
                "role": "user",
                "content": {
                    "content_type": "text",
                    "parts": %s
                }
            }
        ],
        "model": "text-davinci-002-render",
		"parent_message_id": "%s"%s
    }`, uuid.NewString(), string(messages), parentMessageId, conversationIdString)
	req, err := http.NewRequest("POST", c.URL, strings.NewReader(body))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create request: %v", err))
	}
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	// 在请求头和响应头设置text/event-stream是实现SSE的关键。
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	http := &http.Client{}
	resp, err := http.Do(req)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to connect to SSE: %v", err))
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("failed to connect to SSE: %v", resp.Status))
	}
	// eventsource是一种单向通信的方式,只能由服务端向客户端推送消息,可以自动重连接,可以发送随机事件
	decoder := eventsource.NewDecoder(resp.Body)
	go func() {
		defer resp.Body.Close()
		defer close(c.EventChannel)

		for {
			event, err := decoder.Decode()
			if err != nil {
				log.Println(errors.New(fmt.Sprintf("failed to decode event: %v", err)))
				break
			}
			if event.Data() == "[DONE]" || event.Data() == "" {
				break
			}

			c.EventChannel <- event.Data()
		}
	}()

	return nil
}
