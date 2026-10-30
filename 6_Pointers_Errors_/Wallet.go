package main

import (
	"errors"
	"fmt"
)

var ErrInsufficientFunds = errors.New("Cannot withdraw, insufficient funds.")

type Bitcoin int

type Stringer interface {
	String() string // Stringer helps me provide output the way I want upon calling print on a custom type instead of the default.
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

type Wallet struct {
	name    string
	balance Bitcoin /*Capitalisation means variables can be accessed from outside which is why balance variable is lowercase
	but Deposit and Balance methods are upper case */
}

func (w *Wallet) Deposit(amount Bitcoin) { // Had to use a pointer here because otherwise balance of copy was being updated, not the struct.
	if amount > 0 {
		w.balance += amount /* Actually now w is not a wallet but basically a pointer to the address of a wallet. in fact w
		holds and address. so to get any value stored at w we need to dereference. and the correct way will be (*w).balance
		as that implies the wallet at address w ka balance. but this notation is cumbersome and go allows us to write wrongly
		w.balance.*/
	}
}

// note, how for a user only deposit and balance buttons will be availableand he cant directly access balance.

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return ErrInsufficientFunds
	}
	w.balance -= amount
	return nil
}
