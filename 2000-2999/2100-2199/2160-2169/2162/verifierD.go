package main

import (
	"bufio"
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
	refSourceD   = "2000-2999/2100-2199/2160-2169/2162/2162D.go"
	randomTrials = 80
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference(refSourceD)
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
		input := buildInput([]testCase{tc})
		expect, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on case %d: %v", idx+1, err)
		}
		got, err := runCandidate(candidate, input)
		if err != nil {
			fail("candidate failed on case %d: %v", idx+1, err)
		}
		if normalize(expect) != normalize(got) {
			fail("case %d mismatch\nexpected:\n%s\ngot:\n%s", idx+1, expect, got)
		}
	}
	fmt.Printf("All %d cases passed\n", len(tests))
}

type testCase struct {
	n int
	p []int
	l int
	r int
}

func buildReference(src string) (string, error) {
	tmp, err := os.CreateTemp("", "2162D-ref-*")
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
		{n: 3, p: []int{3, 1, 2}, l: 2, r: 2},
		{n: 4, p: []int{2, 1, 3, 4}, l: 2, r: 4},
		{n: 5, p: []int{1, 2, 3, 4, 5}, l: 1, r: 5},
		{n: 6, p: []int{6, 5, 4, 3, 2, 1}, l: 3, r: 3},
	}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})
	l := rng.Intn(n) + 1
	r := rng.Intn(n-l+1) + l
	return testCase{n: n, p: p, l: l, r: r}
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, val := range tc.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.l, tc.r))
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
	lines := strings.Fields(out)
	return strings.Join(lines, "\n")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
