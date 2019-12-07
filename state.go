package lfsm

import (
	"strconv"
	"sync/atomic"
)

type transition struct {
	src, dst uint32
	stateNames StateNameMap
}
type transFn func() error
type transFnMap map[uint32]transFn

type State struct {
	current     uint32
	transitions map[uint32]transFnMap
	stateNames  StateNameMap
	initial     uint32
}

// Current returns the current state.
func (s *State) Current() uint32 {
	return atomic.LoadUint32(&s.current)
}

// CurrentName returns the alias for the current state.
// If no alias is defined, the state integer will be returned in its string version.
func (s *State) CurrentName() string {
	return s.stateNames.find(atomic.LoadUint32(&s.current))
}

// TransitionFrom tries to change the state.
// Returns an error if the transition failed.
func (s *State) TransitionFrom(src, dst uint32) error {
	f, ok := s.transitions[src][dst]
	if !ok {
		return &InvalidTransitionError{src, dst, s.stateNames}
	}
	return f()
}

// Transition tries to change the state.
// It uses the current state as the source state, if you want to specify the source state use TransitionFrom instead.
// Returns an error if the transition failed.
func (s *State) Transition(dst uint32) error {
	return s.TransitionFrom(atomic.LoadUint32(&s.current), dst)
}

func (s *State) makeTransFn(src, dst uint32) transFn {
	return func() error {
		if !atomic.CompareAndSwapUint32(&s.current, src, dst) {
			return &TransitionError{src, dst, s.stateNames}
		}
		return nil
	}
}

// NewState creates a new State Machine.
func NewState(m Constraints, opts ...option) *State {
	s := State{
		transitions: make(map[uint32]transFnMap, len(m)),
		stateNames: make(StateNameMap, len(m)),
	}

	for src, dsts := range m {
		s.transitions[src] = make(transFnMap)
		for _, dst := range dsts {
			s.transitions[src][dst] = s.makeTransFn(src, dst)
		}
	}

	for _,o := range opts {
		o.apply(&s)
	}

	return &s
}

// Constraints defines the possible transition for this state machine.
//
// The map keys describe the source states, and their values are the valid target destinations.
type Constraints map[uint32][]uint32

// StateNameMap holds a mapping between the state (in its integer form) to its alias.
type StateNameMap map[uint32]string
func (m StateNameMap) find(v uint32) string {
	name, ok := m[v]
	if !ok {
		name = strconv.Itoa(int(v))
	}
	return name
}
