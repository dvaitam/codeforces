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
	refSourceA = "./1551A.go"
	randomSets = 120
)

type testCase struct {
	input string
	ns    []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceA)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomCases(rng, randomSets)...)

	for idx, tc := range tests {
		expOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		expPairs, err := parseOutputs(expOut, len(tc.ns))
		if err != nil {
			fail("reference output parse error on case %d: %v\noutput:\n%s", idx+1, err, expOut)
		}
		minDiffs := computeDiffs(tc.ns, expPairs)

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, tc.input)
		}
		gotPairs, err := parseOutputs(gotOut, len(tc.ns))
		if err != nil {
			fail("candidate output parse error on case %d: %v\noutput:\n%s", idx+1, err, gotOut)
		}
		for i, n := range tc.ns {
			c1 := gotPairs[i][0]
			c2 := gotPairs[i][1]
			if c1 < 0 || c2 < 0 {
				fail("case %d test %d: negative coin count (%d,%d)", idx+1, i+1, c1, c2)
			}
			if c1+2*c2 != n {
				fail("case %d test %d: equation violated (%d + 2*%d != %d)", idx+1, i+1, c1, c2, n)
			}
			diff := c1 - c2
			if diff < 0 {
				diff = -diff
			}
			if diff != minDiffs[i] {
				fail("case %d test %d: non-minimal difference (got %d, expected %d)", idx+1, i+1, diff, minDiffs[i])
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1551A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(src))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	return runCommand(cmd, input)
}

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	return runCommand(cmd, input)
}

func runCommand(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicCases() []testCase {
	return []testCase{
		buildCase([]int64{1}),
		buildCase([]int64{2, 3, 4}),
		buildCase([]int64{1000, 30, 999999937}),
		buildCase([]int64{5, 6, 7, 8, 9, 10}),
	}
}

func randomCases(rng *rand.Rand, count int) []testCase {
	cases := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		t := randomT(rng)
		ns := make([]int64, t)
		for j := 0; j < t; j++ {
			ns[j] = randomN(rng)
		}
		cases = append(cases, buildCase(ns))
	}
	return cases
}

func randomT(rng *rand.Rand) int {
	switch rng.Intn(4) {
	case 0:
		return 1
	case 1:
		return rng.Intn(10) + 1
	case 2:
		return rng.Intn(200) + 50
	default:
		return rng.Intn(1000) + 1
	}
}

func randomN(rng *rand.Rand) int64 {
	switch rng.Intn(4) {
	case 0:
		return int64(rng.Intn(10) + 1)
	case 1:
		return int64(rng.Intn(1000) + 1)
	case 2:
		return int64(rng.Intn(1000000) + 1)
	default:
		return int64(rng.Intn(1_000_000_000) + 1)
	}
}

func buildCase(ns []int64) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ns)))
	for _, n := range ns {
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return testCase{
		input: sb.String(),
		ns:    ns,
	}
}

func parseOutputs(out string, t int) ([][2]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*t {
		return nil, fmt.Errorf("expected %d numbers, got %d", 2*t, len(fields))
	}
	res := make([][2]int64, t)
	for i := 0; i < t; i++ {
		v1, err := strconv.ParseInt(fields[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i])
		}
		v2, err := strconv.ParseInt(fields[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i+1])
		}
		res[i] = [2]int64{v1, v2}
	}
	return res, nil
}

func computeDiffs(ns []int64, pairs [][2]int64) []int64 {
	diffs := make([]int64, len(ns))
	for i, n := range ns {
		c1 := pairs[i][0]
		c2 := pairs[i][1]
		if c1+2*c2 != n {
			// reference should not fail; safeguard
			fail("reference produced invalid solution for n=%d: (%d,%d)", n, c1, c2)
		}
		diff := c1 - c2
		if diff < 0 {
			diff = -diff
		}
		diffs[i] = diff
	}
	return diffs
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
