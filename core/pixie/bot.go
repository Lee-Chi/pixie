package pixie

import "github.com/sashabaranov/go-openai"

const (
	Command_List = "/"
)

var Bot B = B{}

type B struct{}

func (b B) Chat(message string) (string, error) {
	var response string = ""

	switch message {
	case Command_List:
		response = "https://hello-english-1738.de.r.appspot.com"
	}

	if response == "" {
		reply, err := Chat([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		})

		if err != nil {
			return "", err
		}

		response = reply
	}

	return response, nil
}
