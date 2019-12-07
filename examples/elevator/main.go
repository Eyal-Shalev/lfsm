package main

import (
	"github.com/Eyal-Shalev/lfsm"
)

type Elevator struct {
	state *lfsm.State
}

func (el Elevator) goTo(dst uint64) {
	el.state.Transition(dst)
}

func makeFloors(size uint64) lfsm.Constraints {
	floors := make(lfsm.Constraints, size)
	for i := uint64(0); i < size; i++ {
		floors[i] = make([]uint64, 0, size-1)
		for j := uint64(0); j < size; j++ {
			if j == i {
				continue
			}
			floors[i] = append(floors[i], j)
		}
	}
	return floors
}

func main() {
	el := Elevator{
		state: lfsm.NewState(0, makeFloors(20)),
	}
}
