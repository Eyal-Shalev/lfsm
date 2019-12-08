package lfsm

import (
	"fmt"
)

// FailedTransitionError reports about a transition that failed because the current state differed from the source
// state.
type FailedTransitionError transition

func (f FailedTransitionError) Error() string {
	return fmt.Sprintf("transition failed (%s -> %s) current state is not %s", f.stateNames.find(f.src), f.stateNames.find(f.dst), f.stateNames.find(f.src))
}

// InvalidTransitionError reports about a transition attempt that was not defined in the Constraints map.
type InvalidTransitionError transition

func (f InvalidTransitionError) Error() string {
	return fmt.Sprintf("invalid transition (%s -> %s)", f.stateNames.find(f.src), f.stateNames.find(f.dst))
}