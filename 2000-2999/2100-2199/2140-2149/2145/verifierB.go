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
	refSourceB   = "2000-2999/2100-2199/2140-2149/2145/2145B.go"
	randomTrials = 200
	maxNTotal    = 200000
)

type testCase struct {
	n int
	k int
	s string
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
	tmp, err := os.CreateTemp("", "2145B-ref-*")
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
		{n: 4, k: 2, s: "01"},
		{n: 3, k: 3, s: "222"},
		{n: 5, k: 5, s: "01221"},
		{n: 1, k: 1, s: "0"},
		{n: 2, k: 2, s: "12"},
		{n: 6, k: 4, s: "1022"},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(200) + 1
	k := rng.Intn(n) + 1
	bytes := make([]byte, k)
	for i := 0; i < k; i++ {
		switch rng.Intn(3) {
		case 0:
			bytes[i] = '0'
		case 1:
			bytes[i] = '1'
		default:
			bytes[i] = '2'
		}
	}
	return testCase{n: n, k: k, s: string(bytes)}
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		sb.WriteString(tc.s)
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
