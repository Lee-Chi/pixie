package pixie

import (
	"context"
	"encoding/json"
	"fmt"
	"pixie/vocabulary"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	EnglishTeacherSkill_Min int = iota
	EnglishTeacherSkill_Vocabulary
	EnglishTeacherSkill_Dialogue
	EnglishTeacherSkill_Correct
	EnglishTeacherSkill_Translate
	EnglishTeacherSkill_Learn
	EnglishTeacherSkill_Max
)

var EnglishTeacherSkillMeans map[int]string = map[int]string{
	EnglishTeacherSkill_Vocabulary: "Vocabulary",
	EnglishTeacherSkill_Dialogue:   "Dialogue",
	EnglishTeacherSkill_Correct:    "Correct",
	EnglishTeacherSkill_Translate:  "Translate",
	EnglishTeacherSkill_Learn:      "Learn",
}

var EnglishTeacherSkillTasks map[int]string = map[int]string{
	EnglishTeacherSkill_Vocabulary: "Explain the following words in English. Present them in a table format, and the table should include the word, part of speech, definition, and example sentences: \"\"\"%s\"\"\"",
	EnglishTeacherSkill_Dialogue:   "Can we have a conversation about %s?",
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
	return "${skill} - $Vocabulary, $Dialogue, $Correct, $Translate, $Learn"
}

func (p *EnglishTeacherPixie) Resolve(ctx context.Context, request Request) (string, error) {
	message := request.Payload

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
			p.skill = EnglishTeacherSkill_Dialogue
			reply = "Ok, 要討論什麼話題?"
			p.turns = []Turn{}
		case "Correct":
			p.skill = EnglishTeacherSkill_Correct
			reply = "Ok, 想校正什麼呢?"
		case "Translate":
			p.skill = EnglishTeacherSkill_Translate
			reply = "Ok, 想翻譯什麼呢?"
		case "Learn":
			p.skill = EnglishTeacherSkill_Learn
			reply = "Ok, 來繼續學習八"
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
		if p.skill == EnglishTeacherSkill_Dialogue {
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
		} else if p.skill == EnglishTeacherSkill_Learn {
			action, payload := func(msg string) (string, string) {
				var action, payload string
				if len(msg) > 0 {
					action = msg[:1]
					payload = msg[1:]
				}

				return action, payload
			}(message)
			switch action {
			case "<":
				voc, err := vocabulary.Previous(ctx, request.UserId)
				if err != nil {
					return "", err
				}
				return voc.Marshal(), nil
			case ">":
				voc, err := vocabulary.Next(ctx, request.UserId)
				if err != nil {
					return "", err
				}
				return voc.Marshal(), nil
			case "#":
				vocabularyId, _ := primitive.ObjectIDFromHex(payload)

				if err := vocabulary.Toggle(ctx, request.UserId, vocabularyId); err != nil {
					return "", err
				}
				return "Ok", nil
			case "^":
				vocabularyId, _ := primitive.ObjectIDFromHex(payload)

				voc, err := vocabulary.At(ctx, request.UserId, vocabularyId)
				if err != nil {
					return "", err
				}
				parts := []string{}
				for _, def := range voc.Definitions {
					parts = append(parts, fmt.Sprintf("%s (%s)", voc.Word, def.PartOfSpeech))
				}
				lines := []string{
					fmt.Sprintf("###%s###", strings.Join(parts, ",")),
					"Provide three example sentences above of text, and present them in this format.",
					"Words (part of speech)",
					"1. Example sentence 1",
					"2. Example sentence 2",
					"3. Example sentence 3",
				}

				message = strings.Join(lines, "\n")
			case "~":
				page, err := strconv.ParseInt(payload, 10, 64)
				if err != nil {
					return "", err
				}
				vocs, err := vocabulary.BrowseBookmark(ctx, request.UserId, page)
				if err != nil {
					return "", err
				}
				lines := []string{}
				for _, voc := range vocs {
					lines = append(lines, voc.Marshal())
				}
				return strings.Join(lines, "\n"), nil
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			})
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

	if p.skill == EnglishTeacherSkill_Dialogue {
		p.turns = append(p.turns, Turn{
			User: message,
			AI:   reply,
		})
	}

	return reply, nil
}
