package agent

import (
	"fmt"
	"pixie/pixie"
	"strings"
	"sync"
)

func debug() string {
	agentsMtx.RLock()
	defer agentsMtx.RUnlock()

	lines := []string{}

	agents := Agents()
	for _, agent := range agents {
		lines = append(lines, fmt.Sprintf("id: %s, pixie: %s", agent.id, agent.Pixie().Debug()))
	}

	return strings.Join(lines, "\n")
}

func pickupAgent(userID string) (*Agent, error) {
	agentsMtx.RLock()
	agents := Agents()
	agent, ok := agents[userID]
	agentsMtx.RUnlock()
	if ok {
		return agent, nil
	}

	firstPixie, err := pixie.God().Pickup(pixie.Name_NormalPixie)
	if err != nil {
		return nil, err
	}

	agent = &Agent{
		id: userID,
		px: firstPixie.Summon(),
	}

	agentsMtx.Lock()
	agents[userID] = agent
	agentsMtx.Unlock()

	return agent, nil
}

func ExecuteCommand(userID string, commandData string) Message {
	agent, err := pickupAgent(userID)
	if err != nil {
		return Message{
			Title:   "Agent Error",
			Content: err.Error(),
		}
	}

	if err := agent.Lock(); err != nil {
		return Message{
			Title:   "Agent Error",
			Content: err.Error(),
		}
	}
	defer agent.Unlock()

	command := ToCommand(commandData)

	handler, ok := commandHandlers[command.Type]
	if !ok {
		return Message{
			Title: "unsupport command",
		}
	}

	return handler(agent, command.Content)
}

var once sync.Once
var agentsMtx sync.RWMutex

func Agents() map[string]*Agent {
	once.Do(func() {
		agents = map[string]*Agent{}
		commandHandlers = map[string]func(*Agent, string) Message{
			CommandType_ListGodPixies: CommandListGodPixies,
			CommandType_FocusPixie:    CommandFocusPixie,
			CommandType_Help:          CommandHelp,
			CommandType_Chat:          CommandChat,
			CommandType_Debug:         CommandDebug,
		}
	})

	return agents
}

var agents map[string]*Agent
var commandHandlers map[string]func(*Agent, string) Message

type Agent struct {
	id string
	px pixie.Pixie

	mtx    sync.Mutex
	isBusy bool
}

func (a *Agent) Lock() error {
	if a.isBusy {
		return fmt.Errorf("agent is busy")
	}

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if a.isBusy {
		return fmt.Errorf("agent is busy")
	}

	a.isBusy = true

	return nil
}

func (a *Agent) Unlock() {
	a.isBusy = false
}

func (a *Agent) SummonPixie(px pixie.Pixie) {
	a.px = px
}

func (a *Agent) Pixie() pixie.Pixie {
	if a.px == nil {
		godPixie, _ := pixie.God().Pickup(pixie.Name_NormalPixie)
		if godPixie != nil {
			a.px = godPixie.Summon()
		}
	}

	return a.px
}
