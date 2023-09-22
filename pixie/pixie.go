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

func chat(messages []openai.ChatCompletionMessage) (string, error) {
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

func welcome(name string) string {
	return fmt.Sprintf("您好，現在是小精靈 %s 為您服務喔", name)
}

type Mode string

type Pixie interface {
	Welcome() string
	IntroduceSelf() string
	Help() string
	ReplyMessage(message string) (string, error)
}

var (
	Name_NormalPixie   = "normal"
	Summon_NormalPixie = NewNormal

	Name_ProgrammerPixie   = "programmer"
	Summon_ProgrammerPixie = NewProgrammer

	Name_MultiTurnConversation   = "chatter"
	Summon_MultiTurnConversation = NewMultiTurnConversation

	Name_EnglishTeacher   = "british"
	Summon_EnglishTeacher = NewEnglishTeacher
)
