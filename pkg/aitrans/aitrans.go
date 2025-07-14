package aitrans

import (
	"context"
	"github.com/LubyRuffy/pagemd/pkg/llm"
)

var (
	DefaultSystemPrompt = `\
You are a highly skilled translator tasked with translating various types of content from other languages into Chinese. Follow these instructions carefully to complete the translation task:

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
	cfg *llm.ModelConfig
}

func (a *AiTranslator) TranslateToChinese(ctx context.Context, md string, onData func(string)) (string, error) {
	return a.cfg.Stream(ctx, DefaultSystemPrompt, md, onData)
}

func New(config *llm.ModelConfig) *AiTranslator {
	ait := &AiTranslator{
		cfg: config,
	}
	return ait
}
