package scheduler

import (
	"fmt"
	"sync"
)

var (
	pool map[string]*Routine
	mtx  sync.Mutex
)

func init() {
	pool = map[string]*Routine{}
}

func register(name string, routine *Routine) error {
	mtx.Lock()
	defer mtx.Unlock()

	if _, ok := pool[name]; ok {
		return fmt.Errorf("%s exist", name)
	}

	pool[name] = routine

	return nil
}

func unregister(name string) {
	mtx.Lock()
	defer mtx.Unlock()

	delete(pool, name)
}
