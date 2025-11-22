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

const refSource = "2000-2999/2000-2099/2040-2049/2042/2042F.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierF.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	answerCount, err := countType3Queries(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	expected, err := parseInts(refOut, answerCount)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseInts(candOut, answerCount)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i := 0; i < answerCount; i++ {
		if got[i] != expected[i] {
			fail("wrong answer on query %d: expected %d, got %d", i+1, expected[i], got[i])
		}
	}

	fmt.Println("OK")
}

func countType3Queries(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	for i := 0; i < n; i++ {
		var tmp int64
		if _, err := fmt.Fscan(reader, &tmp); err != nil {
			return 0, err
		}
	}
	for i := 0; i < n; i++ {
		var tmp int64
		if _, err := fmt.Fscan(reader, &tmp); err != nil {
			return 0, err
		}
	}
	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return 0, err
	}
	cnt := 0
	for i := 0; i < q; i++ {
		var t int
		if _, err := fmt.Fscan(reader, &t); err != nil {
			return 0, err
		}
		switch t {
		case 1, 2:
			var p int
			var x int64
			if _, err := fmt.Fscan(reader, &p, &x); err != nil {
				return 0, err
			}
		case 3:
			var l, r int
			if _, err := fmt.Fscan(reader, &l, &r); err != nil {
				return 0, err
			}
			cnt++
		default:
			return 0, fmt.Errorf("unknown query type %d", t)
		}
	}
	return cnt, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2042F-ref-*")
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

func parseInts(out string, expectedCount int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, 0, expectedCount)
	for len(ans) < expectedCount {
		var v int64
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
