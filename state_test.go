package lfsm_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	lfsm "state"
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
		for i := uint64(0); i < size; i++ {
			logErr(b, s.Transition(i))
		}
	}
}

func newBigState(size uint64) *lfsm.State {
	constraints := make(lfsm.Constraints, size)
	for i := uint64(0); i < size; i++ {
		constraints[i] = []uint64{(i + 1) % size}
	}
	return lfsm.NewState(size-1, constraints)
}

func BenchmarkState10(b *testing.B)   { benchBigState(newBigState(10), b) }
func BenchmarkState100(b *testing.B)  { benchBigState(newBigState(100), b) }
func BenchmarkState1000(b *testing.B) { benchBigState(newBigState(1000), b) }

func TestIntermediateState(t *testing.T) {
	s := lfsm.NewState(0, lfsm.Constraints{
		0: {1},
		1: {0, 2},
		2: {0},
	})

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

func TestFSM_Transition(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	s := lfsm.NewState(0, lfsm.Constraints{
		0: {1, 2},
		1: {0, 2},
		2: {0, 1},
	})

	var tErr [2]atomic.Value

	wg := new(sync.WaitGroup)
	for i := 0; i < 999; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			err := s.Transition(uint64(rand.Intn(3)))
			switch err.(type) {
			case *lfsm.InvalidTransitionError:
				tErr[1].Store(err)
			case *lfsm.TransitionError:
				tErr[0].Store(err)
			}
		}()
	}

	wg.Wait()
	if tErr[0].Load() == nil {
		t.Error("transition error expected")
	}
	if tErr[1].Load() == nil {
		t.Error("invalid transition error expected")
	}
}
