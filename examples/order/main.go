// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

const (
	creating uint64 = iota
	adding
	finalizing
	paying
	paid
	processing
	shipped
	delivered
	canceled
)

func stateString(s uint64) string {
	switch s {
	case creating:
		return fmt.Sprintf("creating(%d)", creating)
	case finalizing:
		return fmt.Sprintf("finalizing(%d)", finalizing)
	case adding:
		return fmt.Sprintf("adding(%d)", adding)
	case paying:
		return fmt.Sprintf("paying(%d)", paying)
	case paid:
		return fmt.Sprintf("paid(%d)", paid)
	case processing:
		return fmt.Sprintf("processing(%d)", processing)
	case shipped:
		return fmt.Sprintf("shipped(%d)", shipped)
	case delivered:
		return fmt.Sprintf("delivered(%d)", delivered)
	case canceled:
		return fmt.Sprintf("canceled(%d)", canceled)
	}
	return fmt.Sprintf("unknown(%d)", s)
}

func printStates() {
	fmt.Println("States:")
	for _, s := range []uint64{creating, finalizing, adding, paid, processing, shipped, delivered, canceled} {
		fmt.Println("\t", stateString(s))
	}
	fmt.Println()
}

type order struct {
	state *lfsm.State

	items map[string]int
	mu    sync.Mutex
}

func (o *order) String() string {
	itemsStr, err := json.Marshal(o.items)
	fatalIfErr(err)
	return fmt.Sprintf("order{items: %s, state: %s}", itemsStr, stateString(o.state.Current()))
}

func (o *order) incItem(item string, count int) error {

	if err := o.state.Transition(adding); err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	if _, ok := o.items[item]; !ok {
		o.items[item] = 0
	}
	o.items[item] += count
	if o.items[item] < 0 {
		log.Printf("Invalid item count.")
		return o.state.Transition(canceled)
	}

	return o.state.Transition(creating)
}

func (o *order) checkout() error {
	return o.state.Transition(finalizing)
}

func (o *order) pay() error {
	err := o.state.Transition(paying)
	if err != nil {
		return err
	}
	if rand.Intn(2) == 1 {
		log.Println("payment failed")
		return o.state.Transition(finalizing)
	}
	time.AfterFunc(time.Second, func() {
		logErr(o.state.Transition(processing))
	})
	time.AfterFunc(time.Second * 2, func() {
		if rand.Intn(2) == 1 {
			log.Println("processing failed")
			logErr(o.state.Transition(canceled))
			return
		}
		logErr(o.state.Transition(shipped))
	})
	return o.state.Transition(paid)
}

func (o *order) cancel() error {
	return o.state.Transition(canceled)
}

var transitionMap = lfsm.Constraints{
	creating:   {adding, finalizing, canceled},
	adding:     {creating, canceled},
	finalizing: {adding, paying, canceled},
	paying:     {paid, finalizing},
	paid:       {processing, canceled},
	processing: {shipped, canceled},
	shipped:    {delivered},
	canceled:   {canceled},
}

func newOrder() *order {
	return &order{
		state: lfsm.NewState(transitionMap, lfsm.InitialState(creating)),
		items: map[string]int{},
	}
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func logErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	printStates()
	o := newOrder()

	// wg := new(sync.WaitGroup)
	//
	// fmt.Println("A", o.String())
	//
	// randFn(wg, func() {
	// 	fmt.Println("B", 0, o.String())
	// 	logErr(o.incItem("foo", 3))
	// 	fmt.Println("B", 1, o.String())
	// })
	//
	// randFn(wg, func() {
	// 	fmt.Println("C", 0, o.String())
	// 	logErr(o.incItem("bar", 2))
	// 	fmt.Println("C", 1, o.String())
	// })
	//
	// randFn(wg, func() {
	// 	fmt.Println("D", 0, o.String())
	// 	logErr(o.incItem("foo", 1))
	// 	fmt.Println("D", 1, o.String())
	// })
	//
	// randFn(wg, func() {
	// 	fmt.Println("E", 0, o.String())
	// 	logErr(o.checkout())
	// 	fmt.Println("E", 1, o.String())
	// })
	//
	// wg.Wait()
	//
	// fmt.Println("F", 0, o.String())
	// logErr(o.checkout())
	// fmt.Println("G", 1, o.String())
	// logErr(o.pay())
	// fmt.Println("H", 1, o.String())
	// time.Sleep(time.Second * 3)
	// fmt.Println("I", 1, o.String())

	wg := new(sync.WaitGroup)
	for i:=0;i<999999;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

			if err := o.state.Transition(canceled); err != nil {
				log.Println(err)
				return
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

			logErr(o.state.Transition(canceled))
		}()
	}

	wg.Wait()
}
