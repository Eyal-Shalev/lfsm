package main

import (
	"log"
	"os"

	"github.com/Eyal-Shalev/lfsm"
)

func main() {
	l := log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
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

	l.Println(s.CurrentName()) // 03:48:26.109387 example_test.go:29: closed

	if err := s.Transition(opened); err != nil {
		l.Fatal(err)
	}

	l.Println(s.CurrentName()) // 03:48:26.132388 example_test.go:35: opened

	if err := s.Transition(closed); err != nil {
		l.Fatal(err)
	}

	l.Println(s.CurrentName()) // 03:48:26.132388 example_test.go:41: closed
}
