package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource2095B = "./2095B.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	// Sanity: ensure reference solution itself passes validation.
	for idx, tc := range tests {
		out, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		val, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if !validColumn(val) {
			fmt.Fprintf(os.Stderr, "reference output fails validity on test %d (%s): column %d not winning\ninput:\n%soutput:\n%s",
				idx+1, tc.name, val, tc.input, out)
			os.Exit(1)
		}
	}

	for idx, tc := range tests {
		out, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		val, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if !validColumn(val) {
			fmt.Fprintf(os.Stderr, "candidate loses on test %d (%s): column %d does not guarantee a win\ninput:\n%soutput:\n%s",
				idx+1, tc.name, val, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2095B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2095B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2095B)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.Write(errBuf.Bytes())
	}
	return out.String(), err
}

func parseOutput(output string) (int, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	v, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("cannot parse integer: %v", err)
	}
	return v, nil
}

func validColumn(c int) bool {
	// Any column strictly outside [0,16] guarantees a win (problem trick).
	return c < 0 || c > 16
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 10)
	for g := 1; g <= 10; g++ {
		tests = append(tests, testCase{
			name:  fmt.Sprintf("game_%d", g),
			input: fmt.Sprintf("Game %d\n", g),
		})
	}
	return tests
}
