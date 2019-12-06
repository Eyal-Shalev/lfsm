package lfsm

type Option func(s *State)

func InitialState(v uint64) Option {
	return func(s *State) {
		s.current = v
	}
}

func StateName(v uint64, name string) Option {
	return func(s *State) {
		s.stateNames[v] = name
	}
}

func StateNames(m StateNameMap) Option {
	return func(s *State) {
		for v,name := range m {
			s.stateNames[v] = name
		}
	}
}