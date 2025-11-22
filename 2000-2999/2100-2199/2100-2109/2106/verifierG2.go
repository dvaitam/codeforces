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

const refSource = "2000-2999/2100-2199/2100-2109/2106/2106G2.go"

type testCase struct {
	n int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierG2.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(inputData)
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
	expected, err := parseOutputs(refOut, tests)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i, exp := range expected {
		act := got[i]
		if len(exp) != len(act) {
			fail("test %d: expected %d values, got %d", i+1, len(exp), len(act))
		}
		for j := range exp {
			if exp[j] != act[j] {
				fail("test %d, position %d: expected %d, got %d", i+1, j+1, exp[j], act[j])
			}
		}
	}

	fmt.Println("OK")
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		var root int
		if _, err := fmt.Fscan(reader, &n, &root); err != nil {
			return nil, err
		}
		tests[i].n = n
		_ = root

		for j := 0; j < n; j++ {
			var val int
			if _, err := fmt.Fscan(reader, &val); err != nil {
				return nil, err
			}
		}
		for j := 0; j < n-1; j++ {
			var u, v int
			if _, err := fmt.Fscan(reader, &u, &v); err != nil {
				return nil, err
			}
			_ = u
			_ = v
		}
	}
	return tests, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2106G2-ref-*")
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

func parseOutputs(out string, tests []testCase) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	res := make([][]int, len(tests))
	for idx, tc := range tests {
		res[idx] = make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &res[idx][i]); err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("test %d: expected %d integers, got %d", idx+1, tc.n, i)
				}
				return nil, err
			}
		}
	}
	if extra, err := readToken(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return res, nil
}

func readToken(r *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch > ' ' {
			b.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			return b.String(), nil
		}
		b.WriteByte(ch)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
