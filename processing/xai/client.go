package xai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.x.ai/v1"
	maxToolRounds  = 16
)

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

type ChunkCallback func(chunk string) error
type ToolHandler func(ctx context.Context, arguments json.RawMessage) (any, error)

type Tool struct {
	Type        string         `json:"type"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Handler     ToolHandler    `json:"-"`
}

type ToolCall struct {
	ID       string           `json:"id,omitempty"`
	Type     string           `json:"type,omitempty"`
	Function ToolCallFunction `json:"function,omitempty"`
}

type ToolCallFunction struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type Client struct {
	model      string
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type chatCompletionRequest struct {
	Model             string     `json:"model"`
	Messages          []Message  `json:"messages"`
	Stream            bool       `json:"stream"`
	Tools             []chatTool `json:"tools,omitempty"`
	ParallelToolCalls *bool      `json:"parallel_tool_calls,omitempty"`
}

type chatTool struct {
	Type     string           `json:"type"`
	Function chatToolFunction `json:"function"`
}

type chatToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

type chatCompletionChunk struct {
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Choices []struct {
		Delta struct {
			Role      string          `json:"role"`
			Content   string          `json:"content"`
			ToolCalls []toolCallDelta `json:"tool_calls"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type toolCallDelta struct {
	Index    int    `json:"index"`
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

func New(model, apiKey string) *Client {
	return &Client{
		model:   model,
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: time.Hour,
		},
	}
}

func Stream(ctx context.Context, model, apiKey string, messages []Message, tools []Tool, callback ChunkCallback) error {
	return New(model, apiKey).StreamChatCompletions(ctx, messages, tools, callback)
}

func (c *Client) StreamChatCompletions(ctx context.Context, messages []Message, tools []Tool, callback ChunkCallback) error {
	if c == nil {
		return errors.New("xai client is nil")
	}
	if c.model == "" {
		return errors.New("xai model is required")
	}
	if c.apiKey == "" {
		return errors.New("xai api key is required")
	}
	if len(messages) == 0 {
		return errors.New("at least one message is required")
	}
	if callback == nil {
		return errors.New("chunk callback is required")
	}

	toolMap, err := toolMapFor(tools)
	if err != nil {
		return err
	}

	conversation := append([]Message(nil), messages...)
	for round := 0; round < maxToolRounds; round++ {
		assistantMessage, err := c.streamChatCompletionOnce(ctx, conversation, tools, callback)
		if err != nil {
			return err
		}
		if len(assistantMessage.ToolCalls) == 0 {
			return nil
		}

		conversation = append(conversation, assistantMessage)

		toolResults, err := executeToolCalls(ctx, assistantMessage.ToolCalls, toolMap)
		if err != nil {
			return err
		}
		conversation = append(conversation, toolResults...)
	}

	return fmt.Errorf("xai tool loop exceeded %d rounds", maxToolRounds)
}

func (c *Client) streamChatCompletionOnce(ctx context.Context, messages []Message, tools []Tool, callback ChunkCallback) (Message, error) {
	reqBody := chatCompletionRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   true,
	}

	if len(tools) > 0 {
		reqBody.Tools = make([]chatTool, 0, len(tools))
		for _, tool := range tools {
			reqBody.Tools = append(reqBody.Tools, chatTool{
				Type: "function",
				Function: chatToolFunction{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  tool.Parameters,
				},
			})
		}
		reqBody.ParallelToolCalls = boolPtr(false)
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return Message{}, fmt.Errorf("marshal xai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return Message{}, fmt.Errorf("build xai request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("send xai request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return Message{}, fmt.Errorf("xai request failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return consumeStream(resp.Body, callback)
}

func consumeStream(r io.Reader, callback ChunkCallback) (Message, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	message := Message{Role: "assistant"}
	toolCallsByIndex := map[int]*ToolCall{}
	toolCallOrder := make([]int, 0, 4)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") || !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			return finalizeStreamMessage(message, toolCallsByIndex, toolCallOrder), nil
		}

		var chunk chatCompletionChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			return Message{}, fmt.Errorf("decode xai chunk: %w", err)
		}

		if chunk.Error != nil {
			return Message{}, errors.New(chunk.Error.Message)
		}

		for _, choice := range chunk.Choices {
			if choice.Delta.Role != "" {
				message.Role = choice.Delta.Role
			}
			for _, delta := range choice.Delta.ToolCalls {
				toolCall, ok := toolCallsByIndex[delta.Index]
				if !ok {
					toolCall = &ToolCall{}
					toolCallsByIndex[delta.Index] = toolCall
					toolCallOrder = append(toolCallOrder, delta.Index)
				}
				if delta.ID != "" {
					toolCall.ID = delta.ID
				}
				if delta.Type != "" {
					toolCall.Type = delta.Type
				}
				if delta.Function.Name != "" {
					toolCall.Function.Name = delta.Function.Name
				}
				if delta.Function.Arguments != "" {
					toolCall.Function.Arguments += delta.Function.Arguments
				}
			}
			if choice.Delta.Content != "" {
				if err := callback(choice.Delta.Content); err != nil {
					return Message{}, err
				}
				message.Content += choice.Delta.Content
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return Message{}, fmt.Errorf("read xai stream: %w", err)
	}

	return finalizeStreamMessage(message, toolCallsByIndex, toolCallOrder), nil
}

func finalizeStreamMessage(message Message, toolCallsByIndex map[int]*ToolCall, toolCallOrder []int) Message {
	if len(toolCallsByIndex) == 0 {
		return message
	}

	message.ToolCalls = make([]ToolCall, 0, len(toolCallsByIndex))
	for _, index := range toolCallOrder {
		if toolCall := toolCallsByIndex[index]; toolCall != nil {
			message.ToolCalls = append(message.ToolCalls, *toolCall)
		}
	}

	return message
}

func toolMapFor(tools []Tool) (map[string]Tool, error) {
	if len(tools) == 0 {
		return nil, nil
	}

	toolMap := make(map[string]Tool, len(tools))
	for _, tool := range tools {
		if tool.Type == "" {
			tool.Type = "function"
		}
		if tool.Type != "function" {
			return nil, fmt.Errorf("unsupported xai tool type %q", tool.Type)
		}
		if tool.Name == "" {
			return nil, errors.New("xai tool name is required")
		}
		if tool.Handler == nil {
			return nil, fmt.Errorf("xai tool %q handler is required", tool.Name)
		}
		if _, exists := toolMap[tool.Name]; exists {
			return nil, fmt.Errorf("duplicated xai tool %q", tool.Name)
		}
		toolMap[tool.Name] = tool
	}

	return toolMap, nil
}

func executeToolCalls(ctx context.Context, toolCalls []ToolCall, toolMap map[string]Tool) ([]Message, error) {
	if len(toolCalls) == 0 {
		return nil, nil
	}

	results := make([]Message, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		tool, ok := toolMap[toolCall.Function.Name]
		var output any
		if !ok {
			output = map[string]any{"error": fmt.Sprintf("unknown tool %q", toolCall.Function.Name)}
		} else {
			result, err := tool.Handler(ctx, json.RawMessage(toolCall.Function.Arguments))
			if err != nil {
				output = map[string]any{"error": err.Error()}
			} else {
				output = result
			}
		}

		payload, err := json.Marshal(output)
		if err != nil {
			return nil, fmt.Errorf("marshal tool result for %q: %w", toolCall.Function.Name, err)
		}

		results = append(results, Message{
			Role:       "tool",
			ToolCallID: toolCall.ID,
			Content:    string(payload),
		})
	}

	return results, nil
}

func boolPtr(v bool) *bool {
	return &v
}
