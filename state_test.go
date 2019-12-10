package lfsm_test

import (
	"testing"

	"github.com/Eyal-Shalev/lfsm"
)

func fatalIfErr(tb testing.TB, err error) {
	if err != nil {
		tb.Fatal(err)
	}
}
func logErr(tb testing.TB, err error) {
	if err != nil {
		tb.Error(err)
	}
}

func benchBigState(s *lfsm.State, b *testing.B) {
	size := s.Current() + 1
	for n := 0; n < b.N; n++ {
		for i := uint32(0); i < size; i++ {
			logErr(b, s.Transition(i))
		}
	}
}

func newBigState(size uint32) *lfsm.State {
	constraints := make(lfsm.Constraints, size)
	for i := uint32(0); i < size; i++ {
		constraints[i] = []uint32{(i + 1) % size}
	}
	return lfsm.NewState(constraints, lfsm.InitialState(size-1))
}

func BenchmarkState10(b *testing.B)   { benchBigState(newBigState(10), b) }
func BenchmarkState100(b *testing.B)  { benchBigState(newBigState(100), b) }
func BenchmarkState1000(b *testing.B) { benchBigState(newBigState(1000), b) }

func TestIntermediateState(t *testing.T) {
	s := lfsm.NewState(lfsm.Constraints{
		0: {1},
		1: {0, 2},
		2: {0},
	}, lfsm.InitialState(0))

	if err := s.Transition(0); err == nil {
		t.Error("Invalid transition error expected.")
	}
	if err := s.Transition(2); err == nil {
		t.Error("Invalid transition error expected.")
	}

	fatalIfErr(t, s.Transition(1))
	if err := s.Transition(1); err == nil {
		t.Error("Invalid transition error expected.")
	}

	fatalIfErr(t, s.Transition(0))

	if err := s.Transition(2); err == nil {
		t.Error("Invalid transition error expected.")
	}

	fatalIfErr(t, s.Transition(1))
	fatalIfErr(t, s.Transition(2))

	if err := s.Transition(1); err == nil {
		t.Error("Invalid transition error expected.")
	}
	if err := s.Transition(2); err == nil {
		t.Error("Invalid transition error expected.")
	}
}
