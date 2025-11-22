package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceE = "2000-2999/2100-2199/2100-2109/2106/2106E.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	totalAns, err := countAnswers(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runCommand(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	expected, err := parseInts(refOut, totalAns)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	actual, err := parseInts(userOut, totalAns)
	if err != nil {
		fail("failed to parse candidate output: %v\n%s", err, userOut)
	}

	for i := 0; i < totalAns; i++ {
		if expected[i] != actual[i] {
			fail("answer mismatch at position %d: expected %d got %d", i+1, expected[i], actual[i])
		}
	}

	fmt.Println("OK")
}

func countAnswers(input []byte) (int, error) {
	fields := strings.Fields(string(input))
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty input")
	}
	idx := 0
	t, err := strconv.Atoi(fields[idx])
	if err != nil {
		return 0, fmt.Errorf("invalid t: %v", err)
	}
	idx++
	total := 0
	for tc := 0; tc < t; tc++ {
		if idx+1 >= len(fields) {
			return 0, fmt.Errorf("input ended prematurely at test %d", tc+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return 0, fmt.Errorf("invalid n at test %d: %v", tc+1, err)
		}
		q, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return 0, fmt.Errorf("invalid q at test %d: %v", tc+1, err)
		}
		idx += 2
		need := n
		if idx+need > len(fields) {
			return 0, fmt.Errorf("input ended while reading permutation at test %d", tc+1)
		}
		idx += need // permutation
		need = 3 * q
		if idx+need > len(fields) {
			return 0, fmt.Errorf("input ended while reading queries at test %d", tc+1)
		}
		idx += need
		total += q
	}
	if idx != len(fields) {
		return 0, fmt.Errorf("unexpected extra tokens in input")
	}
	return total, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2106E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceE))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
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

func runCommand(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseInts(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", f, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
