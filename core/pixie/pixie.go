package pixie

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client = nil

func Build(openaiToken string) {
	openaiClient = openai.NewClient(openaiToken)
}

func Chat(messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("ChatGPT is busy")
	}

	return resp.Choices[0].Message.Content, nil
}
