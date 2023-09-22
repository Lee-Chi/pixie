package agent

import (
	"fmt"
	"strings"
)

type Message struct {
	Title   string
	Content string
}

func (message Message) Marshal() string {
	lines := []string{}
	if message.Title != "" {
		lines = append(lines, fmt.Sprintf("*%s*\n", message.Title))
	}

	lines = append(lines, message.Content)

	return strings.Join(lines, "\n")
}
