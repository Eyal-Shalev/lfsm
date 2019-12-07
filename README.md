# Lock-Free State Machine [![Build Status](https://travis-ci.org/Eyal-Shalev/lfsm.svg?branch=master)](https://travis-ci.org/Eyal-Shalev/lfsm)
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
		opened uint64 = iota
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

## Concurrent access
```go
package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

func main() {
	l := log.New(os.Stdout, "", log.Lshortfile)
	wg := new(sync.WaitGroup)
	s := lfsm.NewState(lfsm.Constraints{
		0: {1, 2},
		1: {2},
		2: {2},
	}, lfsm.StateName(0, "start"), lfsm.StateName(1, "intermediate"), lfsm.StateName(2, "final"))

	l.Printf("Current state: %s", s.CurrentName()) // Current state: start

	if err := s.Transition(1); err != nil {
		l.Fatalln(err) // no reason to fail.
	}
	l.Printf("Current state: %s", s.CurrentName()) // Current state: intermediate

	wg.Add(1)
	time.AfterFunc(time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
		l.Printf("Current state: %s", s.CurrentName()) // Current state: final
	})

	wg.Add(1)
	time.AfterFunc(2*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
		l.Printf("Current state: %s", s.CurrentName()) // Current state: final
	})

	wg.Add(1)
	time.AfterFunc(3*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
		l.Printf("Current state: %s", s.CurrentName()) // Current state: final
	})

	if err := s.TransitionFrom(0, 2); err != nil {
		l.Printf("expected error: %s", err) // expected error: transition failed (start -> final)
	}

	if err := s.Transition(1); err != nil {
		l.Printf("expected error: %s", err) // expected error: transition failed (intermediate -> intermediate)
	}

	wg.Add(1)
	time.AfterFunc(4*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(0); err != nil {
			l.Printf("expected error: %s", err) // expected error: transition failed (final -> start)
		}
	})

	wg.Wait()
}
```