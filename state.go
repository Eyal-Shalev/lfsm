package lfsm

import (
	"strconv"
	"sync/atomic"
)

type StateNameMap map[uint64]string
func (m StateNameMap) find(v uint64) string {
	name, ok := m[v]
	if !ok {
		name = strconv.Itoa(int(v))
	}
	return name
}

type transition struct {
	src, dst uint64
	stateNames StateNameMap
}

type transFn func() error
type transFnMap map[uint64]transFn

type State struct {
	current     uint64
	transitions map[uint64]transFnMap
	stateNames StateNameMap
}

func (s *State) Current() uint64 {
	return atomic.LoadUint64(&s.current)
}

func (s *State) CurrentName() string {
	return s.stateNames.find(atomic.LoadUint64(&s.current))
}

// Transition tries to change the state.
// Returns an error if the transition failed.
func (s *State) Transition(dst uint64) error {
	return s.TransitionFrom(atomic.LoadUint64(&s.current), dst)
}

// TransitionFrom tries to change the state.
// Returns an error if the transition failed.
func (s *State) TransitionFrom(src, dst uint64) error {
	f, ok := s.transitions[src][dst]
	if !ok {
		return &InvalidTransitionError{src, dst, s.stateNames}
	}
	return f()
}

func (s *State) makeTransFn(src, dst uint64) transFn {
	return func() error {
		if !atomic.CompareAndSwapUint64(&s.current, src, dst) {
			return &TransitionError{src, dst, s.stateNames}
		}
		return nil
	}
}

type Constraints map[uint64][]uint64

func NewState(m Constraints, opts ...Option) *State {
	s := State{
		transitions: make(map[uint64]transFnMap, len(m)),
		stateNames: make(StateNameMap, len(m)),
	}

	for src, dsts := range m {
		s.transitions[src] = make(transFnMap)
		for _, dst := range dsts {
			s.transitions[src][dst] = s.makeTransFn(src, dst)
		}
	}

	for _,o := range opts {
		o(&s)
	}

	return &s
}
