package lfsm

import (
	"fmt"
)

// TransitionError is an error struct for all failed transition attempts.
type TransitionError struct {
	Src, Dst   uint32
	stateNames StateNames
	msg        string
}
func (f *TransitionError) SrcName() string {
	return f.stateNames.find(f.Src)
}
func (f *TransitionError) DstName() string {
	return f.stateNames.find(f.Src)
}
func (f *TransitionError) Error() string {
	if f.msg != "" {
		return f.msg
	}
	return fmt.Sprintf("transition failed (%s -> %s)", f.SrcName(), f.DstName())
}

// NewFailedTransitionError reports that the current state differs from the transition source state.
func NewFailedTransitionError(src, dst uint32, stateNames StateNames) *TransitionError {
	return &TransitionError{
		src,
		dst,
		stateNames,
		fmt.Sprintf("transition failed (%s -> %s) current state is not %s", stateNames.find(src), stateNames.find(dst), stateNames.find(src)),
	}
}

// NewFailedTransitionError reports about a transition attempt that was not defined in the Constraints map.
func NewInvalidTransitionError(src, dst uint32, stateNames StateNames) *TransitionError {
	return &TransitionError{
		src,
		dst,
		stateNames,
		fmt.Sprintf("invalid transition (%s -> %s)", stateNames.find(src), stateNames.find(dst)),
	}
}