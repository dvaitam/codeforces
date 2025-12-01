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
	refSourceB   = "./2150B.go"
	randomTrials = 200
	maxNTotal    = 200000
	mod          = 998244353
)

type testCase struct {
	n int
	a []int
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

	input := buildInput(tests)

	expect, err := runProgram(refBin, input)
	if err != nil {
		fail("reference failed: %v", err)
	}
	got, err := runCandidate(candidate, input)
	if err != nil {
		fail("candidate failed: %v", err)
	}
	if normalize(expect) != normalize(got) {
		fail("output mismatch\nexpected:\n%s\ngot:\n%s", expect, got)
	}
	fmt.Printf("All %d cases passed\n", len(tests))
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "2150B-ref-*")
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
		{n: 2, a: []int{2, 0}},
		{n: 2, a: []int{1, 1}},
		{n: 4, a: []int{2, 2, 1, 0}},
		{n: 5, a: []int{2, 2, 1, 0, 0}},
		{n: 3, a: []int{1, 0, 0}},
		{n: 6, a: []int{1, 1, 1, 1, 1, 1}},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 2
	a := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(n + 1)
		sum += a[i]
	}
	if rng.Intn(4) == 0 {
		for i := 0; i < n; i++ {
			a[i] = 0
		}
		sum = 0
	}
	return testCase{n: n, a: a}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	totalN := 0
	valid := make([]testCase, 0, len(cases))
	for _, tc := range cases {
		if totalN+tc.n > maxNTotal {
			break
		}
		totalN += tc.n
		valid = append(valid, tc)
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(valid)))
	for _, tc := range valid {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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

func normalize(out string) string {
	return strings.TrimSpace(out)
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
