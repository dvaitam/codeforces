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
	refSource  = "./1266E.go"
	refBinary  = "ref1266E.bin"
	totalTests = 80
)

type milestoneUpdate struct {
	s int
	t int
	u int
}

type testCase struct {
	n       int
	goals   []int64
	updates []milestoneUpdate
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut, len(tc.updates))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, len(tc.updates))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}

		mismatch := false
		for i := range refVals {
			if refVals[i] != candVals[i] {
				mismatch = true
				break
			}
		}
		if mismatch {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sreference:\n%s\ncandidate:\n%s\n", idx+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref-1266E-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1266E.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, g := range tc.goals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(g, 10))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.updates)))
	for _, u := range tc.updates {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u.s, u.t, u.u))
	}
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{
			n:     1,
			goals: []int64{5},
			updates: []milestoneUpdate{
				{s: 1, t: 1, u: 0},
				{s: 1, t: 2, u: 1},
			},
		},
		{
			n:     2,
			goals: []int64{3, 4},
			updates: []milestoneUpdate{
				{s: 1, t: 1, u: 2},
				{s: 2, t: 2, u: 1},
				{s: 1, t: 1, u: 0},
			},
		},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		tests = append(tests, randomTest(rng))
	}
	tests = append(tests,
		heavyTest(200000, 100000, rand.New(rand.NewSource(1))),
		heavyTest(200000, 100000, rand.New(rand.NewSource(2))),
	)
	return tests
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	goals := make([]int64, n)
	for i := range goals {
		goals[i] = int64(rng.Intn(20) + 1)
	}
	q := rng.Intn(10) + 1
	updates := make([]milestoneUpdate, q)
	for i := 0; i < q; i++ {
		s := rng.Intn(n) + 1
		t := rng.Intn(int(goals[s-1]-1)) + 1
		var u int
		if rng.Intn(3) == 0 {
			u = 0
		} else {
			u = rng.Intn(n) + 1
		}
		updates[i] = milestoneUpdate{s: s, t: t, u: u}
	}
	return testCase{n: n, goals: goals, updates: updates}
}

func heavyTest(n, q int, rng *rand.Rand) testCase {
	goals := make([]int64, n)
	for i := range goals {
		goals[i] = int64(rng.Intn(1_000_000_000) + 1)
	}
	updates := make([]milestoneUpdate, q)
	for i := 0; i < q; i++ {
		s := rng.Intn(n) + 1
		t := rng.Intn(int(goals[s-1]-1)) + 1
		var u int
		if rng.Intn(5) == 0 {
			u = 0
		} else {
			u = rng.Intn(n) + 1
		}
		updates[i] = milestoneUpdate{s: s, t: t, u: u}
	}
	return testCase{n: n, goals: goals, updates: updates}
}
