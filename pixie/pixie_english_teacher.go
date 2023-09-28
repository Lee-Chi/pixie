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
	EnglishTeacherSkill_Correct
	EnglishTeacherSkill_Explain
	EnglishTeacherSkill_TranslateIntoEnglish
	EnglishTeacherSkill_TranslateIntoChinese
	EnglishTeacherSkill_Dialogue
	EnglishTeacherSkill_Max
)

var EnglishTeacherSkillMeans map[int]string = map[int]string{
	EnglishTeacherSkill_Vocabulary:           "Vocabulary",
	EnglishTeacherSkill_Correct:              "Correct",
	EnglishTeacherSkill_Explain:              "Explain",
	EnglishTeacherSkill_TranslateIntoEnglish: "TranslateIntoEnglish",
	EnglishTeacherSkill_TranslateIntoChinese: "TranslateIntoChinese",
	EnglishTeacherSkill_Dialogue:             "Dialogue",
}

var EnglishTeacherSkillTasks map[int]string = map[int]string{
	EnglishTeacherSkill_Correct:              "Check the following text for grammar or spelling errors: %s",
	EnglishTeacherSkill_Explain:              "Explain the following text: %s",
	EnglishTeacherSkill_TranslateIntoEnglish: "Translate the following text into English: %s",
	EnglishTeacherSkill_TranslateIntoChinese: "Translate the following text into Tranditional Chinese: %s",
	EnglishTeacherSkill_Dialogue:             "Can we have a conversation about %s?",
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
	return "${skill} - $Vocabulary, $Correct, $Explain, $TranslateIntoEnglish, $TranslateIntoChinese"
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
		case "Correct":
			p.skill = EnglishTeacherSkill_Correct
			reply = "Ok, 想校正什麼呢?"
		case "Explain":
			p.skill = EnglishTeacherSkill_Explain
			reply = "Ok, 想知道什麼呢?"
		case "TranslateIntoEnglish":
			p.skill = EnglishTeacherSkill_TranslateIntoEnglish
			reply = "Ok, 想翻譯什麼呢?"
		case "TranslateIntoChinese":
			p.skill = EnglishTeacherSkill_TranslateIntoChinese
			reply = "Ok, 想翻譯什麼呢?"
		case "Dialogue":
			p.skill = EnglishTeacherSkill_Dialogue
			reply = "Ok, 要討論什麼話題?"
			p.turns = []Turn{}
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
		Content: "你現在是精通英文和中文的語言專家",
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
		} else if p.skill == EnglishTeacherSkill_Vocabulary {
			action, payload := func(msg string) (string, string) {
				var action, payload string
				if len(msg) > 0 {
					action = msg[:1]

					payload = msg[1:]
					payload = strings.Trim(payload, " ")
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
					fmt.Sprintf(`對以下英文單字做解釋並提供3個例句:%s。以json方式輸出，key包含part_of_speech,explain,sentences。範例: {"bark":{"part_of_speech":"vt","explain":"to shout or speak loudly and insistently","sentences":["The dog barked at the intruder.","The coach barked orders at the players.","He barked out a command to stop."]}}，且美化格式後再回覆。`, strings.Join(parts, ",")),
				}

				message = strings.Join(lines, "\n")
			case "&":
				voc, err := vocabulary.Find(ctx, payload)
				if err != nil {
					return "", err
				}

				return voc.Marshal(), nil
			case "*":
				message = fmt.Sprintf(`對以下英文單字做解釋並提供3個例句:%s。以json方式輸出，key包含part_of_speech,explain,sentences。範例: {"bark":{"part_of_speech":"vt","explain":"to shout or speak loudly and insistently","sentences":["The dog barked at the intruder.","The coach barked orders at the players.","He barked out a command to stop."]}}，且美化格式後再回覆。`, payload)
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
