package lfsm

// Can be used to alter the state struct during initialization.
type option interface {
	apply(s *State)
}

type optionFn func(s *State)
func (o optionFn) apply(s *State) {
	o(s)
}

// InitialState sets the initial state of the state machine
func InitialState(v uint32) option {
	return optionFn(func(s *State) {
		s.initial = v
		s.current = v
	})
}

// StateName sets an alias to a state integer.
func StateName(v uint32, name string) option {
	return optionFn(func(s *State) {
		s.stateNames[v] = name
	})
}

func (m StateNameMap) apply(s *State) {
	for v,name := range m {
		s.stateNames[v] = name
	}
}