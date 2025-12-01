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
	refSourceF   = "2000-2999/2100-2199/2150-2159/2152/2152F.go"
	randomTrials = 80
	maxNTotal    = 250000
	maxQTotal    = 250000
)

type testCase struct {
	n       int
	z       int
	x       []int
	queries [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceF)
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
	tmp, err := os.CreateTemp("", "2152F-ref-*")
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
		{n: 1, z: 1, x: []int{1}, queries: [][2]int{{1, 1}}},
		{n: 2, z: 1, x: []int{1, 2}, queries: [][2]int{{1, 2}}},
		{n: 6, z: 10, x: []int{1, 5, 7, 8, 11, 12}, queries: [][2]int{{1, 6}, {2, 5}}},
		{n: 6, z: 1, x: []int{1, 1, 1, 3, 3, 3}, queries: [][2]int{{1, 6}, {1, 3}, {4, 6}}},
		{n: 10, z: 5, x: []int{1, 4, 5, 6, 8, 12, 15, 18, 19, 21}, queries: [][2]int{{1, 10}, {3, 7}, {6, 10}}},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5000) + 1
	z := rng.Intn(1_000_000_000) + 1
	x := make([]int, n)
	x[0] = rng.Intn(10) + 1
	for i := 1; i < n; i++ {
		x[i] = x[i-1] + rng.Intn(5) + 1
	}
	q := rng.Intn(5000) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	return testCase{n: n, z: z, x: x, queries: queries}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	totalN := 0
	totalQ := 0
	valid := make([]testCase, 0, len(cases))
	for _, tc := range cases {
		if totalN+tc.n > maxNTotal || totalQ+len(tc.queries) > maxQTotal {
			break
		}
		totalN += tc.n
		totalQ += len(tc.queries)
		valid = append(valid, tc)
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(valid)))
	for _, tc := range valid {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.z))
		for i, val := range tc.x {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
		for _, qr := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qr[0], qr[1]))
		}
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
