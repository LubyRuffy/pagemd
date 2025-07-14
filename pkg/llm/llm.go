package llm

import (
	"bytes"
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type ModelConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
	APIKey  string `yaml:"api_key" json:"api_key"`
	Model   string `yaml:"model" json:"model"`
}

type Config struct {
	Model ModelConfig `yaml:"model" json:"model"`
}

func New(baseUrl string, model string, key string) *ModelConfig {
	return &ModelConfig{
		BaseURL: baseUrl,
		APIKey:  key,
		Model:   model,
	}
}

// Load 加载配置
func Load(cfgFile string) (*ModelConfig, error) {
	var config Config
	// Read configuration file
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	yamlFile, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))

	decoder := yaml.NewDecoder(bytes.NewReader(yamlFile))
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config.Model, nil
}

func (a *ModelConfig) MustGetModel(ctx context.Context) model.ToolCallingChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: a.BaseURL,
		APIKey:  a.APIKey,
		Model:   a.Model,
	})
	if err != nil {
		log.Fatalf("failed to create chat model: %v", err)
	}
	return chatModel
}

func (a *ModelConfig) Stream(ctx context.Context, systemPrompt, userPrompt string, onData func(string)) (string, error) {
	result := ""
	s, err := a.MustGetModel(ctx).Stream(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: systemPrompt,
		},
		{
			Role:    schema.User,
			Content: userPrompt,
		},
	})
	if err != nil {
		return "", err
	}
	for {
		v, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}
		result += v.Content
		if onData != nil {
			onData(v.Content)
		}
	}
	return result, nil
}
