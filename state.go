package lfsm

import (
	"sync/atomic"
)

type transition struct {
	src, dst uint64
}

type transFn func() error
type transFnMap map[uint64]transFn

type State struct {
	current     uint64
	transitions map[uint64]transFnMap
}

func (s *State) Current() uint64 {
	return atomic.LoadUint64(&s.current)
}

// Transition tries to change the state.
// Returns an error if the transition failed.
func (s *State) Transition(dst uint64) error {
	src := atomic.LoadUint64(&s.current)
	f, ok := s.transitions[src][dst]
	if !ok {
		return &InvalidTransitionError{src, dst}
	}
	return f()
}

func (s *State) makeTransFn(src, dst uint64) transFn {
	return func() error {
		if !atomic.CompareAndSwapUint64(&s.current, src, dst) {
			return &TransitionError{src, dst}
		}
		return nil
	}
}

type Constraints map[uint64][]uint64

func NewState(initial uint64, m Constraints) *State {
	s := State{
		current:     initial,
		transitions: make(map[uint64]transFnMap, len(m)),
	}

	for src, dsts := range m {
		s.transitions[src] = make(transFnMap)
		for _, dst := range dsts {
			s.transitions[src][dst] = s.makeTransFn(src, dst)
		}
	}

	return &s
}
