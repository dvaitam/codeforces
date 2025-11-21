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

const refSource = "2000-2999/2000-2099/2030-2039/2030/2030E.go"

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierE.go /path/to/candidate")
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	t, err := parseTestCount(input)
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
	refAns, err := parseAnswers(refOut, t)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAns, err := parseAnswers(candOut, t)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	for i := 0; i < t; i++ {
		if candAns[i] != refAns[i] {
			fail("mismatch on test %d: expected %d got %d", i+1, refAns[i], candAns[i])
		}
	}

	fmt.Println("OK")
}

func parseTestCount(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		token, err := readToken(reader)
		if err != nil {
			return nil, fmt.Errorf("expected %d answers, got %d: %v", t, i, err)
		}
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("answer %d is not an integer: %v", i+1, err)
		}
		ans[i] = val
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("extra token %q after %d answers", extra, t)
		}
		return nil, err
	}
	return ans, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2030E-ref-*")
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
	var sb strings.Builder
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
