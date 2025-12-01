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

const refSourceF = "./533F.go"

type result struct {
	count     int
	positions []int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierF.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
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
	expected, err := parseOutput(refOut)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	candOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	actual, err := parseOutput(candOut)
	if err != nil {
		fail("failed to parse candidate output: %v\n%s", err, candOut)
	}

	if err := compareResults(expected, actual); err != nil {
		fail("%v", err)
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "533F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceF))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runCommand(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
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

func parseOutput(out string) (result, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return result{}, fmt.Errorf("empty output")
	}
	cnt, err := strconv.Atoi(fields[0])
	if err != nil {
		return result{}, fmt.Errorf("invalid count %q", fields[0])
	}
	numFields := len(fields) - 1
	if cnt != numFields {
		return result{}, fmt.Errorf("declared %d positions but found %d numbers", cnt, numFields)
	}
	pos := make([]int, numFields)
	prev := 0
	for i := 0; i < numFields; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return result{}, fmt.Errorf("invalid position %q", fields[i+1])
		}
		if i > 0 && val <= prev {
			return result{}, fmt.Errorf("positions not strictly increasing")
		}
		prev = val
		pos[i] = val
	}
	return result{count: cnt, positions: pos}, nil
}

func compareResults(expected, actual result) error {
	if expected.count != actual.count {
		return fmt.Errorf("wrong number of positions: expected %d got %d", expected.count, actual.count)
	}
	if len(expected.positions) != len(actual.positions) {
		return fmt.Errorf("position list length mismatch: expected %d got %d", len(expected.positions), len(actual.positions))
	}
	for i, exp := range expected.positions {
		if actual.positions[i] != exp {
			return fmt.Errorf("positions differ at index %d: expected %d got %d", i+1, exp, actual.positions[i])
		}
	}
	return nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
