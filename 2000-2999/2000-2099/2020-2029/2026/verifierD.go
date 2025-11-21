package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	a []int64
	q []query
}

type query struct {
	l int64
	r int64
}

const maxTotalN = 300000
const maxQueries = 300000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := serializeInput(tests)
	totalQ := totalQueries(tests)

	expected, err := runAndParse(refBin, input, totalQ)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}

	got, err := runAndParse(candidate, input, totalQ)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	idx := 0
	for tIdx, tc := range tests {
		for qIdx := 0; qIdx < len(tc.q); qIdx++ {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "Mismatch in test %d query %d: expected %d got %d\n", tIdx+1, qIdx+1, expected[idx], got[idx])
				fmt.Fprintf(os.Stderr, "array: %v\nquery: %d %d\n", tc.a, tc.q[qIdx].l, tc.q[qIdx].r)
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2026D.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2026D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, expectedCount int) ([]int64, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d outputs, got %d (output: %q)", expectedCount, len(fields), out)
	}
	res := make([]int64, expectedCount)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	totalN := totalLength(tests)
	totalQ := totalQueries(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for totalN < maxTotalN || totalQ < maxQueries {
		if totalN >= maxTotalN && totalQ >= maxQueries {
			break
		}
		n := rng.Intn(2000) + 1
		if totalN+n > maxTotalN {
			n = maxTotalN - totalN
		}
		if n <= 0 {
			break
		}
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = int64(rng.Intn(21) - 10)
		}
		maxLen := int64(n) * int64(n+1) / 2
		qCount := rng.Intn(2000) + 1
		if totalQ+qCount > maxQueries {
			qCount = maxQueries - totalQ
		}
		if qCount <= 0 {
			qCount = 1
		}
		qs := make([]query, qCount)
		for i := 0; i < qCount; i++ {
			l := randRange(rng, 1, maxLen)
			r := randRange(rng, l, maxLen)
			qs[i] = query{l: l, r: r}
		}
		tests = append(tests, testCase{a: arr, q: qs})
		totalN += n
		totalQ += qCount
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{a: []int64{1}, q: []query{{l: 1, r: 1}}},
		{a: []int64{1, 2, 5, 10}, q: []query{{l: 1, r: 10}, {l: 3, r: 7}}},
		{a: []int64{-10, -10, -10}, q: []query{{l: 1, r: 6}}},
		{a: []int64{10, -10, 10, -10, 10}, q: []query{{l: 3, r: 12}, {l: 5, r: 15}}},
	}
}

func serializeInput(tests []testCase) string {
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
		sb.WriteString(strconv.Itoa(len(tc.q)))
		sb.WriteByte('\n')
		for _, qu := range tc.q {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu.l, qu.r))
		}
	}
	return sb.String()
}

func totalLength(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.a)
	}
	return total
}

func totalQueries(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += len(tc.q)
	}
	return total
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + int64(rng.Int63n(hi-lo+1))
}
