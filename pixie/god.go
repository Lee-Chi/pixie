package pixie

import (
	"fmt"
	"strings"
	"sync"
)

type GodPool struct {
	names  []string
	pixies map[string]GodPixie
}

type GodPixie struct {
	Pixie
	Summon func() Pixie
}

var godPool *GodPool

func BuildGodPool() error {

	return nil
}

var once sync.Once

func God() *GodPool {
	once.Do(func() {
		godPool = &GodPool{
			names: []string{
				Name_NormalPixie,
				Name_ProgrammerPixie,
				Name_MultiTurnConversation,
			},
			pixies: map[string]GodPixie{
				Name_NormalPixie: {
					Pixie:  Summon_NormalPixie(),
					Summon: Summon_NormalPixie,
				},
				Name_ProgrammerPixie: {
					Pixie:  Summon_ProgrammerPixie(),
					Summon: Summon_ProgrammerPixie,
				},
				Name_MultiTurnConversation: {
					Pixie:  Summon_MultiTurnConversation(),
					Summon: Summon_MultiTurnConversation,
				},
			},
		}
	})

	return godPool
}
func (p GodPool) ListPixies() string {
	lines := []string{}
	for _, name := range p.names {
		line := p.pixies[name].IntroduceSelf()
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (p GodPool) Pickup(name string) (*GodPixie, error) {
	if found, ok := p.pixies[name]; ok {
		return &found, nil
	}

	return nil, fmt.Errorf("pixie %s is not found", name)
}

func (gp GodPixie) Fork() Pixie {
	return gp.Summon()
}
