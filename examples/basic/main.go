// +build ignore

package main

import (
	"log"
	"os"

	"github.com/Eyal-Shalev/lfsm"
)

func main() {
	l := log.New(os.Stdout, "", log.Lshortfile)
	const (
		opened uint32 = iota
		closed
	)
	s := lfsm.NewState(
		lfsm.Constraints{
			opened: {closed},
			closed: {opened},
		},
		lfsm.InitialState(closed),
		lfsm.StateNames{opened: "opened",closed: "closed"},
	)

	l.Println(s) // digraph g{s[label="",shape=none,height=.0,width=.0];s->n1;n0[label="opened"];n1[label="closed"];n0->n1;n1->n0;}

	l.Printf("Current state: %s", s.CurrentName()) // Current state: closed

	if err := s.Transition(opened); err != nil {
		l.Fatal(err)
	}

	if err := s.Transition(opened); err != nil {
		l.Printf("Expected error: %s", err) // Expected error: invalid transition (opened -> opened)
	}

	if err := s.TransitionFrom(closed, opened); err != nil {
		l.Printf("Expected error: %s", err) // Expected error: transition failed (closed -> opened)
	}

	l.Printf("Current state: %s", s.CurrentName()) // Current state: opened

	if err := s.Transition(closed); err != nil {
		l.Fatal(err)
	}

	l.Printf("Current state: %s", s.CurrentName()) // Current state: closed
}
