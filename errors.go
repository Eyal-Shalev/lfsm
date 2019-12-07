package lfsm

import (
	"fmt"
)

type TransitionError transition

func (f TransitionError) Error() string {
	return fmt.Sprintf("transition failed (%s -> %s)", f.stateNames.find(f.src), f.stateNames.find(f.dst))
}

type InvalidTransitionError TransitionError

func (f InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition (%s -> %s)", f.stateNames.find(f.src), f.stateNames.find(f.dst))
}