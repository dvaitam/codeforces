package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceF = "2000-2999/2000-2099/2050-2059/2055/2055F.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	tc, err := parseTestCount(input)
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
	exp, err := parseAnswers(refOut, tc)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	got, err := parseAnswers(userOut, tc)
	if err != nil {
		fail("failed to parse candidate output: %v\n%s", err, userOut)
	}

	for i := 0; i < tc; i++ {
		if exp[i] != got[i] {
			fail("answer mismatch at test %d: expected %s got %s", i+1, exp[i], got[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2055F-ref-*")
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

func parseTestCount(input []byte) (int, error) {
	fields := strings.Fields(string(input))
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty input")
	}
	var t int
	_, err := fmt.Sscan(fields[0], &t)
	if err != nil {
		return 0, fmt.Errorf("invalid t: %v", err)
	}
	return t, nil
}

func parseAnswers(out string, t int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	ans := make([]string, t)
	for i, f := range fields {
		up := strings.ToUpper(f)
		if up != "YES" && up != "NO" {
			return nil, fmt.Errorf("invalid answer %q at position %d", f, i+1)
		}
		ans[i] = up
	}
	return ans, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
