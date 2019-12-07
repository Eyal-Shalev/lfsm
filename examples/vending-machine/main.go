// +build ignore

package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

const (
	idle uint64 = iota
	addingCash
	buying
	withdrawing
)

type VendingMachine struct {
	state    *lfsm.State
	balance  uint64
	products map[string]uint64
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (vm *VendingMachine) addCash(amount uint64) error {
	err := vm.state.Transition(addingCash)
	if err != nil {
		return fmt.Errorf("cannot add %v$ because %w", amount/100, err)
	}
	atomic.AddUint64(&vm.balance, amount)
	time.Sleep(time.Second / 2)
	log.Printf("balance: %v$", float64(atomic.LoadUint64(&vm.balance))/100)
	return vm.state.Transition(idle)
}

func (vm *VendingMachine) withdraw() error {
	err := vm.state.Transition(withdrawing)
	if err != nil {
		return fmt.Errorf("cannot withdraw because %w", err)
	}
	withdrawn := atomic.SwapUint64(&vm.balance, 0)
	log.Printf("withdrawn %v$, balance: 0$", float64(withdrawn)/100)
	return vm.state.Transition(idle)
}

func (vm *VendingMachine) buy(product string) error {
	err := vm.state.Transition(buying)
	if err != nil {
		panicIfErr(vm.state.Transition(idle))
		return fmt.Errorf("cannot buy because %w", err)
	}
	balance := atomic.LoadUint64(&vm.balance)
	if balance < vm.products[product] {
		panicIfErr(vm.state.Transition(idle))
		return fmt.Errorf("cannot buy %s - missing %v$", product, float64(vm.products[product]-balance)/100)
	}
	atomic.AddUint64(&vm.balance, -vm.products[product])
	balance = atomic.LoadUint64(&vm.balance)
	vm.dropItem()
	log.Printf("%s was purchased, %v$ is left in the machine", product, float64(balance)/100)
	return vm.state.Transition(idle)
}

func (vm *VendingMachine) dropItem() {
	time.Sleep(time.Second)
}

func logIfErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	vm := VendingMachine{
		state: lfsm.NewState(
			lfsm.Constraints{
				idle:        {addingCash, buying, withdrawing},
				addingCash:  {idle},
				buying:      {idle},
				withdrawing: {idle},
			},
			lfsm.InitialState(idle),
			lfsm.StateNames(map[uint64]string{
				idle: "Idle",
				addingCash: "Adding cash",
				buying: "Buying",
				withdrawing: "Withdrawing",
			}),
		),
		products: map[string]uint64{
			"coke cola":    200,
			"pepsi cola":   150,
			"orange juice": 250,
		},
	}

	logIfErr(vm.buy("coke cola"))

	fatalIfErr(vm.addCash(250))

	time.AfterFunc(time.Second/2, func() {
		logIfErr(vm.withdraw())
	})

	fatalIfErr(vm.buy("pepsi cola"))

	fatalIfErr(vm.withdraw())
}
