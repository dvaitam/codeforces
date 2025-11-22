package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource = "2095I.go"
	refBinary = "ref_2095I.bin"
)

type testCase struct {
	a int64
	b int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		return
	}
	candidatePath := os.Args[1]

	ref, err := buildGoBinary(refSource, refBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	candBin, cleanup, err := prepareCandidate(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)

	for i, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		refVal, err := parseSingleInt(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candVal, err := parseSingleInt(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, candOut)
			os.Exit(1)
		}

		if refVal != candVal {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\ninput:\n%s", i+1, refVal, candVal, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		const candBinary = "candidate_2095I.bin"
		bin, err := buildGoBinary(path, candBinary)
		if err != nil {
			return "", func() {}, err
		}
		return bin, func() { os.Remove(bin) }, nil
	}
	return path, func() {}, nil
}

func buildGoBinary(source, output string) (string, error) {
	cmd := exec.Command("go", "build", "-o", output, source)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", output), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSingleInt(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer %q: %v", fields[0], err)
	}
	return val, nil
}

func formatInput(tc testCase) string {
	return fmt.Sprintf("%d %d\n", tc.a, tc.b)
}

func generateTests(rng *rand.Rand) []testCase {
	tests := []testCase{
		{a: 1, b: 1},
		{a: 1, b: 1_000_000_000},
		{a: 1_000_000_000, b: 1_000_000_000},
		{a: 42, b: 58},
		{a: 123456789, b: 987654321},
	}

	for len(tests) < 200 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func randomCase(rng *rand.Rand) testCase {
	a := rng.Int63n(1_000_000_000) + 1
	b := rng.Int63n(1_000_000_000) + 1
	return testCase{a: a, b: b}
}
