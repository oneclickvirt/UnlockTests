//go:build !(linux || darwin || freebsd || windows)

package executor

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readLine reads a line of input from stdin using a basic bufio reader.
// This fallback is used on platforms where go-rl is not supported (e.g. openbsd, netbsd).
func readLine(prompt string) (string, error) {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}
