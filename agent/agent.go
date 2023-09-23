package agent

import (
	"context"
	"fmt"
	"pixie/db"
	"pixie/db/model"
	"pixie/db/mongo"
	"pixie/pixie"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Field_UserId mongo.F = mongo.Field("user_id")
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

func pickupAgent(ctx context.Context, userId string) (*Agent, error) {
	agentsMtx.RLock()
	agents := Agents()
	agent, ok := agents[userId]
	agentsMtx.RUnlock()
	if ok {
		return agent, nil
	}

	agent, err := NewAgent(ctx, userId)
	if err != nil {
		return nil, err
	}

	agentsMtx.Lock()
	agents[userId] = agent
	agentsMtx.Unlock()

	return agent, nil
}

func ExecuteCommand(ctx context.Context, userId string, commandData string) Message {
	agent, err := pickupAgent(ctx, userId)
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

	message := handler(agent, command.Content)

	agent.Save(ctx)

	return message
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

func NewAgent(ctx context.Context, userId string) (*Agent, error) {
	md := struct {
		Id          primitive.ObjectID `bson:"_id"`
		model.Agent `bson:"-,inline"`
	}{}
	if err := db.Pixie().Collection(model.CAgent).FindOneOrZero(
		ctx,
		Field_UserId.Equal(userId),
		&md,
	); err != nil {
		return nil, err
	}

	agent := &Agent{
		id: userId,
		px: nil,
	}

	if md.Id.IsZero() {

		firstPixie, err := pixie.God().Pickup(pixie.Name_NormalPixie)
		if err != nil {
			return nil, err
		}

		agent.px = firstPixie.Summon()

		md.UserId = userId
		md.Pixie = model.Pixie{
			Name:    agent.px.Name(),
			Payload: agent.px.Marshal(),
		}

		if err := db.Pixie().Collection(model.CAgent).Insert(
			ctx,
			md.Agent,
		); err != nil {
			return nil, err
		}

		return agent, nil
	}

	if err := agent.px.Unmarshal(md.Pixie.Payload); err != nil {
		return nil, err
	}

	return agent, nil
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

func (a *Agent) Save(ctx context.Context) error {
	if !a.px.NeedSave() {
		return nil
	}

	md := model.Agent{
		UserId: a.id,
		Pixie: model.Pixie{
			Name:    a.px.Name(),
			Payload: a.px.Marshal(),
		},
	}

	if err := db.Pixie().Collection(model.CAgent).Update(
		ctx,
		Field_UserId.Equal(a.id),
		md,
	); err != nil {
		return err
	}

	return nil
}
