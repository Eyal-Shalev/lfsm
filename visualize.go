package lfsm

import (
	"bytes"
	"fmt"
)

// String returns the Graphviz representation of this state machine.
//
// Example:
// 	fmt.Println(lfsm.NewState(lfsm.Constraints{0:{1},1:{0}}, lfsm.StateNames{0:"opened",1:"closed"}))
//  // Output: digraph g{s[label="",shape=none,height=.0,width=.0];s->n0;n0[label="opened",style=filled];n1[label="closed"];n0->n1;n1->n0;}
//  // See: https://dreampuf.github.io/GraphvizOnline/#digraph%20g%7Bs%5Blabel%3D%22%22%2Cshape%3Dnone%2Cheight%3D.0%2Cwidth%3D.0%5D%3Bs-%3En0%3Bn0%5Blabel%3D%22opened%22%2Cstyle%3Dfilled%5D%3Bn1%5Blabel%3D%22closed%22%5D%3Bn0-%3En1%3Bn1-%3En0%3B%7D
// See: https://www.graphviz.org/
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
		for dst := range dsts {
			_,_ = fmt.Fprintf(buf, "n%d->n%d;", src, dst)
		}
	}

	_,_ = fmt.Fprint(buf, "}")
	return buf.String()
}
