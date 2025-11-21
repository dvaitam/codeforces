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

const refSource = "2000-2999/2000-2099/2060-2069/2061/2061C.go"

type testCase struct {
	a []int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2061C-ref-*")
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
		{a: []int{0}},
		{a: []int{0, 0}},
		{a: []int{1, 1}},
		{a: []int{0, 1, 2}},
		{a: []int{0, 0, 0, 0, 0}},
		{a: []int{0, 0, 0, 0, 0, 0}},
		{a: []int{0, 1, 2, 3, 4}},
		{a: []int{0, 1, 1, 1, 5}},
		{a: []int{0, 1, 2, 3, 4, 5, 6}},
		{a: []int{5, 1, 5, 2, 5}},
		{a: []int{10, 4, 2, 3, 1, 1}},
		{a: []int{2, 3, 1, 1}},
	}
}

func randomTest(rng *rand.Rand, maxN int, allowLarge bool) testCase {
	n := rng.Intn(maxN-1) + 2 // at least 2
	a := make([]int, n)
	limit := n
	if allowLarge && rng.Intn(2) == 0 {
		limit = n + 50
	}
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(limit + 1)
	}
	return testCase{a: a}
}

func heavyTest(n int) testCase {
	a := make([]int, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			a[i] = i
		} else {
			a[i] = n - i
		}
	}
	return testCase{a: a}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.a)))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		res[i] = sc.Text()
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
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
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, 50, true))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, 300, false))
	}
	tests = append(tests, heavyTest(200000))
	tests = append(tests, randomTest(rng, 200000, true))

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\nn=%d a=%v\n",
				i+1, want[i], got[i], len(tests[i].a), tests[i].a)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
