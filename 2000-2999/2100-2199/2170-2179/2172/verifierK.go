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

type testCase struct {
	n, m, q int
	grid    []string
	targets []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2172K-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2172K.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(target string) *exec.Cmd {
	switch filepath.Ext(target) {
	case ".go":
		return exec.Command("go", "run", target)
	case ".py":
		return exec.Command("python3", target)
	default:
		return exec.Command(target)
	}
}

func runProgram(target, input string) (string, error) {
	cmd := commandFor(target)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.q))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	for _, v := range tc.targets {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return sb.String()
}

func manualTests() []string {
	return []string{
		`1 1 1
5
5
`,
		`2 2 2
11
11
1
2
`,
		`3 3 3
1+2
3*4
567
3
7
123
`,
	}
}

func randomGrid(rng *rand.Rand, n, m int) []string {
	chars := []byte("123456789+*")
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = chars[rng.Intn(len(chars))]
		}
		grid[i] = string(row)
	}
	return grid
}

func randomTargets(rng *rand.Rand, q int) []int {
	targets := make([]int, q)
	set := make(map[int]struct{}, q)
	for i := 0; i < q; i++ {
		for {
			val := rng.Intn(1_000_000) + 1
			if _, exists := set[val]; !exists {
				set[val] = struct{}{}
				targets[i] = val
				break
			}
		}
	}
	return targets
}

func randomTest(rng *rand.Rand, n, m, q int) string {
	tc := testCase{
		n:       n,
		m:       m,
		q:       q,
		grid:    randomGrid(rng, n, m),
		targets: randomTargets(rng, q),
	}
	return buildInput(tc)
}

func buildTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for i := 0; i < 40; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(5) + 1
		for n*m > 20 {
			n = rng.Intn(5) + 1
			m = rng.Intn(5) + 1
		}
		q := rng.Intn(10) + 1
		tests = append(tests, randomTest(rng, n, m, q))
	}
	for i := 0; i < 30; i++ {
		n := rng.Intn(15) + 5
		m := rng.Intn(15) + 5
		for n*m > 200 {
			n = rng.Intn(15) + 5
			m = rng.Intn(15) + 5
		}
		q := rng.Intn(100) + 1
		tests = append(tests, randomTest(rng, n, m, q))
	}
	tests = append(tests, randomTest(rng, 100, 300, 1000))
	tests = append(tests, randomTest(rng, 300, 100, 1000))
	tests = append(tests, randomTest(rng, 1000, 30, 500))
	tests = append(tests, randomTest(rng, 30, 1000, 500))
	tests = append(tests, randomTest(rng, 173, 173, 1000))
	tests = append(tests, randomTest(rng, 200, 150, 1000))
	return tests
}

func compareOutputs(expected, got string) error {
	expLines := strings.Split(expected, "\n")
	gotLines := strings.Split(got, "\n")
	var expVals, gotVals []string
	for _, line := range expLines {
		line = strings.TrimSpace(line)
		if line != "" {
			expVals = append(expVals, line)
		}
	}
	for _, line := range gotLines {
		line = strings.TrimSpace(line)
		if line != "" {
			gotVals = append(gotVals, line)
		}
	}
	if len(expVals) != len(gotVals) {
		return fmt.Errorf("expected %d lines, got %d", len(expVals), len(gotVals))
	}
	for i := range expVals {
		if expVals[i] != gotVals[i] {
			return fmt.Errorf("line %d mismatch: expected %s, got %s", i+1, expVals[i], gotVals[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	for idx, input := range tests {
		expect, err := runProgram(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(expect, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, err, input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
