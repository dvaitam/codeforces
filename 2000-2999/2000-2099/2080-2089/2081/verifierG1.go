package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2080-2089/2081/2081G1.go"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for idx, n := range tests {
		input := fmt.Sprintf("%d\n", n)

		refOut, err := runWithInput(exec.Command(refBin), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (n=%d): %v\noutput:\n%s\n", idx+1, n, err, refOut)
			os.Exit(1)
		}
		candOut, err := runWithInput(commandFor(candidate), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (n=%d): %v\noutput:\n%s\n", idx+1, n, err, candOut)
			os.Exit(1)
		}

		expVal, err := parseSingle(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		gotVal, err := parseSingle(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\ninput:\n%s", idx+1, err, candOut, input)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (n=%d): expected %d got %d\n", idx+1, n, expVal, gotVal)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2081G1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseSingle(output string) (uint64, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no output")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %v", err)
	}
	return val, nil
}

func generateTests() []int {
	tests := []int{
		5,      // sample 1
		10,     // small
		1234,   // medium
		99991,  // prime-ish
		200000, // upper moderate
		1,      // minimal
		2,      // tiny
		3,
		4,
		6,
		12,
	}

	rng := rand.New(rand.NewSource(2081))
	const maxN = 800000
	for len(tests) < 60 {
		n := rng.Intn(maxN) + 1
		tests = append(tests, n)
	}
	return tests
}
