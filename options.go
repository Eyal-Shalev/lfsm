package lfsm

// Can be used to alter the state struct during initialization.
type option func(s *State)

// InitialState sets the initial state of the state machine
func InitialState(v uint32) option {
	return func(s *State) {
		s.initial = v
		s.current = v
	}
}

// StateName sets an alias to a state integer.
func StateName(v uint32, name string) option {
	return func(s *State) {
		s.stateNames[v] = name
	}
}

// StateNames is like StateName but sets aliases to multiple state integers.
func StateNames(m StateNameMap) option {
	return func(s *State) {
		for v,name := range m {
			s.stateNames[v] = name
		}
	}
}