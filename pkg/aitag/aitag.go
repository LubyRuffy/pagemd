package aitag

import (
	"context"
	"encoding/json"
	"github.com/ollama/ollama/api"
	"net/http"
	"net/url"
)

var (
	OllamaEndpoint = "http://localhost:11434"
	//OllamaModel    = "command-r7b"
	//OllamaModel    = "phi4"
	//OllamaModel = "qwen2.5:7b"
	//OllamaModel = "qwen2.5:14b"
	OllamaModel         = "qwen2.5:32b"
	DefaultSystemPrompt = `
You are an expert IT engineer specializing in cybersecurity and vulnerability analysis with a meticulous approach to detail. Your task is to meticulously extract classifications from the provided text content, strictly adhering to the predefined classification structure below.

## Classification Structure

Classifications must be an *exact* match to the classification names in the following structure (case-sensitive). The part before the semicolon is the classification, and the part after is the description.

- technical
	- cybersecurity
		- pentest: Penetration testing
		- vulnerability exploitation: Vulnerability exploit 漏洞利用专题文章，涉及到富有启发意义的、具有创新性的漏洞利用技术。如果有对应的编号，需要提取编号信息，如"CVE-2024-11320"，"ms17-010"等
		- security tools: 安全产品或安全工具的发布（更新）或者说明（包括用法和示例）等。包括但不限于如下列表。
			- IDA
			- burp
			- metasploit
			- nmap
			- shodan
			- censys
			- fofa
		- APT: APT(Advanced persistent threat)组织和事件的追踪，溯源，反制，拓线。如果有对应的编号和组织信息，需要进行提取，比如"APT-C-39"
		- blockchain security: 区块链安全相关的文章，包括但不限于合约安全、协议安全、节点/客户端安全、交易所安全、钓鱼欺诈等。
		- supply-chain security: 由于第三方的商业或开源软件/系统组件存在安全漏洞而导致的安全事件或报告。
		- browser security: 主流网页浏览器相关的安全文章。
		- kernel security: Linux、Windows、MacOS、iOS、Android等操作系统内核相关的安全技术文章。
		- CTF: CTF(Capture the flag)解题报告、比赛通知和赛况等。
		- RCE of Core Components: 可以达成远程代码执行的核心组件的漏洞。核心组件是指在软件栈中具有基础地位库、中间件、服务，如Apache Dubbo、Log4j等。
		- network infrastructure vulnerabilities: 网络基础设施组件（如CDN）以及基础网络协议（如TLS，TCP）相关的漏洞
		- cryptography: 密码学相关的问题，如因密码误用导致的漏洞、密码算法本身的漏洞等。
		- fuzzing: 与模糊测试相关的研究
	- programming: 开发，包括但不限于如下列表。
		- golang
		- python
		- rust
		- javascript
		- assembly
		- software engineering: 讨论代码质量提升，代码重构（refactor），软件架构，系统建设方案，最佳实践之类的。
	- AI
		- LLM: 大模型，包括但不限于如下列表。
			- qwen
			- llama
			- deepseek
			- gemini
			- openai
			- chatgpt
			- claude
		- AI tools: AI工具，包括但不限于如下列表。
			- ollama
			- llama.cpp
			- cline
	- news
		- data breach: 个人、企业或政府机构的敏感信息泄露事件。
		- hacking incidents
		- bug bounty: 由软件厂商公布的漏洞赏金计划，或安全研究人员撰写的赏金获取经验。
		- ransomware: 与勒索软件行动相关的新闻、勒索软件样本分析等
		- standards and laws & regulations: 新的网络安全相关的标准和法律法规的发布。
- non-technical

## Input

Input is in Markdown format.

## Output Requirements

- Extract classifications based on the provided classification structure.
- Ensure each extracted classification includes all its parent classifications.
    - For example, if "fofa" is identified, the output must include: "fofa", "security tools", "cybersecurity", "technical".
    - For example, if "python" is identified, the output must include: "python", "programming", "technical".
- Multiple tags can be selected within a major category. For example, both "ollama" and "python" can be selected if relevant, with their respective parent nodes included: "ollama", "AI tools", "AI", "technical", "python", "programming", "technical".
- Only one of the major classifications, "technical" or "non-technical", can be selected. They are mutually exclusive.
- The output must be a JSON formatted string with a key named "tags" whose value is a one-dimensional string array. For example: {"tags": ["fofa", "security tools", "cybersecurity", "technical"]}
- Extract all relevant classifications exhaustively and accurately. For instance, if "fuzzing" technology is used in the process of "vulnerability exploitation", both classifications are required.
- All classifications must be output exactly as they appear in the classification structure, maintaining the original letter case. No translation, additional explanations, or modifications to the casing are permitted.
- Ensure all classifications in the output array are unique.
`
)

type AiTagger struct {
	ollamaClient *api.Client
}

var (
	True = true
)

type tag struct {
	Tags []string `json:"tags"`
}

func (a *AiTagger) Tag(ctx context.Context, md string, onData func(string)) ([]string, error) {
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
		Format: json.RawMessage(`"json"`),
		Options: map[string]interface{}{
			"num_ctx":     102400,
			"temperature": 0,
		},
	}, func(resp api.ChatResponse) error {
		result += resp.Message.Content
		if onData != nil {
			onData(resp.Message.Content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var t tag
	if err = json.Unmarshal([]byte(result), &t); err != nil {
		return nil, err
	}

	return t.Tags, nil
}

func New() *AiTagger {
	u, _ := url.Parse(OllamaEndpoint)
	client := api.NewClient(u, http.DefaultClient)
	ait := &AiTagger{
		ollamaClient: client,
	}
	return ait
}
