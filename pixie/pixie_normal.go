package pixie

import (
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
)

type NormalPixie struct {
	name string
	role string

	wrapper struct {
		Name string `bson:"name"`
		Role string `bson:"role"`
	}
	needSave bool
}

func NewNormal() Pixie {
	return &NormalPixie{
		name: "pixie",
		role: "",
	}
}

func (p NormalPixie) Name() string {
	return p.name
}

func (p NormalPixie) NeedSave() bool {
	return p.needSave
}
func (p *NormalPixie) Unmarshal(marshal string) error {
	if err := bson.Unmarshal([]byte(marshal), &p.wrapper); err != nil {
		return err
	}

	p.name = p.wrapper.Name
	p.role = p.wrapper.Role

	return nil
}

func (p *NormalPixie) Marshal() string {
	p.needSave = false

	p.wrapper.Name = p.name
	p.wrapper.Role = p.role

	marshal, _ := bson.Marshal(p.wrapper)

	return string(marshal)
}

func (p NormalPixie) Debug() string {
	return fmt.Sprintf("{name: %s, role: %s}", p.name, p.role)
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
		p.needSave = true

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
