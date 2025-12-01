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

const refSource = "./2124I.go"

type testCase struct {
	n int
	x []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2124I-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, x: []int{1, 2}},
		{n: 2, x: []int{1, 1}},
		{n: 3, x: []int{1, 2, 3}},
		{n: 3, x: []int{1, 1, 1}},
	}
}

func randomTest(rng *rand.Rand, n int) testCase {
	x := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = rng.Intn(i+1) + 1
	}
	return testCase{n: n, x: x}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.x {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

type parsedCase struct {
	yes  bool
	perm []int
}

func parseOutput(out string, tests []testCase) ([]parsedCase, error) {
	fields := strings.Fields(out)
	pos := 0
	res := make([]parsedCase, len(tests))
	for i, tc := range tests {
		if pos >= len(fields) {
			return nil, fmt.Errorf("test %d: missing verdict", i+1)
		}
		verdict := strings.ToUpper(fields[pos])
		pos++
		if verdict == "YES" {
			if pos+tc.n > len(fields) {
				return nil, fmt.Errorf("test %d: expected %d permutation values, got %d", i+1, tc.n, len(fields)-pos)
			}
			perm := make([]int, tc.n)
			seen := make([]bool, tc.n+1)
			for j := 0; j < tc.n; j++ {
				v, err := strconv.Atoi(fields[pos+j])
				if err != nil {
					return nil, fmt.Errorf("test %d: invalid integer %q", i+1, fields[pos+j])
				}
				if v < 1 || v > tc.n {
					return nil, fmt.Errorf("test %d: permutation value out of range %d", i+1, v)
				}
				if seen[v] {
					return nil, fmt.Errorf("test %d: duplicate permutation value %d", i+1, v)
				}
				seen[v] = true
				perm[j] = v
			}
			pos += tc.n
			res[i] = parsedCase{yes: true, perm: perm}
		} else if verdict == "NO" {
			res[i] = parsedCase{yes: false}
		} else {
			return nil, fmt.Errorf("test %d: invalid verdict %q", i+1, verdict)
		}
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra output detected after %d testcases", len(tests))
	}
	return res, nil
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 && totalN(tests) < 150000 {
		n := rng.Intn(200) + 2
		tests = append(tests, randomTest(rng, n))
	}
	for len(tests) < 60 && totalN(tests) < 190000 {
		n := rng.Intn(2000) + 500
		tests = append(tests, randomTest(rng, n))
	}
	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	if len(want) != len(got) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d got %d\n", len(want), len(got))
		os.Exit(1)
	}

	for i := range want {
		if want[i].yes != got[i].yes {
			fmt.Fprintf(os.Stderr, "test %d: verdict mismatch expected %v got %v\n", i+1, want[i].yes, got[i].yes)
			os.Exit(1)
		}
		if want[i].yes {
			// Only need to ensure permutation validity (handled in parse). No need to match exact permutation.
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
