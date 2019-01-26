// TODO
// Use reflection to automatically get the function name so they don't need to
// supply it.

// Package clean provides a simple API for global program cleanup when exiting.
package clean

import (
	"os"
	"sync"
)

var (
	// cleanup functions to execute with Do
	funcs = make(map[string]func())
	mu    sync.Mutex // protect funcs from concurrentcy
)

// Add a cleanup function to be executed
func Add(f func(), name string) {
	mu.Lock()
	defer mu.Unlock()

	funcs[name] = f
}

// Remove a cleanup function
func Remove(index string) {
	mu.Lock()
	defer mu.Unlock()

	delete(funcs, index)
}

// Do the cleanup, panics will be caught and ignored.
func Do() {
	mu.Lock()
	defer mu.Unlock()
	for _, f := range funcs {
		func() {
			defer func() { recover() }() // the show must go on

			f()
		}()
	}
}

// Exit wraps os.Exit except it executes cleanup before exiting.
func Exit(code int) {
	Do()
	os.Exit(code)
}
