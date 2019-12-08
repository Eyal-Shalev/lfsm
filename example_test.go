package lfsm_test

import (
	"fmt"

	"github.com/Eyal-Shalev/lfsm"
)

func ExampleState_String() {
	fmt.Println(lfsm.NewState(lfsm.Constraints{0:{0}}, lfsm.StateNames{0:"foo"}))
	// Output: digraph g{s[label="",shape=none,height=.0,width=.0];s->n0;n0[label="foo",style=filled];n0->n0;}
}

func ExampleStateNames() {
	s := lfsm.NewState(lfsm.Constraints{0:{}}, lfsm.StateNames{0:"foo"})
	fmt.Printf("Current state: %s(%d).\n", s.CurrentName(), s.Current())
	// Output: Current state: foo(0).
}

func ExampleInitialState() {
	s := lfsm.NewState(lfsm.Constraints{0:{1},1:{0}}, lfsm.InitialState(1))
	fmt.Printf("Current state: %d.\n", s.Current())
	// Output: Current state: 1.
}

func ExampleState_CurrentName() {
	s := lfsm.NewState(lfsm.Constraints{0:{1},1:{0}}, lfsm.StateName(1, "foo"))
	fmt.Printf("Current state: %s(%d).\n", s.CurrentName(), s.Current())
	_ = s.Transition(1)
	fmt.Printf("Current state: %s(%d).\n", s.CurrentName(), s.Current())

	// Output:
	// Current state: 0(0).
	// Current state: foo(1).
}

func ExampleState_TransitionFrom() {
	s := lfsm.NewState(lfsm.Constraints{0:{1},1:{0}}, lfsm.InitialState(1))
	err := s.TransitionFrom(0,1)
	fmt.Printf("Expected error: %s.\n", err)

	err = s.TransitionFrom(1, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current state: %d.\n", s.Current())

	// Output:
	// Expected error: transition failed (0 -> 1) current state is not 0.
	// Current state: 0.
}

func ExampleState_Transition() {
	s := lfsm.NewState(lfsm.Constraints{0:{0,1},1:{0}}, lfsm.InitialState(1))
	err := s.Transition(1)
	fmt.Printf("Expected error: %s.\n", err)

	err = s.Transition(0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current state: %d.\n", s.Current())

	err = s.Transition(0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current state: %d.\n", s.Current())

	err = s.Transition(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current state: %d.\n", s.Current())

	// Output:
	// Expected error: invalid transition (1 -> 1).
	// Current state: 0.
	// Current state: 0.
	// Current state: 1.
}