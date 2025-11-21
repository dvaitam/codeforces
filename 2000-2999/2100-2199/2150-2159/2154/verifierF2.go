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

const refSource2154F2 = "2000-2999/2100-2199/2150-2159/2154/2154F2.go"
const mod = 998244353

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVals, err := parseOutput(tc.input, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(tc.input, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d answers got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, len(refVals), len(candVals), tc.input, refOut, candOut)
			os.Exit(1)
		}
		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
					idx+1, tc.name, i+1, refVals[i], candVals[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2154F2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2154F2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2154F2)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(input, output string) ([]int64, error) {
	inFields := strings.Fields(input)
	if len(inFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inFields[0])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", inFields[0])
	}
	outFields := strings.Fields(output)
	if len(outFields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(outFields))
	}
	res := make([]int64, t)
	for i, tok := range outFields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		if v < 0 || v >= mod {
			return nil, fmt.Errorf("answer %d out of range: %d", i+1, v)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		buildSample(),
		makeManual("all_known_sorted", []int{1, 2, 3, 4, 5}),
		makeManual("all_unknown", []int{-1, -1, -1, -1}),
		makeManual("prefix_known", []int{1, 2, -1, -1, -1}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 140; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	tests = append(tests, largeTest())
	return tests
}

func buildSample() testCase {
	const sampleInput = `7
5
-1 -1 -1 -1 -1
4
1 2 3 4
5
-1 -1 -1 2 -1
6
-1 3 2 1 -1 -1
18
11 -1 2 -1 -1 -1 -1 6 -1 -1 14 8 9 15 -1 -1 -1 -1
6
-1 3 -1 4 -1 5
3
-1 2 1
`
	return testCase{name: "statement_sample", input: sampleInput}
}

func makeManual(name string, perm []int) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(perm)))
	writeArray(&sb, perm)
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	totalN := 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(50000) + 2
		totalN += n
		perm := buildRandomPartialPermutation(rng, n)
		sb.WriteString(fmt.Sprintf("%d\n", n))
		writeArray(&sb, perm)
		sb.WriteByte('\n')
	}
	_ = totalN
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}

func buildRandomPartialPermutation(rng *rand.Rand, n int) []int {
	values := rng.Perm(n)[:rng.Intn(n)+1]
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = -1
	}
	for i, v := range values {
		pos := rng.Intn(n)
		for perm[pos] != -1 {
			pos = rng.Intn(n)
		}
		perm[pos] = v + 1
		if i == n-1 {
			break
		}
	}
	return perm
}

func largeTest() testCase {
	n := 200000
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = -1
	}
	perm[0] = 1
	perm[n-1] = n
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	writeArray(&sb, perm)
	sb.WriteByte('\n')
	return testCase{name: "large_partial", input: sb.String()}
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
}
