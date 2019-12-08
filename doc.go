/*
Package lfsm provides a light-weight lock-free state machine.

This state machine uses atomic operations to transition between states.

Example usage:
	package main

	import (
		"log"

		"github.com/Eyal-Shalev/lfsm"
	)

	func main() {
		s := lfsm.NewState(
			lfsm.Constraints{
				0: {1},
				1: {0},
			},
			lfsm.InitialState(1),
			lfsm.StateNames{0: "opened", 1: "closed"},
		)

		log.Printf("Current state: %s", s.CurrentName()) // Current state: opened
		if err := s.Transition(0); err != nil {
			panic(err)
		}
		log.Printf("Current state: %s", s.CurrentName()) // Current state: opened
	}
*/
package lfsm
