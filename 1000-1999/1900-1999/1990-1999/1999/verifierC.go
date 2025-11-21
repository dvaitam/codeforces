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

type interval struct {
	l int64
	r int64
}

type testCase struct {
	n    int
	s    int64
	m    int64
	segs []interval
}

const maxTotalN = 200000

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

	expected, err := runAndParse(refBin, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	got, err := runAndParse(candidate, input, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %s got %s\nn=%d s=%d m=%d\nintervals=%v\n", i+1, expected[i], got[i], tests[i].n, tests[i].s, tests[i].m, tests[i].segs)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_1999C.bin"
	cmd := exec.Command("go", "build", "-o", refName, "1999C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runAndParse(target, input string, count int) ([]string, error) {
	out, err := runProgram(target, input)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(out)
	if len(fields) != count {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", count, len(fields), out)
	}
	for i := range fields {
		fields[i] = strings.ToUpper(fields[i])
	}
	return fields, nil
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
	total := totalIntervals(tests)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for total < maxTotalN {
		n := rng.Intn(min(2000, maxTotalN-total)) + 1
		s := randRange(rng, 1, 1_000_000_000)
		minM := int64(2*n + 5)
		if minM < s+1 {
			minM = s + 1
		}
		m := randRange(rng, minM, 1_000_000_000)
		segs := generateIntervals(rng, n, m)
		tests = append(tests, testCase{
			n:    n,
			s:    s,
			m:    m,
			segs: segs,
		})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, s: 3, m: 10, segs: []interval{{l: 3, r: 5}}},
		{n: 2, s: 2, m: 10, segs: []interval{{l: 2, r: 4}, {l: 7, r: 9}}},
	}
}

func generateIntervals(rng *rand.Rand, n int, m int64) []interval {
	segs := make([]interval, n)
	curr := int64(0)
	for i := 0; i < n; i++ {
		remaining := int64(n - i - 1)
		var maxStart int64
		if remaining > 0 {
			maxStart = m - (2*remaining + 2)
		} else {
			maxStart = m - 1
		}
		if maxStart < curr {
			maxStart = curr
		}
		start := curr
		if maxStart > curr {
			start += randRange(rng, 0, maxStart-curr)
		}
		var maxLen int64
		if remaining > 0 {
			maxLen = m - start - (2*remaining + 1)
		} else {
			maxLen = m - start
		}
		if maxLen < 1 {
			maxLen = 1
		}
		length := randRange(rng, 1, maxLen)
		end := start + length
		segs[i] = interval{l: start, r: end}
		curr = end + 1
	}
	return segs
}

func serializeInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.s, tc.m))
		for _, seg := range tc.segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", seg.l, seg.r))
		}
	}
	return sb.String()
}

func totalIntervals(tests []testCase) int {
	total := 0
	for _, tc := range tests {
		total += tc.n
	}
	return total
}

func randRange(rng *rand.Rand, lo, hi int64) int64 {
	if lo == hi {
		return lo
	}
	return lo + int64(rng.Int63n(hi-lo+1))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
