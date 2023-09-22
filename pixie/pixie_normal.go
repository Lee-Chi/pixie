package pixie

import (
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type NormalPixie struct {
	name string
	role string
}

func NewNormal() Pixie {
	return &NormalPixie{
		name: "pixie",
		role: "",
	}
}

func (p NormalPixie) Welcome() string {
	return welcome(p.name)
}
func (p NormalPixie) IntroduceSelf() string {
	return "@normal |- 找不到適合的小精靈，選我就對了"
}
func (p NormalPixie) Help() string {
	return "#{角色} - #專家"
}

func (p *NormalPixie) ReplyMessage(message string) (string, error) {
	if strings.HasPrefix(message, "!") {
		return p.Help(), nil
	} else if strings.HasPrefix(message, "#") {
		p.role = strings.TrimPrefix(message, "#")

		return "ok", nil
	}

	messages := make([]openai.ChatCompletionMessage, 0)

	if p.role != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("你現在是%s", p.role),
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	reply, err := chat(messages)
	if err != nil {
		return "", err
	}

	return reply, nil
}
