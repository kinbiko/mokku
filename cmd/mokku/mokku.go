package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/kinbiko/mokku"
)

func main() {
	s, err := clipboard.ReadAll()
	if err != nil {
		errorOut(err)
	}

	mock, err := mokku.Mock([]byte(s))
	if err != nil {
		errorOut(err)
	}

	if err = clipboard.WriteAll(string(mock)); err != nil {
		errorOut(err)
	}
}

func errorOut(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
