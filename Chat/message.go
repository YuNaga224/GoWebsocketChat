package main

import (
	"time"
)

//messageは一つのメッセージを表す

type message struct {
	Name    string
	Message string
	When    time.Time
}
