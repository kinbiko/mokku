package main

import (
	"github.com/atotto/clipboard"
	"github.com/kinbiko/mokku"
)

func main() {
	s, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}
	mock, err := mokku.Mock([]byte(s))
	if err != nil {
		panic(err)
	}
	clipboard.WriteAll(string(mock))
}
