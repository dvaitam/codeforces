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

const refSource = "2000-2999/2000-2099/2060-2069/2060/2060D.go"

type testCase struct {
	a []int64
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2060D-ref-*")
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

func runProgram(bin string, input string) (string, error) {
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
		{a: []int64{1, 2, 3, 4, 5}},
		{a: []int64{4, 3, 2, 1}},
		{a: []int64{4, 5, 2, 3}},
		{a: []int64{4, 5, 4, 5, 4, 5, 4, 5}},
		{a: []int64{9, 9, 8, 2, 4, 4, 3, 5, 3}},
		{a: []int64{10, 10}},
		{a: []int64{1, 1000000000}},
		{a: []int64{1000000000, 1}},
	}
}

func randomTest(rng *rand.Rand, maxN int, maxVal int64) testCase {
	n := rng.Intn(maxN-1) + 2 // at least 2
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(maxVal) + 1
	}
	return testCase{a: a}
}

func heavyTest(n int) testCase {
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			a[i] = int64(n - i)
		} else {
			a[i] = int64(i + 1)
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
			sb.WriteString(strconv.FormatInt(v, 10))
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
		token := strings.ToUpper(strings.TrimSpace(sc.Text()))
		if token != "YES" && token != "NO" {
			return nil, fmt.Errorf("invalid verdict on test %d: %q", i+1, sc.Text())
		}
		res[i] = token
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
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
		tests = append(tests, randomTest(rng, 50, 1_000_000_000))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, 300, 10_000))
	}
	tests = append(tests, heavyTest(200000))
	tests = append(tests, randomTest(rng, 200000, 1_000_000_000))

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
