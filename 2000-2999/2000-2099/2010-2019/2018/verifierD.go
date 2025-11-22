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

const refSourceD = "2000-2999/2000-2099/2010-2019/2018/2018D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
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
	expected, err := parseInts(refOut)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	actual, err := parseInts(userOut)
	if err != nil {
		fail("failed to parse candidate output: %v\n%s", err, userOut)
	}

	if len(actual) != len(expected) {
		fail("wrong number of outputs: expected %d got %d", len(expected), len(actual))
	}
	for i := range expected {
		if actual[i] != expected[i] {
			fail("outputs differ at line %d: expected %d got %d", i+1, expected[i], actual[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2018D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceD))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
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

func parseInts(out string) ([]int64, error) {
	fields := strings.Fields(out)
	res := make([]int64, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
