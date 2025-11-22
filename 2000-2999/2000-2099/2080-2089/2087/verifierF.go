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
)

const (
	refSource = "2087F.go"
	refBinary = "ref_2087F.bin"
)

type testCase struct {
	n int
	a []int
	b []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	ref, err := buildGoBinary(refSource, refBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	candBin, candCleanup, err := prepareCandidate(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer candCleanup()

	rng := rand.New(rand.NewSource(2087))
	tests := generateTests(rng)

	for i, tc := range tests {
		input := formatInput(tc)

		refOutStr, err := runProgram(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		refAns, err := parseSingleInt(refOutStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, refOutStr)
			os.Exit(1)
		}

		candOutStr, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		candAns, err := parseSingleInt(candOutStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, candOutStr)
			os.Exit(1)
		}

		if refAns != candAns {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d, got %d\ninput:\n%s", i+1, refAns, candAns, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		const candBinary = "candidate_2087F.bin"
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

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSingleInt(output string) (int, error) {
	fields := strings.Fields(output)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer %q: %v", fields[0], err)
	}
	return val, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests(rng *rand.Rand) []testCase {
	tests := deterministicTests()
	for len(tests) < 200 {
		tests = append(tests, randomCase(rng))
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, a: []int{1}, b: []int{1}},
		{n: 1, a: []int{5}, b: []int{5}},
		{n: 2, a: []int{2, 2}, b: []int{2, 2}},
		{n: 3, a: []int{2, 1, 4}, b: []int{1, 4, 2}},
		{n: 4, a: []int{2, 1, 5, 3}, b: []int{5, 1, 3, 2}},                         // sample 1
		{n: 4, a: []int{3, 1, 5, 3}, b: []int{5, 1, 3, 2}},                         // sample 2
		{n: 3, a: []int{1, 1, 1}, b: []int{1, 2, 3}},                               // sample 3
		{n: 8, a: []int{6, 4, 1, 2, 3, 5, 1, 6}, b: []int{6, 3, 2, 1, 2, 3, 4, 3}}, // sample 4
		{n: 5, a: []int{10, 10, 10, 10, 10}, b: []int{10, 10, 10, 10, 10}},
		{n: 6, a: []int{1, 3, 3, 5, 5, 7}, b: []int{4, 1, 6, 1, 6, 1}},
		{n: 6, a: []int{3, 3, 3, 3, 3, 3}, b: []int{1, 2, 3, 4, 5, 6}},
		{n: 6, a: []int{6, 5, 4, 3, 2, 1}, b: []int{1, 2, 3, 4, 5, 6}},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(500) + 1
	a := make([]int, n)
	b := make([]int, n)

	for i := 0; i < n; i++ {
		switch rng.Intn(6) {
		case 0:
			a[i] = 1
			b[i] = rng.Intn(500) + 1
		case 1:
			a[i] = rng.Intn(500) + 1
			b[i] = 1
		case 2:
			val := rng.Intn(4) + 2
			a[i], b[i] = val, val
		case 3:
			high := rng.Intn(300) + 50
			if high > 500 {
				high = 500
			}
			a[i] = high
			b[i] = high
		default:
			a[i] = rng.Intn(500) + 1
			b[i] = rng.Intn(500) + 1
		}
	}
	return testCase{n: n, a: a, b: b}
}
