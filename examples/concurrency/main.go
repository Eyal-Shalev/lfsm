package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

func main() {
	l := log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
	wg := new(sync.WaitGroup)
	s := lfsm.NewState(lfsm.Constraints{
		0: {1, 2},
		1: {2},
		2: {2},
	}, lfsm.StateName(0, "start"), lfsm.StateName(1, "intermediate"), lfsm.StateName(2, "final"))

	if err := s.Transition(1); err != nil {
		l.Fatalln(err) // no reason to fail.
	}

	wg.Add(1)
	time.AfterFunc(time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
	})

	wg.Add(1)
	time.AfterFunc(2*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
	})

	wg.Add(1)
	time.AfterFunc(3*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(2); err != nil {
			l.Fatalln(err) // no reason to fail.
		}
	})

	if err := s.TransitionFrom(0, 2); err != nil {
		l.Printf("expected error: %s", err) // Invalid because current state is not 0
	}

	if err := s.Transition(1); err != nil {
		l.Printf("expected error: %s", err) // invalid because transitioning to 1 is only available from 0.
	}

	wg.Add(1)
	time.AfterFunc(4*time.Millisecond, func() {
		defer wg.Done()
		if err := s.Transition(1); err != nil {
			l.Printf("expected error: %s", err) // invalid because transitioning to 1 is only available from 0.
		}
	})

	wg.Wait()
}
