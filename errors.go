package lfsm

import (
	"fmt"
)

type TransitionError transition

func (f TransitionError) Error() string {
	return fmt.Sprintf("trasition failed (%d -> %d)", f.src, f.dst)
}

type InvalidTransitionError TransitionError

func (f InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition (%d -> %d)", f.src, f.dst)
}