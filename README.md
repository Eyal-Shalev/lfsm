# Lock-Free State Machine [![GoDoc](https://godoc.org/github.com/Eyal-Shalev/lfsm?status.svg)](https://godoc.org/github.com/Eyal-Shalev/lfsm) [![Build Status](https://travis-ci.org/Eyal-Shalev/lfsm.svg?branch=master)](https://travis-ci.org/Eyal-Shalev/lfsm)
LFSM is a light-weight State Machine implementation that doesn't use any locks.

## Basic Example
```go
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
		lfsm.StateName(opened, "opened"),
		lfsm.StateName(closed, "closed"),
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
```
