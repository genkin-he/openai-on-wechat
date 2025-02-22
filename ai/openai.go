package ai

import (
	"context"
	"log"
	"strings"

	"github.com/riba2534/openai-on-wechat/config"
	"github.com/sashabaranov/go-openai"
)

var SystemMessage openai.ChatCompletionMessage

var textOpenAIClient *openai.Client
var imageOpenAIClient *openai.Client

func Init() {
	// text init
	textConfig := openai.DefaultConfig(config.C.WechatConfig.TextConfig.AuthToken)
	textConfig.BaseURL = config.C.WechatConfig.TextConfig.OpenApiUrl // 使用反向代理的地址
	textOpenAIClient = openai.NewClientWithConfig(textConfig)
	// image init
	imageConfig := openai.DefaultConfig(config.C.WechatConfig.ImageConfig.AuthToken)
	imageConfig.BaseURL = config.C.WechatConfig.ImageConfig.OpenApiUrl // 使用反向代理的地址
	imageOpenAIClient = openai.NewClientWithConfig(imageConfig)
	// Prompt init
	SystemMessage = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: config.Prompt,
	}
}

func GetOpenAITextReply(q string) string {
	resp, err := textOpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				SystemMessage,
				{
					Role:    openai.ChatMessageRoleUser,
					Content: q,
				},
			},
		},
	)
	if err != nil {
		log.Printf("openAIClient.CreateChatCompletion err=%+v\n", err)
		return "抱歉，出错了，请稍后重试~"
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

func CreateImage(q string) string {
	resp, err := imageOpenAIClient.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt: q,
			N:      1,
			Size:   "512x512",
		},
	)
	if err != nil {
		log.Printf("openAIClient.CreateImage err=%+v\n", err)
		return ""
	}
	return resp.Data[0].URL
}
