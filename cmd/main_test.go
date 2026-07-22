package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ut", "-h"}
	main()
}

func TestMainRejectsInvalidMode(t *testing.T) {
	expectMainExit(t, []string{"-m", "5", "-f", "0"}, 2)
}

func TestMainRejectsInvalidLanguage(t *testing.T) {
	expectMainExit(t, []string{"-L", "ja", "-f", "0"}, 2)
}

func TestParseCLIRejectsIgnoredConflictingAndPositionalOptions(t *testing.T) {
	for _, args := range [][]string{
		{"-json", "-s"},
		{"-json", "-b"},
		{"-timeout", "1s"},
		{"-f", "1", "-test", "Netflix"},
		{"unexpected"},
	} {
		if _, err := parseCLI(args); err == nil {
			t.Fatalf("expected arguments %v to be rejected", args)
		}
	}
}

func TestParseCLINormalizesLanguageAndPreservesStructuredOptions(t *testing.T) {
	opts, err := parseCLI([]string{"-L", " EN ", "-f", "0"})
	if err != nil {
		t.Fatal(err)
	}
	if opts.language != "en" {
		t.Fatalf("unexpected options: %#v", opts)
	}
	structured, err := parseCLI([]string{"-json", "-timeout", "5s", "-conc", "4"})
	if err != nil || !structured.jsonOutput || structured.concurrency != 4 {
		t.Fatalf("unexpected structured options: %#v err=%v", structured, err)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	separator := 0
	for index, arg := range os.Args {
		if arg == "--" {
			separator = index + 1
			break
		}
	}
	os.Args = append([]string{"ut"}, os.Args[separator:]...)
	main()
	os.Exit(0)
}

func expectMainExit(t *testing.T, args []string, want int) {
	t.Helper()
	commandArgs := append([]string{"-test.run=TestHelperProcess", "--"}, args...)
	command := exec.Command(os.Args[0], commandArgs...)
	command.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	err := command.Run()
	if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != want {
		t.Fatalf("main(%s) exit=%v, want %d", strings.Join(args, " "), err, want)
	}
}
