package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const refSource = "2000-2999/2000-2099/2050-2059/2051/2051C.go"

type testCase struct {
	n int
	m int
	k int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierC.go /path/to/candidate")
	}
	target := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInputStructure(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAns, err := parseOutputs(refOut, tests)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(target), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAns, err := parseOutputs(candOut, tests)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if len(refAns) != len(candAns) {
		fail("test count mismatch: reference %d candidate %d", len(refAns), len(candAns))
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fail("test %d mismatch\nexpected: %s\ngot: %s", i+1, refAns[i], candAns[i])
		}
	}

	fmt.Println("OK")
}

func parseInputStructure(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		tc := testCase{}
		if _, err := fmt.Fscan(reader, &tc.n, &tc.m, &tc.k); err != nil {
			return nil, err
		}
		for j := 0; j < tc.m; j++ {
			var dummy int
			if _, err := fmt.Fscan(reader, &dummy); err != nil {
				return nil, err
			}
		}
		for j := 0; j < tc.k; j++ {
			var dummy int
			if _, err := fmt.Fscan(reader, &dummy); err != nil {
				return nil, err
			}
		}
		tests[i] = tc
	}
	return tests, nil
}

func parseOutputs(out string, tests []testCase) ([]string, error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(out)))
	results := make([]string, len(tests))
	for idx, tc := range tests {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("missing output for test %d: %v", idx+1, err)
		}
		if len(token) != tc.m {
			return nil, fmt.Errorf("test %d: expected string of length %d, got %q", idx+1, tc.m, token)
		}
		for pos, ch := range token {
			if ch != '0' && ch != '1' {
				return nil, fmt.Errorf("test %d: invalid character %q at position %d", idx+1, ch, pos+1)
			}
		}
		results[idx] = token
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("extra token %q after outputs", extra)
		}
		return nil, err
	}
	return results, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2051C-ref-*")
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

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
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

func readToken(r *bufio.Reader) (string, error) {
	var sb bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		if !isSpace(b) {
			sb.WriteByte(b)
			break
		}
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(b) {
			break
		}
		sb.WriteByte(b)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v' || b == '\f'
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
