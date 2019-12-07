// +build ignore

package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/Eyal-Shalev/lfsm"
)

var l = log.New(os.Stdout, "", log.Lshortfile)

type BankAccount struct {
	balance *big.Int
	state *lfsm.State
	name string
}

const (
	accountIdle uint32 = iota
	accountWithdraw
	accountDeposit
)

func (ba *BankAccount) Deposit(amount *big.Int) error {
	if err := ba.state.Transition(accountDeposit); err != nil {
		return err
	}
	defer ba.state.Transition(accountIdle)

	amount = big.NewInt(0).Set(amount) // clone amount.

	if amount.Sign() != 1 {
		return fmt.Errorf("cannot deposit non-positive amounts")
	}

	time.Sleep(time.Second / 2)
	ba.balance.Add(ba.balance, amount)
	l.Printf("[%s] Deposited %v (new balance %v)", ba.name, amount, ba.balance)
	return nil
}

func (ba *BankAccount) Withdraw(amount *big.Int) (func(), error) {
	if err := ba.state.Transition(accountWithdraw); err != nil {
		return nil, err
	}
	defer ba.state.Transition(accountIdle)

	amount = big.NewInt(0).Set(amount) // clone amount.

	if amount.Sign() != 1 {
		return nil, fmt.Errorf("cannot withdraw non-positive amounts")
	}

	if ba.balance.Cmp(amount) == -1 {
		return nil, fmt.Errorf("cannot withdraw more than the current balance")
	}

	time.Sleep(time.Second / 2)
	ba.balance.Sub(ba.balance, amount)
	l.Printf("[%s] Withdrawn %v (new balance %v)", ba.name, amount, ba.balance)

	return func() {
		for {
			if err := ba.state.Transition(accountDeposit); err == nil {
				ba.balance.Add(ba.balance, amount)
				l.Printf("[%s] Returned %v (new balance %v)", ba.name, amount, ba.balance)
				_ = ba.state.Transition(accountIdle)
				return
			}
			time.Sleep(time.Nanosecond)
		}
	}, nil
}

func NewBankAccount(balance *big.Int, name string) *BankAccount {
	return &BankAccount{
		name: name,
		balance: big.NewInt(0).Set(balance),
		state: lfsm.NewState(lfsm.Constraints{
			accountIdle: {accountWithdraw, accountDeposit},
			accountWithdraw: {accountIdle},
			accountDeposit: {accountIdle},
		}, lfsm.StateNameMap{
			accountIdle: "Idle",
			accountWithdraw: "Withdrawing",
			accountDeposit: "Depositing",
		}),
	}
}

type wireTransfer struct {
	state *lfsm.State
}

func (wt *wireTransfer) Transfer(from, to *BankAccount, amount *big.Int) error {
	if err := wt.state.Transition(transferAwaitFrom); err != nil {
		return err
	}
	defer wt.state.Transition(transferIdle)

	cancel, err := from.Withdraw(amount);
	if err != nil {
		return err
	}

	if err := to.Deposit(amount); err != nil {
		go cancel()
		return err
	}

	return nil
}

const (
	transferIdle uint32 = iota
	transferAwaitFrom
	transferAwaitTo
)

var Wire = wireTransfer{lfsm.NewState(
	lfsm.Constraints{
		transferIdle:      {transferAwaitFrom},
		transferAwaitFrom: {transferAwaitTo, transferIdle},
		transferAwaitTo:   {transferIdle},
	},
	lfsm.StateNameMap{
		transferIdle:      "Idle",
		transferAwaitFrom: "Await From",
		transferAwaitTo:   "Await To",
	},
)}

func main() {
	bank1 := NewBankAccount(big.NewInt(1000000), "bank1")
	bank2 := NewBankAccount(big.NewInt(1000000), "bank2")

	if err := Wire.Transfer(bank1, bank2, big.NewInt(1000001)); err != nil {
		l.Printf("Expected error: %s", err) // Expected error: cannot withdraw more than the current balance
	}

	if err := Wire.Transfer(bank1, bank2, big.NewInt(1000001)); err != nil {
		l.Printf("Expected error: %s", err) // Expected error: cannot withdraw more than the current balance
	}

	go func() {
		if err := Wire.Transfer(bank1, bank2, big.NewInt(500000)); err != nil {
			l.Fatalln(err)
		}
	}()

	if err := Wire.Transfer(bank2, bank1, big.NewInt(1500000)); err != nil {
		l.Printf("Expected error: %s", err) // Expected error: cannot withdraw more than the current balance
	}
	time.Sleep(time.Second + time.Millisecond)
	if err := Wire.Transfer(bank2, bank1, big.NewInt(1500000)); err != nil {
		l.Fatalln(err)
	}
}
