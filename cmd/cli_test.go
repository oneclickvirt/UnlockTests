package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParseCLIStructuredOptions(t *testing.T) {
	opts, err := parseCLI([]string{"--structured", "-m", "6", "-f", "0", "-conc", "3", "--timeout", "1s"})
	if err != nil {
		t.Fatalf("parseCLI returned error: %v", err)
	}
	if !opts.jsonOutput || opts.mode != 6 || opts.selection != "0" || opts.concurrency != 3 || opts.timeout != time.Second {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestHelpRetainsLegacyFlags(t *testing.T) {
	var output bytes.Buffer
	newFlagSet(&cliOptions{}, &output).PrintDefaults()
	for _, legacy := range []string{"-b", "-f string", "-h", "-I string", "-L string", "-m int", "-s", "-test string", "-v"} {
		if !strings.Contains(output.String(), legacy) {
			t.Fatalf("help is missing legacy flag %q: %s", legacy, output.String())
		}
	}
}
