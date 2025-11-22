package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/900-999/930-939/936/936E.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierE.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	answerCount, err := countType2Queries(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference solution: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference solution failed: %v", err)
	}
	expected, err := parseAnswers(refOut, answerCount)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseAnswers(candOut, answerCount)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i, exp := range expected {
		if got[i] != exp {
			fail("wrong answer on query %d: expected %d, got %d", i+1, exp, got[i])
		}
	}

	fmt.Println("OK")
}

func countType2Queries(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	for i := 0; i < n; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return 0, err
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return 0, err
	}
	cnt := 0
	for i := 0; i < q; i++ {
		var t, x, y int
		if _, err := fmt.Fscan(reader, &t, &x, &y); err != nil {
			return 0, err
		}
		if t == 2 {
			cnt++
		} else if t != 1 {
			return 0, fmt.Errorf("unexpected query type %d", t)
		}
	}
	return cnt, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "936E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseAnswers(out string, expectedCount int) ([]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int, 0, expectedCount)
	for len(ans) < expectedCount {
		var v int
		if _, err := fmt.Fscan(reader, &v); err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d answers, got %d", expectedCount, len(ans))
			}
			return nil, err
		}
		ans = append(ans, v)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return ans, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
