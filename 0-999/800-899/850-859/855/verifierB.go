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

const (
	refSourceB   = "0-999/800-899/850-859/855/855B.go"
	randomTrials = 120
	maxAbs       = 1000000000
)

type testCase struct {
	n int
	p int64
	q int64
	r int64
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceB)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, randomCase(rng))
	}

	for idx, input := range tests {
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
		}
		if normalize(expect) != normalize(got) {
			fail("case %d mismatch\ninput:\n%s\nexpected: %s\ngot: %s", idx+1, input, expect, got)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "855B-ref-*")
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
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicCases() []string {
	var tests []string
	tests = append(tests, buildInput(testCase{
		n: 1,
		p: 1,
		q: 1,
		r: 1,
		a: []int64{5},
	}))
	tests = append(tests, buildInput(testCase{
		n: 5,
		p: 1,
		q: 2,
		r: 3,
		a: []int64{-5, 0, 5, -2, 5},
	}))
	tests = append(tests, buildInput(testCase{
		n: 4,
		p: -2,
		q: -3,
		r: -4,
		a: []int64{3, -1, 2, -2},
	}))
	tests = append(tests, buildInput(testCase{
		n: 6,
		p: 3,
		q: -4,
		r: 2,
		a: []int64{-maxAbs, maxAbs, -maxAbs + 3, 7, 0, 11},
	}))
	large := make([]int64, 100000)
	for i := range large {
		large[i] = int64(i%2*2-1) * int64(1e9-int64(i))
	}
	tests = append(tests, buildInput(testCase{
		n: len(large),
		p: 987654321,
		q: -123456789,
		r: 555555555,
		a: large,
	}))
	return tests
}

func randomCase(rng *rand.Rand) string {
	tc := testCase{}
	tc.n = randomN(rng)
	tc.p = randomCoeff(rng)
	tc.q = randomCoeff(rng)
	tc.r = randomCoeff(rng)
	tc.a = make([]int64, tc.n)
	for i := 0; i < tc.n; i++ {
		tc.a[i] = randomCoeff(rng)
	}
	return buildInput(tc)
}

func randomN(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		return 1
	case 1:
		return rng.Intn(5) + 2
	case 2:
		return rng.Intn(50) + 10
	case 3:
		return rng.Intn(1000) + 100
	default:
		return rng.Intn(90001) + 1000
	}
}

func randomCoeff(rng *rand.Rand) int64 {
	return int64(rng.Int63n(2*maxAbs+1)) - maxAbs
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", tc.n, tc.p, tc.q, tc.r)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
