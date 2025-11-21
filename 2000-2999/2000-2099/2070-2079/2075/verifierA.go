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

const refSource2075A = "2000-2999/2000-2099/2070-2079/2075/2075A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
	dir, err := os.MkdirTemp("", "cf-2075A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2075A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource2075A)
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
	inputFields := strings.Fields(input)
	if len(inputFields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(inputFields[0])
	if err != nil || t < 1 || t > 10000 {
		return nil, fmt.Errorf("invalid test count %q", inputFields[0])
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
		if v < 0 {
			return nil, fmt.Errorf("answers must be non-negative, got %d", v)
		}
		res[i] = v
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeManual("sample_like", [][2]int64{{39, 7}, {9, 3}, {6, 3}, {999967802, 35}, {5, 5}, {999999999, 3}, {1000000000, 3}}),
		makeManual("small_odd_even", [][2]int64{{3, 3}, {4, 3}}),
		makeManual("k_equals_n", [][2]int64{{7, 7}, {8, 7}}),
		makeManual("mix_parity", [][2]int64{{1001, 9}, {1000, 9}, {999999937, 999999937}}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func makeManual(name string, pairs [][2]int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs)))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
	}
	return testCase{name: name, input: sb.String()}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(6) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := int64(rng.Intn(1_000_000_000-2) + 3)
		k := randomOddLE(rng, int(n))
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	}
	return testCase{name: fmt.Sprintf("random_%d", idx+1), input: sb.String()}
}

func randomOddLE(rng *rand.Rand, limit int) int {
	k := rng.Intn(limit-2) + 3
	if k%2 == 0 {
		k++
		if k > limit {
			k -= 2
		}
	}
	return k
}
