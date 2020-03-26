package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/kinbiko/mokku"
)

const usage = `Usage:
1. Copy the interface you want to mock
2. Run 'mokku'
3. Paste the mocked implementation that has been written to your clipboard`

func main() {
	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

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
	flag.Usage()
	os.Exit(1)
}
