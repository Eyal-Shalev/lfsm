package lfsm_test

import (
	"log"
	"testing"

	"github.com/Eyal-Shalev/lfsm"
)

func TestExample1(t *testing.T) {

	const (
		closed uint64 = iota
		opened
	)
	s := lfsm.NewState(lfsm.Constraints{
		closed: {opened},
		opened: {closed},
	}, lfsm.InitialState(closed))

	log.Println(s.Current()) // 2019/12/06 17:34:52 0

	if err := s.Transition(opened); err != nil {
		log.Fatal(err)
	}

	log.Println(s.Current()) // 2019/12/06 17:34:52 1

	if err := s.Transition(closed); err != nil {
		log.Fatal(err)
	}

	log.Println(s.Current()) // 2019/12/06 17:34:52 0
}
