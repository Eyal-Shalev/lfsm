/*
Package lfsm provides a light-weight lock-free state machine.

This state machine uses atomic operations to transition between states.

Basic Example:
	package main

	import (
		"log"
		"github.com/Eyal-Shalev/lfsm"
	)

	func main() {
		s := lfsm.NewState(lfsm.Constraints{0: {1}, 1: {0}})

		log.Printf("Current state: %s", s.Current()) // Current state: 1
		if err := s.Transition(0); err != nil {
			panic(err)
		}
		log.Printf("Current state: %s", s.Current()) // Current state: 0
	}

You may want to label your states, so the lfsm.StateNames struct and lfsm.StateName function can be used to supply
options for the State Machine. For extra convenience you can use constants for your states.
Named Example:
	package main

	import (
		"log"
		"github.com/Eyal-Shalev/lfsm"
	)
	const (
		opened uint64 = iota
		closed
	)

	func main() {
		s := lfsm.NewState(
			lfsm.Constraints{
				opened: {closed},
				closed: {opened},
			},
			lfsm.StateNames{opened: "opened", closed: "closed"},
		)

		log.Printf("Current state: %s", s.Current()) // Current state: closed
		if err := s.Transition(0); err != nil {
			panic(err)
		}
		log.Printf("Current state: %s", s.Current()) // Current state: opened
	}
*/
package lfsm
