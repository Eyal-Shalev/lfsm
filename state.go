package lfsm

import (
	"strconv"
	"sync/atomic"
)

type transitionMap map[uint32]map[uint32]bool

// State is the structs that holds the current state, the available transitions and other options.
type State struct {
	current     uint32
	transitions transitionMap
	stateNames  StateNames
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
	if _, ok := s.transitions[src][dst]; !ok {
		return NewInvalidTransitionError(src, dst, s.stateNames)
	}
	if !atomic.CompareAndSwapUint32(&s.current, src, dst) {
		return NewFailedTransitionError(src, dst, s.stateNames)
	}
	return nil
}

// Transition tries to change the state.
// It uses the current state as the source state, if you want to specify the source state use TransitionFrom instead.
// Returns an error if the transition failed.
func (s *State) Transition(dst uint32) error {
	return s.TransitionFrom(atomic.LoadUint32(&s.current), dst)
}

// NewState creates a new State Machine.
func NewState(m Constraints, opts ...option) *State {
	s := State{
		transitions: make(transitionMap, len(m)),
		stateNames: make(StateNames, len(m)),
	}

	for src, dsts := range m {
		s.transitions[src] = make(map[uint32]bool)
		for _, dst := range dsts {
			s.transitions[src][dst] = true
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

// StateNames holds a mapping between the state (in its integer form) to its alias.
type StateNames map[uint32]string
func (m StateNames) find(v uint32) string {
	name, ok := m[v]
	if !ok {
		name = strconv.Itoa(int(v))
	}
	return name
}
func (m StateNames) apply(s *State) {
	for v,name := range m {
		s.stateNames[v] = name
	}
}
