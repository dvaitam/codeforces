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

const refSource = "2000-2999/2000-2099/2020-2029/2027/2027E1.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierE1.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	testCount, err := countTests(inputData)
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
		fail("reference solution failed: %v", err)
	}
	expected, err := parseWinners(refOut, testCount)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseWinners(candOut, testCount)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i := 0; i < testCount; i++ {
		if got[i] != expected[i] {
			fail("mismatch at test %d: expected %s, got %s", i+1, expected[i], got[i])
		}
	}

	fmt.Println("OK")
}

func countTests(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2027E1-ref-*")
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

func parseWinners(out string, expectedCount int) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]string, 0, expectedCount)
	for len(ans) < expectedCount {
		token, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("expected %d answers, got %d", expectedCount, len(ans))
			}
			return nil, err
		}
		low := strings.ToLower(token)
		if low != "alice" && low != "bob" {
			return nil, fmt.Errorf("invalid answer token %q", token)
		}
		ans = append(ans, low)
	}
	if extra, err := readToken(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return ans, nil
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
