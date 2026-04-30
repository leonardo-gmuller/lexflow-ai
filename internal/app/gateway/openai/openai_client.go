package open_ai

import (
	"errors"
	"time"

	"github.com/leonardo-gmuller/lexflow-ai/internal/app/config"
	"github.com/leonardo-gmuller/lexflow-ai/internal/app/domain/port"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Ajuste seu wrapper para v3 do SDK
type OpenAIClient struct {
	cli        *openai.Client
	model      string
	temp       float64
	topP       float64
	timeout    time.Duration
	apiKey     string
	toolEngine port.ToolEngine
}

func NewOpenAI(cfg config.OpenAIConfig) (*OpenAIClient, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("OPENAI_API_KEY vazio")
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-4.1-mini"
	}

	aai := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)

	return &OpenAIClient{
		cli:     &aai,
		model:   cfg.Model,
		temp:    float64(cfg.Temperature) / 100.0,
		topP:    float64(cfg.TopP) / 100.0,
		timeout: time.Duration(cfg.Timeout) * time.Second,
		apiKey:  cfg.APIKey,
	}, nil
}
