package aitrans

import (
	"context"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
)

var (
	OllamaEndpoint = "http://localhost:11434"
	OllamaModel    = "qwen2.5:7b"
	//OllamaModel         = "qwen2.5:32b"
	DefaultSystemPrompt = `You are a highly skilled translator tasked with translating various types of content from other languages into Chinese. Follow these instructions carefully to complete the translation task:

## Glossary

Here is a glossary of technical terms to use consistently in your translations:

- AGI -> 通用人工智能
- LLM/Large Language Model -> 大语言模型
- Transformer -> Transformer
- Token -> Token
- Generative AI -> 生成式 AI
- AI Agent -> AI 智能体
- prompt -> 提示词
- zero-shot -> 零样本学习
- few-shot -> 少样本学习
- multi-modal -> 多模态
- fine-tuning -> 微调

## INPUTS

- 输入格式为Markdown格式。

## OUTPUTS

- 输出为Markdown格式
- 输出为中文
- 输出格式要保留跟输入一样的格式
- 对于代码块（用` + "```lang ```" + `符号包含）中的内容不要翻译 
- 保证内容的完整和准确性
- 除了翻译的内容，不要增加任何其他的说明和注释
- 如果输入本身的语言就是中文，你就直接输出“无需翻译”
`
)

type AiTranslator struct {
	ollamaClient *api.Client
}

var (
	True = true
)

func (a *AiTranslator) TranslateToChinese(ctx context.Context, md string, onData func(string)) (string, error) {
	result := ""
	err := a.ollamaClient.Chat(ctx, &api.ChatRequest{
		Model: OllamaModel,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: DefaultSystemPrompt,
			},
			{
				Role:    "user",
				Content: md,
			},
		},
		Stream: &True,
		Options: map[string]interface{}{
			"num_ctx": 102400,
		},
	}, func(resp api.ChatResponse) error {
		result += resp.Message.Content
		if onData != nil {
			onData(resp.Message.Content)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func New() *AiTranslator {
	u, _ := url.Parse(OllamaEndpoint)
	client := api.NewClient(u, http.DefaultClient)
	ait := &AiTranslator{
		ollamaClient: client,
	}
	return ait
}
