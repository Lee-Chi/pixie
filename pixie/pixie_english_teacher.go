package pixie

import (
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	EnglishTeacherSkill_Min int = iota
	EnglishTeacherSkill_Vocabulary
	EngiishTeacherSkill_Dialogue
	EnglishTeacherSkill_Correct
	EnglishTeacherSkill_Translate
	EnglishTeacherSkill_Max
)

var EnglishTeacherSkillMeans map[int]string = map[int]string{
	EnglishTeacherSkill_Vocabulary: "Vocabulary",
	EngiishTeacherSkill_Dialogue:   "Dialogue",
	EnglishTeacherSkill_Correct:    "Correct",
	EnglishTeacherSkill_Translate:  "Translate",
}

var EnglishTeacherSkillTasks map[int]string = map[int]string{
	EnglishTeacherSkill_Vocabulary: "Explain the following words in English. Present them in a table format, and the table should include the word, part of speech, definition, and example sentences: \"\"\"%s\"\"\"",
	EngiishTeacherSkill_Dialogue:   "Can we have a conversation about %s?",
	EnglishTeacherSkill_Correct:    "Check the following text for grammar or spelling errors: \"\"\"%s\"\"\"",
	EnglishTeacherSkill_Translate:  "Translate the following text into english: \"\"\"%s\"\"\"",
}

type EnglishTeacherPixie struct {
	name  string
	skill int
	turns []Turn

	wrapper struct {
		Name  string `json:"name"`
		Skill int    `json:"skill"`
	}
	needSave bool
}

func NewEnglishTeacher() Pixie {
	return &EnglishTeacherPixie{
		name:  "british",
		skill: EnglishTeacherSkill_Min,
		turns: []Turn{},
	}
}

func (p EnglishTeacherPixie) Name() string {
	return p.name
}

func (p EnglishTeacherPixie) NeedSave() bool {
	return p.needSave
}

func (p *EnglishTeacherPixie) Unmarshal(marshal string) error {
	if err := json.Unmarshal([]byte(marshal), &p.wrapper); err != nil {
		return err
	}

	p.name = p.wrapper.Name
	p.skill = p.wrapper.Skill
	p.turns = []Turn{}

	return nil
}

func (p *EnglishTeacherPixie) Marshal() string {
	p.needSave = false

	p.wrapper.Name = p.name
	p.wrapper.Skill = p.skill

	marshal, _ := json.Marshal(p.wrapper)

	return string(marshal)
}

func (p EnglishTeacherPixie) Debug() string {
	skill := EnglishTeacherSkillMeans[p.skill]
	return fmt.Sprintf("{name: %s, skill: %s, size_of_turns: %d}", p.name, skill, len(p.turns))
}
func (p EnglishTeacherPixie) Welcome() string {
	return welcome(p.name)
}

func (p EnglishTeacherPixie) IntroduceSelf() string {
	return "@british |- 我是您的英文小老師"
}
func (p EnglishTeacherPixie) Help() string {
	return "${skill} - $Vocabulary, $Dialogue, $Correct, $Translate"
}

func (p *EnglishTeacherPixie) ReplyMessage(message string) (string, error) {
	if strings.HasPrefix(message, "!") {
		return p.Help(), nil
	} else if strings.HasPrefix(message, "$") {
		p.needSave = true
		skill := strings.TrimPrefix(message, "$")
		reply := ""

		switch skill {
		case "Vocabulary":
			p.skill = EnglishTeacherSkill_Vocabulary
			reply = "Ok, 想知道什麼單字?"
		case "Dialogue":
			p.skill = EngiishTeacherSkill_Dialogue
			reply = "Ok, 要討論什麼話題?"
			p.turns = []Turn{}
		case "Correct":
			p.skill = EnglishTeacherSkill_Correct
			reply = "Ok, 想校正什麼呢?"
		case "Translate":
			p.skill = EnglishTeacherSkill_Translate
			reply = "Ok, 想翻譯什麼呢?"
		}

		if reply == "" {
			p.needSave = false
			reply = "Oh no"
		}

		return reply, nil
	}

	messages := make([]openai.ChatCompletionMessage, 0)

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You're an English teacher now.",
	})

	if p.skill > EnglishTeacherSkill_Min && p.skill < EnglishTeacherSkill_Max {
		if p.skill == EngiishTeacherSkill_Dialogue {
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

			if len(messages) == 1 {
				// first
				message = fmt.Sprintf(EnglishTeacherSkillTasks[p.skill], message)
				messages = append(messages, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				})
			}
		} else {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(EnglishTeacherSkillTasks[p.skill], message),
			})
		}
	} else {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})
	}

	reply, err := chat(messages)
	if err != nil {
		return "", err
	}

	if p.skill == EngiishTeacherSkill_Dialogue {
		p.turns = append(p.turns, Turn{
			User: message,
			AI:   reply,
		})
	}

	return reply, nil
}
