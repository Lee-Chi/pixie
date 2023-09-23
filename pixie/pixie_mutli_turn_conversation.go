package pixie

import (
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Turn struct {
	User string
	AI   string
}
type MultiTurnConversationPixie struct {
	name  string
	turns []Turn
	role  string
}

func NewMultiTurnConversation() Pixie {
	return &MultiTurnConversationPixie{
		name:  "chatter",
		turns: []Turn{},
		role:  "",
	}
}

func (p MultiTurnConversationPixie) Debug() string {
	return fmt.Sprintf("{name: %s, role: %s, size_of_turns: %d}", p.name, p.role, len(p.turns))
}

func (p MultiTurnConversationPixie) Welcome() string {
	return welcome(p.name)
}
func (p MultiTurnConversationPixie) IntroduceSelf() string {
	return "@chatter |- 我不是金魚腦，我會記得你說的話唷"
}
func (p MultiTurnConversationPixie) Help() string {
	return "#{角色} - #面試官 ..."
}

func (p *MultiTurnConversationPixie) ReplyMessage(message string) (string, error) {
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
			Content: p.role,
		})
	}
	for _, turn := range p.turns {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: turn.User,
		})
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: turn.AI,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	reply, err := chat(messages)
	if err != nil {
		return "", nil
	}

	p.turns = append(p.turns, Turn{
		User: message,
		AI:   reply,
	})

	return reply, nil
}
