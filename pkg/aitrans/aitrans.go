package aitrans

import (
	"context"
	"fmt"
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
`
)

type AiTranslator struct {
	//ollamaClient *openai.Client
	ollamaClient *api.Client
}

var (
	False = false
	True  = true
)

func (a *AiTranslator) TranslateToChinese(ctx context.Context, md string, onData func(string)) (string, error) {
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
		fmt.Print(resp.Message.Content)
		return nil
	})
	if err != nil {
		return "", err
	}
	return "", nil
}

//func (a *AiTranslator) TranslateToChinese(ctx context.Context, md string, onData func(string)) (string, error) {
//	stream, err := a.ollamaClient.CreateChatCompletionStream(ctx,
//		openai.ChatCompletionRequest{
//			Model: OllamaModel, // Use the model parameter
//			Messages: []openai.ChatCompletionMessage{
//				{Role: openai.ChatMessageRoleSystem, Content: DefaultSystemPrompt},
//				{Role: openai.ChatMessageRoleUser, Content: md},
//			},
//			MaxTokens:           102400,
//			MaxCompletionTokens: 102400,
//			//ResponseFormat: &openai.ChatCompletionResponseFormat{
//			//	Type: openai.ChatCompletionResponseFormatTypeJSONObject,
//			//},
//		},
//	)
//	if err != nil {
//		return "", err
//	}
//	defer stream.Close()
//
//	var output string
//_out:
//	for {
//		response, err := stream.Recv()
//		if errors.Is(err, io.EOF) {
//			//log.Println("stream closed")
//			break
//		}
//		if err != nil {
//			return "", err
//		}
//		//log.Println(response)
//		for _, choice := range response.Choices {
//			//fmt.Printf("%s", choice.Delta.Content)
//			output += choice.Delta.Content
//			onData(choice.Delta.Content)
//
//			if choice.FinishReason != "" {
//				//log.Println(choice.FinishReason)
//				break _out
//			}
//		}
//
//	}
//
//	return output, nil
//}

func New() *AiTranslator {
	// Initialize the OpenAI client with your API key
	//openaiConfig := openai.DefaultConfig("ollama-token")
	//openaiConfig.BaseURL = fmt.Sprintf("%s/v1", OllamaEndpoint)
	//client := openai.NewClientWithConfig(openaiConfig)

	u, _ := url.Parse(OllamaEndpoint)
	client := api.NewClient(u, http.DefaultClient)
	ait := &AiTranslator{
		//ollamaClient: client,
		ollamaClient: client,
	}
	return ait
}
