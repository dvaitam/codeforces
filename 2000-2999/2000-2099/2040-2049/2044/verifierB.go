package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\n%s", err, expected)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, actual)
		os.Exit(1)
	}

	exp := strings.TrimSpace(expected)
	act := strings.TrimSpace(actual)
	if exp != act {
		fmt.Fprintf(os.Stderr, "output mismatch.\nExpected:\n%s\nGot:\n%s\n", exp, act)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2044B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2044B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(cases []string) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, s := range cases {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []string {
	return []string{
		"p",
		"q",
		"w",
		"pqwpqw",
		"pppppppppp",
		"qqqqqqqqqq",
		"wpqpwqpwq",
	}
}

func randomTests() []string {
	const letters = "pqw"
	cases := make([]string, 0, 200)
	seed := int64(1)
	for len(cases) < cap(cases) {
		seed = seed * 48271 % 2147483647
		n := int(seed%100) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			seed = seed * 48271 % 2147483647
			sb.WriteByte(letters[seed%3])
		}
		cases = append(cases, sb.String())
	}
	return cases
}
