// +build ignore

package main

import (
	"log"
	"sync"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

const (
	unlocked uint32 = iota
	locked
)

type Mutex struct {
	state *lfsm.State
}

var constraints = lfsm.Constraints{
	unlocked: {locked},
	locked:   {unlocked},
}

func NewMutex() *Mutex {
	return &Mutex{lfsm.NewState(constraints, lfsm.InitialState(unlocked))}
}

func (m Mutex) Lock() {
	for {
		if m.state.Transition(locked) == nil {
			return
		}
		time.Sleep(time.Nanosecond)
	}
}

func (m Mutex) Unlock() {
	if err := m.state.Transition(unlocked); err != nil {
		panic("unlock of unlocked mutex")
	}
}

func HammerMutex(m *Mutex, loops int, counter *int, wg *sync.WaitGroup) {
	for i := 0; i < loops; i++ {
		m.Lock()
		*counter += 1
		m.Unlock()
	}
	wg.Done()
}

func main() {
	m := NewMutex()
	counter := new(int)
	wg := new(sync.WaitGroup)
	I, J := 1000, 1000
	for i := 0; i < I; i++ {
		wg.Add(1)
		go HammerMutex(m, J, counter, wg)
	}
	wg.Wait()
	if *counter != I*J {
		log.Fatalf("counter is: %d instead of %d", *counter, I*J)
	}
}
