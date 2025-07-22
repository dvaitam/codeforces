package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed testsA.json
var testData []byte

type Tests struct {
	Inputs  []string `json:"inputs"`
	Outputs []string `json:"outputs"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	var t Tests
	if err := json.Unmarshal(testData, &t); err != nil {
		panic(err)
	}
	if len(t.Inputs) != len(t.Outputs) {
		panic("test data mismatch")
	}
	exe := os.Args[1]
	for i, in := range t.Inputs {
		expect := strings.TrimSpace(t.Outputs[i])
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(t.Inputs))
}
