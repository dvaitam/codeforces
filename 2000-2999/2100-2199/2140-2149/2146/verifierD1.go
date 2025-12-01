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
	refSourceD1  = "2000-2999/2100-2199/2140-2149/2146/2146D1.go"
	randomTrials = 150
	maxR         = 200000
)

type testCase struct {
	l int
	r int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceD1)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomCase(rng))
	}

	for idx, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.l, tc.r)
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		expVal, _, err := parseOutput(expect, tc)
		if err != nil {
			fail("reference output invalid on case %d: %v\noutput:\n%s", idx+1, err, expect)
		}

		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		gotVal, gotPerm, err := parseOutput(got, tc)
		if err != nil {
			fail("candidate output invalid on case %d: %v\noutput:\n%s", idx+1, err, got)
		}

		if gotVal != expVal {
			fail("case %d: mismatched score (expected %d, got %d)", idx+1, expVal, gotVal)
		}

		if !isPermutation(gotPerm, tc.l, tc.r) {
			fail("case %d: candidate permutation is invalid", idx+1)
		}

		score := computeScore(gotPerm, tc.l, tc.r)
		if score != gotVal {
			fail("case %d: reported score %d doesn't match computed score %d", idx+1, gotVal, score)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "2146D1-ref-*")
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

func deterministicCases() []testCase {
	return []testCase{
		{l: 0, r: 0},
		{l: 0, r: 1},
		{l: 0, r: 3},
		{l: 0, r: 7},
		{l: 0, r: 15},
		{l: 0, r: 31},
		{l: 0, r: 42},
		{l: 0, r: 100},
		{l: 0, r: 127},
		{l: 0, r: 255},
	}
}

func randomCase(rng *rand.Rand) testCase {
	r := rng.Intn(maxR)
	if r < 1 {
		r = 1
	}
	if rng.Intn(5) == 0 {
		r = (1 << rng.Intn(18)) - 1
		if r < 0 {
			r = 0
		}
		if r >= maxR {
			r = maxR - 1
		}
	}
	return testCase{l: 0, r: r}
}

func parseOutput(out string, tc testCase) (int64, []int, error) {
	lines := strings.Fields(out)
	if len(lines) < 1+(tc.r-tc.l+1) {
		return 0, nil, fmt.Errorf("insufficient tokens")
	}
	val, err := parseInt64(lines[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid score: %v", err)
	}
	n := tc.r - tc.l + 1
	if len(lines) < 1+n {
		return 0, nil, fmt.Errorf("missing permutation values")
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		x, err := parseInt(lines[1+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid permutation value at position %d: %v", i, err)
		}
		perm[i] = x
	}
	return val, perm, nil
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func isPermutation(perm []int, l, r int) bool {
	n := r - l + 1
	if len(perm) != n {
		return false
	}
	seen := make([]bool, r+1)
	for _, x := range perm {
		if x < l || x > r || seen[x] {
			return false
		}
		seen[x] = true
	}
	return true
}

func computeScore(perm []int, l, r int) int64 {
	n := r - l + 1
	score := int64(0)
	for i := 0; i < n; i++ {
		score += int64(perm[i] | (l + i))
	}
	return score
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func runCandidate(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
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

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
