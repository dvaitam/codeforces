package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2100-2109/2103/2103B.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2103B-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string, expected int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(fields))
	}
	return fields, nil
}

func fixedTests() []string {
	return []string{
		"6\n3\n000\n3\n111\n3\n011\n3\n100\n5\n10101\n19\n1010101010010011011\n",
		"4\n1\n0\n1\n1\n2\n01\n2\n10\n",
	}
}

func randomString(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = '0'
		} else {
			b[i] = '1'
		}
	}
	return string(b)
}

func randomTestInput(rng *rand.Rand) string {
	var cases []string
	remaining := 200000
	for remaining > 0 && len(cases) < 60 {
		maxLen := remaining
		if maxLen > 20000 {
			maxLen = 20000
		}
		n := rng.Intn(maxLen) + 1
		s := randomString(rng, n)
		cases = append(cases, fmt.Sprintf("%d\n%s\n", n, s))
		remaining -= n
		if rng.Intn(4) == 0 {
			break
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(c)
	}
	return sb.String()
}

func structuredCases() string {
	var sb strings.Builder
	sb.WriteString("5\n")
	sb.WriteString("5\n00000\n")       // all zero
	sb.WriteString("5\n11111\n")       // all one
	sb.WriteString("6\n001111\n")      // single transition
	sb.WriteString("6\n111000\n")      // transition other way
	sb.WriteString("10\n0101010101\n") // alternating
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := fixedTests()
	tests = append(tests, structuredCases())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// include a maximum size case with random string
	maxStr := randomString(rng, 200000)
	tests = append(tests, fmt.Sprintf("1\n200000\n%s\n", maxStr))
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTestInput(rng))
	}

	for idx, input := range tests {
		var t int
		if _, err := fmt.Fscan(strings.NewReader(input), &t); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse test count for case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		expOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expLines, err := parseOutput(expOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse failed on case %d: %v\noutput:\n%s", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotLines, err := parseOutput(gotOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse failed on case %d: %v\noutput:\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := 0; i < t; i++ {
			if expLines[i] != gotLines[i] {
				fmt.Fprintf(os.Stderr, "mismatch on case %d test %d: expected %s got %s\ninput:\n%s", idx+1, i+1, expLines[i], gotLines[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Println("All tests passed")
}
