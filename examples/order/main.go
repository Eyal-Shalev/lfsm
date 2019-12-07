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

type order struct {
	state *lfsm.State

	items map[string]int
	mu    sync.Mutex
}

func (o *order) String() string {
	itemsStr, err := json.Marshal(o.items)
	fatalIfErr(err)
	return fmt.Sprintf("order{items: %s, state: %s}", itemsStr, o.state.CurrentName())
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
	time.AfterFunc(time.Second*2, func() {
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

func newOrder() *order {
	return &order{
		state: lfsm.NewState(
			lfsm.Constraints{
				creating:   {adding, finalizing, canceled},
				adding:     {creating, canceled},
				finalizing: {adding, paying, canceled},
				paying:     {paid, finalizing},
				paid:       {processing, canceled},
				processing: {shipped, canceled},
				shipped:    {delivered},
				canceled:   {canceled},
			},
			lfsm.InitialState(creating),
			lfsm.StateNames(map[uint64]string{
				creating: "creating",
				adding: "adding",
				finalizing: "finalizing",
				paying: "paying",
				paid: "paid",
				processing: "processing",
				shipped: "shipped",
				delivered: "delivered",
				canceled: "canceled",
			}),
		),
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
	o := newOrder()
	fmt.Println(o.state)

	fatalIfErr(o.incItem("foo", 3))
	fatalIfErr(o.checkout())
	fatalIfErr(o.incItem("foo", 1))

}
