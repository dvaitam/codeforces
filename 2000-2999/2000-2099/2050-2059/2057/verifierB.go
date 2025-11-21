package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	k int
	a []int
}

func runBinary(bin string, input string) (string, error) {
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
		// Samples from the statement
		{n: 1, k: 0, a: []int{48843}},
		{n: 3, k: 1, a: []int{2, 3, 2}},
		{n: 5, k: 3, a: []int{1, 2, 3, 4, 5}},
		{n: 7, k: 0, a: []int{4, 7, 1, 3, 2, 4, 1}},
		{n: 11, k: 4, a: []int{3, 2, 1, 4, 4, 3, 4, 2, 1, 3, 3}},
		{n: 5, k: 5, a: []int{1, 2, 3, 4, 5}},
		// Edge scenarios
		{n: 4, k: 2, a: []int{1, 1, 2, 2}},
		{n: 3, k: 10, a: []int{1, 2, 3}},
		{n: 6, k: 0, a: []int{5, 5, 5, 5, 5, 5}},
		{n: 6, k: 3, a: []int{1, 2, 2, 3, 3, 4}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(30) + 1
	k := rng.Intn(n + 1)
	var limit int
	if rng.Intn(2) == 0 {
		limit = 5
	} else {
		limit = n*2 + 3
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(limit) + 1
		if rng.Intn(20) == 0 {
			// Sprinkle in large values to ensure wide value coverage
			a[i] = rng.Intn(1_000_000_000) + 1
		}
	}
	return testCase{n: n, k: k, a: a}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
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

func parseOutput(out string, t int) ([]int, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]int, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		val, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid integer on test %d: %v", i+1, err)
		}
		res[i] = val
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func minOperations(tc testCase) int {
	freq := make(map[int]int)
	for _, v := range tc.a {
		freq[v]++
	}
	counts := make([]int, 0, len(freq))
	for _, c := range freq {
		counts = append(counts, c)
	}
	sort.Ints(counts)
	distinct := len(counts)
	k := tc.k
	for _, c := range counts {
		if k >= c {
			k -= c
			distinct--
		} else {
			break
		}
	}
	if distinct == 0 {
		return 1
	}
	return distinct
}

func expectedOutputs(tests []testCase) []int {
	ans := make([]int, len(tests))
	for i, tc := range tests {
		ans[i] = minOperations(tc)
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	input := buildInput(tests)
	out, err := runBinary(candidate, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	got, err := parseOutput(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse output: %v\n", err)
		os.Exit(1)
	}
	want := expectedOutputs(tests)
	for i := range tests {
		if got[i] != want[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\nn=%d k=%d a=%v\n",
				i+1, want[i], got[i], tests[i].n, tests[i].k, tests[i].a)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
