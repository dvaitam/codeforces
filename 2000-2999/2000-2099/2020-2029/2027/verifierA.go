package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2027/2027A.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	if len(input) == 0 {
		fail("empty input")
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAns, err := parseOutput(refOut, input)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAns, err := parseOutput(candOut, input)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if len(refAns) != len(candAns) {
		fail("number of answers mismatch: expected %d got %d", len(refAns), len(candAns))
	}

	for i := range refAns {
		if refAns[i] != candAns[i] {
			fail("test %d: expected %d got %d", i+1, refAns[i], candAns[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2027A-ref-*")
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

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, input []byte) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))

	var t int64
	inputReader := bufio.NewReader(bytes.NewReader(input))
	if _, err := fmt.Fscan(inputReader, &t); err != nil {
		return nil, fmt.Errorf("failed to reread t from input: %v", err)
	}

	ans := make([]int64, 0, t)
	for i := int64(0); i < t; i++ {
		var n int
		if _, err := fmt.Fscan(inputReader, &n); err != nil {
			return nil, fmt.Errorf("failed to reread n: %v", err)
		}
		for j := 0; j < n; j++ {
			var w, h int
			fmt.Fscan(inputReader, &w, &h)
		}

		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return nil, fmt.Errorf("expected answer for test %d: %v", i+1, err)
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid integer %q", i+1, token)
		}
		ans = append(ans, val)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected: %s", extra)
	}
	return ans, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
