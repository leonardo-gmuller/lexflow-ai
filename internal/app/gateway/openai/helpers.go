package open_ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/port"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

var temperatureModels = map[string]bool{
	"gpt-3.5-turbo":        true,
	"gpt-4":                true,
	"gpt-4.1-mini":         true,
	"gpt-4-vision-preview": true,
}

func modelSupportsTemperature(m string) bool {
	return temperatureModels[m]
}

type functionCall struct {
	CallID    string
	Name      string
	Arguments string
}

func (c *OpenAIClient) extractFunctionCalls(resp *responses.Response) []functionCall {
	var calls []functionCall

	for _, item := range resp.Output {
		if item.Type != "function_call" {
			continue
		}

		calls = append(calls, functionCall{
			CallID:    item.CallID,
			Name:      item.Name,
			Arguments: item.Arguments.OfString,
		})
	}

	return calls
}

func (c *OpenAIClient) executeToolCalls(calls []functionCall) ([]responses.ResponseInputItemUnionParam, error) {
	out := make([]responses.ResponseInputItemUnionParam, 0, len(calls))

	for _, call := range calls {
		args := map[string]any{}
		if strings.TrimSpace(call.Arguments) != "" {
			if err := json.Unmarshal([]byte(call.Arguments), &args); err != nil {
				return nil, fmt.Errorf("erro parse args tool %s: %w", call.Name, err)
			}
		}

		result, err := c.toolEngine.Execute(call.Name, args)
		if err != nil {
			result = fmt.Sprintf("tool error: %v", err)
		}

		out = append(out, responses.ResponseInputItemUnionParam{
			OfCustomToolCallOutput: &responses.ResponseCustomToolCallOutputParam{
				CallID: call.CallID,
				Output: responses.ResponseCustomToolCallOutputOutputUnionParam{OfString: openai.String(result)},
			},
		})
	}

	return out, nil
}

func (c *OpenAIClient) responsesNewWithRetry(ctx context.Context, params responses.ResponseNewParams) (*responses.Response, error) {
	const maxRetries = 3

	for attempt := 0; attempt < maxRetries; attempt++ {

		resp, err := c.cli.Responses.New(ctx, params)
		if err == nil {
			return resp, nil
		}

		// erro de conversation_locked -> backoff e tenta de novo
		backoff := time.Second * time.Duration(1<<attempt)
		log.Printf("[openai] conversation_locked, retrying in %s (attempt %d/%d)", backoff, attempt+1, maxRetries)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
			// tenta de novo
		}
	}

	return nil, fmt.Errorf("conversation_locked após %d tentativas", maxRetries)
}

func (c *OpenAIClient) SetToolEngine(toolEngine port.ToolEngine) {
	c.toolEngine = toolEngine
}

func (c *OpenAIClient) buildResponseContents(
	ctx context.Context,
	userMsgs []string,
) (responses.ResponseInputMessageContentListParam, error) {
	var contents []responses.ResponseInputContentUnionParam
	for _, userMsg := range userMsgs {
		userMsg = strings.TrimSpace(userMsg)
		if userMsg == "" {
			continue
		}

		contents = append(contents,
			responses.ResponseInputContentUnionParam{
				OfInputText: &responses.ResponseInputTextParam{
					Text: userMsg,
					Type: "input_text",
				},
			},
		)
	}

	return contents, nil
}

func (c *OpenAIClient) buildTools() []responses.ToolUnionParam {
	defs := c.toolEngine.GetTools()
	out := make([]responses.ToolUnionParam, 0, len(defs))

	for _, d := range defs {
		name, _ := d["name"].(string)
		description, _ := d["description"].(string)
		parameters, _ := d["parameters"].(map[string]any)

		out = append(out, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        name,
				Description: openai.String(description),
				Parameters:  parameters,
				Type:        "function",
			},
		})
	}

	return out
}
