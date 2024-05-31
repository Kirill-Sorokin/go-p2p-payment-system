package main

import "time"

type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	Timestamp time.Time
}
