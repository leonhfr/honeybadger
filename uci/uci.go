// Package uci implements the Universal Chess Interface.
package uci

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/leonhfr/honeybadger/engine"
)

// Run runs the program in UCI mode.
//
// Run parses command from the reader, executes them with the provided engine
// and writes the responses on the writer.
func Run(e *engine.Engine, r io.Reader, w io.Writer) {
	cc := make(chan Command)
	rc := make(chan Response)

	var wg sync.WaitGroup
	wg.Add(3)

	go parser(&wg, r, cc)
	go worker(&wg, e, cc, rc)
	go logger(&wg, w, rc)

	wg.Wait()
}

// parser reads from the reader and sends commands to the channel
func parser(wg *sync.WaitGroup, r io.Reader, cc chan<- Command) {
	defer wg.Done()
	defer close(cc)

	for scanner := bufio.NewScanner(r); scanner.Scan(); {
		command := strings.Fields(scanner.Text())
		c := Parse(command)
		cc <- c
		if _, ok := c.(CommandQuit); ok {
			return
		}
	}
}

// worker receives commands from the channel and execute them
func worker(wg *sync.WaitGroup, e *engine.Engine, cc <-chan Command, rc chan<- Response) {
	defer wg.Done()
	defer close(rc)

	for c := range cc {
		c.Run(e, rc)
	}
}

// logger receives responses from the channel and logs them
func logger(wg *sync.WaitGroup, w io.Writer, rc <-chan Response) {
	defer wg.Done()

	l := log.New(w, "", 0)
	for r := range rc {
		l.Println(r.String())
	}
}
