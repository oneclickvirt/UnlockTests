//go:build linux || darwin || freebsd || windows

package executor

import (
	rl "github.com/mattn/go-rl"
)

// readLine reads a line of input from the terminal with the given prompt.
// On linux/darwin/freebsd/windows, go-rl provides history and line-editing support.
func readLine(prompt string) (string, error) {
	return rl.ReadLine(prompt)
}
