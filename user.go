package main

import (
	"sync"
)

type User struct {
	name    string
	balance int
	mu      sync.Mutex
}

func NewUser(name string) *User {
	return &User{name: name, balance: 0}
}

func (u *User) Deposit(amount int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.balance += amount
}

func (u *User) Withdraw(amount int) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.balance >= amount {
		u.balance -= amount
		return true
	}
	return false
}

func (u *User) GetBalance() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.balance
}

func (u *User) Send(receiver *User, amount int) bool {
	if u.Withdraw(amount) {
		receiver.Deposit(amount)
		return true
	}
	return false
}
