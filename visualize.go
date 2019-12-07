package lfsm

import (
	"bytes"
	"fmt"
)

// String returns the Graphviz representation of this state machine.
func (s *State) String() string {
	buf := &bytes.Buffer{}
	_,_ = fmt.Fprint(buf, "digraph g{")
	_,_ = fmt.Fprintf(buf, `s[label="",shape=none,height=.0,width=.0];s->n%d;`, s.initial)

	for v,name := range s.stateNames {
		if v == s.Current() {
			_,_ = fmt.Fprintf(buf, `n%d[label="%s",style=filled];`, v, name)
		} else {
			_,_ = fmt.Fprintf(buf, `n%d[label="%s"];`, v, name)
		}
	}

	for src,dsts := range s.transitions {
		for dst,_ := range dsts {
			_,_ = fmt.Fprintf(buf, "n%d->n%d;", src, dst)
		}
	}

	_,_ = fmt.Fprint(buf, "}")
	return buf.String()
}
