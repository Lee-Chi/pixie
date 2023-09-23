package pixie

import (
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const (
	ProgrammerSkill_Min int = iota
	ProgrammerSkill_WriteCode
	ProgrammerSkill_ReadCode
	ProgrammerSkill_RefactorCode
	ProgrammerSkill_FixBug
	ProgrammerSkill_WriteTest
	ProgrammerSkill_WriteRegex
	ProgrammerSkill_Max
)

var ProgrammerSkillMeans map[int]string = map[int]string{
	ProgrammerSkill_WriteCode:    "Write",
	ProgrammerSkill_ReadCode:     "Read",
	ProgrammerSkill_RefactorCode: "Refactor",
	ProgrammerSkill_FixBug:       "Fix",
	ProgrammerSkill_WriteTest:    "Test",
	ProgrammerSkill_WriteRegex:   "Regex",
}

var ProgrammerSkillTasks map[int]string = map[int]string{
	ProgrammerSkill_WriteCode:    "%s",
	ProgrammerSkill_ReadCode:     "請告訴我以下程式碼在做什麼。%s",
	ProgrammerSkill_RefactorCode: "你也是個clean code專家，我有以下的程式碼，請用更乾淨簡潔的方式改寫，並且說明為什麼要這樣重構。%s",
	ProgrammerSkill_FixBug:       "檢查以下程式碼有什麼問題。%s",
	ProgrammerSkill_WriteTest:    "對以下程式碼寫一個測試，提供5個案例，案例要包含到極端的狀況",
	ProgrammerSkill_WriteRegex:   "你也是個regex專家，寫一個regex，需求是%s",
}

type ProgrammerPixie struct {
	name     string
	language string
	skill    int
}

func NewProgrammer() Pixie {
	return &ProgrammerPixie{
		name:     "programmer",
		language: "",
		skill:    -1,
	}
}

func (p ProgrammerPixie) Debug() string {
	skill := ProgrammerSkillMeans[p.skill]
	return fmt.Sprintf("{name: %s, language: %s, skill: %s}", p.name, p.language, skill)
}

func (p ProgrammerPixie) Welcome() string {
	return welcome(p.name)
}

func (p ProgrammerPixie) IntroduceSelf() string {
	return "@programmer |- 我是程式小天才"
}

func (p ProgrammerPixie) Help() string {
	return strings.Join([]string{
		"#{程式語言} - #golang, #javascript, #python ...",
		"${skill} - $Write, $Read, $Refactor, $Fix, Test, Regex",
	}, "\n")
}

func (p *ProgrammerPixie) ReplyMessage(message string) (string, error) {
	if strings.HasPrefix(message, "!") {
		return p.Help(), nil
	} else if strings.HasPrefix(message, "#") {
		p.language = strings.TrimPrefix(message, "#")

		return "ok", nil
	} else if strings.HasPrefix(message, "$") {
		skill := strings.TrimPrefix(message, "$")
		switch skill {
		case "Write":
			p.skill = ProgrammerSkill_WriteCode
			return "Ok, 想實現什麼功能?", nil
		case "Read":
			p.skill = ProgrammerSkill_ReadCode
			return "Ok, 給我程式碼", nil
		case "Refactor":
			p.skill = ProgrammerSkill_RefactorCode
			return "Ok, 給我程式碼", nil
		case "Fix":
			p.skill = ProgrammerSkill_FixBug
			return "Ok, 給我程式碼", nil
		case "Test":
			p.skill = ProgrammerSkill_WriteTest
			return "Ok, 給我程式碼", nil
		case "Regex":
			p.skill = ProgrammerSkill_WriteRegex
			return "Ok, 想做什麼需求?", nil
		}

		return "Oh no", nil
	}

	messages := make([]openai.ChatCompletionMessage, 0)

	if p.language != "" && p.skill > ProgrammerSkill_Min && p.skill < ProgrammerSkill_Max {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("你現在是一個%s專家", p.language),
		})

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(ProgrammerSkillTasks[p.skill], message),
		})
	} else {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("你現在是一個程式設計師"),
		})

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		})
	}

	reply, err := chat(messages)
	if err != nil {
		return "", err
	}

	return reply, nil
}
