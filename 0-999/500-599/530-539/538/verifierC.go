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

const refSource = "0-999/500-599/530-539/538/538C.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if !equalOutputs(refOut, candOut) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier538C-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func equalOutputs(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

func buildTests() []testCase {
	tests := []testCase{
		{name: "sample1", input: "8 2\n2 0\n7 0\n"},
		{name: "sample2", input: "8 3\n2 0\n7 0\n8 3\n"},
		{name: "single_note", input: "5 1\n3 2\n"},
	}
	tests = append(tests, randomTests(30)...)
	return tests
}

func randomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(1000) + 1
		m := rng.Intn(20) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		prevDay := 0
		for j := 0; j < m; j++ {
			day := rng.Intn(n-prevDay) + prevDay + 1
			height := rng.Intn(1000)
			fmt.Fprintf(&sb, "%d %d\n", day, height)
			prevDay = day
			if prevDay >= n {
				break
			}
		}
		tests = append(tests, testCase{name: fmt.Sprintf("random_%d", i+1), input: sb.String()})
	}
	return tests
}
