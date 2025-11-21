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
	refSourceC   = "0-999/800-899/880-889/883/883C.go"
	randomTrials = 150
	limitValue   = 10_000_000
)

type inputCase struct {
	f, T, t0   int64
	a1, t1, p1 int64
	a2, t2, p2 int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceC)
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < randomTrials; i++ {
		tests = append(tests, buildInput(randomCase(rng)))
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
	tmp, err := os.CreateTemp("", "883C-ref-*")
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
	var cases []string
	// Already enough time using regular tariff.
	cases = append(cases, buildInput(inputCase{
		f: 10, T: 200, t0: 20,
		a1: 5, t1: 15, p1: 3,
		a2: 10, t2: 15, p2: 5,
	}))
	// Need only first package (similar to statement).
	cases = append(cases, buildInput(inputCase{
		f: 120, T: 960, t0: 10,
		a1: 26, t1: 8, p1: 8,
		a2: 40, t2: 9, p2: 10,
	}))
	// Need combination of packages.
	cases = append(cases, buildInput(inputCase{
		f: 200, T: 800, t0: 10,
		a1: 30, t1: 6, p1: 7,
		a2: 50, t2: 5, p2: 15,
	}))
	// Impossible case (-1).
	cases = append(cases, buildInput(inputCase{
		f: 1000, T: 1500, t0: 3,
		a1: 100, t1: 3, p1: 10,
		a2: 250, t2: 4, p2: 20,
	}))
	// Large boundary case.
	cases = append(cases, buildInput(inputCase{
		f: limitValue, T: limitValue, t0: 2,
		a1: 3000000, t1: 1, p1: 100,
		a2: 1500000, t2: 1, p2: 60,
	}))
	return cases
}

func randomCase(rng *rand.Rand) inputCase {
	tc := inputCase{}
	tc.f = randRange(rng, 1, limitValue)
	tc.t0 = randRange(rng, 1, limitValue)
	// Some tests have T equal to regular time when possible.
	if tc.f*tc.t0 <= limitValue && rng.Intn(3) == 0 {
		tc.T = tc.f * tc.t0
	} else {
		tc.T = randRange(rng, 1, limitValue)
	}
	tc.a1 = randRange(rng, 1, limitValue)
	tc.t1 = randRange(rng, 1, limitValue)
	tc.p1 = randRange(rng, 1, limitValue)
	tc.a2 = randRange(rng, 1, limitValue)
	tc.t2 = randRange(rng, 1, limitValue)
	tc.p2 = randRange(rng, 1, limitValue)
	// Occasionally ensure a very fast package exists.
	if rng.Intn(4) == 0 {
		if tc.t1 > tc.t0 {
			tc.t1 = randRange(rng, 1, tc.t0)
		}
	}
	if rng.Intn(4) == 0 {
		if tc.t2 > tc.t0 {
			tc.t2 = randRange(rng, 1, tc.t0)
		}
	}
	return tc
}

func buildInput(tc inputCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.f, tc.T, tc.t0)
	fmt.Fprintf(&sb, "%d %d %d\n", tc.a1, tc.t1, tc.p1)
	fmt.Fprintf(&sb, "%d %d %d\n", tc.a2, tc.t2, tc.p2)
	return sb.String()
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	if hi < lo {
		lo, hi = hi, lo
	}
	if lo == hi {
		return lo
	}
	return lo + rng.Int63n(hi-lo+1)
}

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
