package agent

import (
	"fmt"
	"pixie/pixie"
	"strings"
)

type Command struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func ToCommand(data string) Command {
	if strings.HasPrefix(data, CommandType_ListGodPixies) {
		return Command{
			Type:    CommandType_ListGodPixies,
			Content: "",
		}
	} else if strings.HasPrefix(data, CommandType_FocusPixie) {
		pixieName := strings.TrimPrefix(data, CommandType_FocusPixie)

		return Command{
			Type:    CommandType_FocusPixie,
			Content: pixieName,
		}
	} else if strings.HasPrefix(data, CommandType_Debug) {
		return Command{
			Type:    CommandType_Debug,
			Content: "",
		}
	} else if strings.HasPrefix(data, CommandType_Help) {
		return Command{
			Type:    CommandType_Help,
			Content: "",
		}
	}

	return Command{
		Type:    CommandType_Chat,
		Content: data,
	}
}

const CommandType_ListGodPixies string = "/"

func CommandListGodPixies(agent *Agent, content string) Message {
	return Message{
		Title:   "可召喚的小精靈",
		Content: pixie.God().ListPixies(),
	}
}

const CommandType_FocusPixie string = "@"

func CommandFocusPixie(agent *Agent, content string) Message {
	pixieName := content
	if pixieName == "" {
		pixieName = pixie.Name_NormalPixie
	}

	godPixie, err := pixie.God().Pickup(pixieName)
	if err != nil {
		return Message{
			Title: fmt.Sprintf("pixie %s is not found", pixieName),
		}
	}

	agent.SummonPixie(godPixie.Summon())

	return Message{
		Title:   agent.px.Welcome(),
		Content: agent.px.Help(),
	}
}

const CommandType_Help string = "?"

func CommandHelp(agent *Agent, context string) Message {
	px := agent.Pixie()

	return Message{
		Title:   px.Welcome(),
		Content: "也可以使用 / 看看有沒有更適合您的小精靈",
	}
}

const CommandType_Chat string = ""

func CommandChat(agent *Agent, content string) Message {
	px := agent.Pixie()

	if px == nil {
		return Message{
			Title: fmt.Sprintf("請先指定小精靈"),
		}
	}

	message, err := px.ReplyMessage(content)
	if err != nil {
		return Message{
			Title:   fmt.Sprintf("哎呀，發生了一點點小錯誤"),
			Content: err.Error(),
		}
	}

	return Message{
		Content: message,
	}
}

const CommandType_Debug string = ":"

func CommandDebug(agent *Agent, content string) Message {
	content = debug()
	return Message{
		Title:   "Debug",
		Content: content,
	}
}
