package main

import (
	"os"
	"testing"
)

func TestMainHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ut", "-h"}
	main()
}

func TestMainRejectsInvalidMode(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ut", "-m", "5", "-f", "0"}
	main()
}

func TestMainRejectsInvalidLanguage(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ut", "-L", "ja", "-f", "0"}
	main()
}
