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
)

const refSource = "2000-2999/2000-2099/2030-2039/2036/2036F.go"

type query struct {
	l, r int64
	i    int
	k    int64
}

type testCase struct {
	name    string
	queries []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(len(tc.queries), refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		expected := simulate(tc)
		if !equalAnswers(refAns, expected) {
			fmt.Fprintf(os.Stderr, "reference mismatch simulation on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(len(tc.queries), candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(candAns, expected) {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, tc.name, input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2036F-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2036F.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(expected int, output string) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		ans[i] = val
	}
	return ans, nil
}

func equalAnswers(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", q.l, q.r, q.i, q.k))
	}
	return sb.String()
}

func simulate(tc testCase) []int64 {
	res := make([]int64, len(tc.queries))
	for idx, q := range tc.queries {
		res[idx] = answerQuery(q.l, q.r, q.i, q.k)
	}
	return res
}

func answerQuery(l, r int64, i int, k int64) int64 {
	total := xorRange(l, r)
	if i == 0 {
		// modulo 1, all numbers discarded if k == 0
		return 0
	}
	m := int64(1) << uint(i)
	k %= m
	delta := (k - l) % m
	if delta < 0 {
		delta += m
	}
	first := l + delta
	if first > r {
		return total
	}
	cnt := (r-first)/m + 1
	start := (first - k) / m
	xorIdx := xorRange(start, start+cnt-1)
	var subset int64
	if cnt&1 == 1 {
		subset = k
	}
	subset ^= xorIdx << uint(i)
	return total ^ subset
}

func xorPrefix(n int64) int64 {
	switch n & 3 {
	case 0:
		return n
	case 1:
		return 1
	case 2:
		return n + 1
	default:
		return 0
	}
}

func xorRange(a, b int64) int64 {
	if a > b {
		return 0
	}
	return xorPrefix(b) ^ xorPrefix(a-1)
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample",
			queries: []query{
				{1, 3, 1, 0},
				{2, 28, 3, 7},
				{15, 43, 1, 0},
				{57, 2007, 1, 0},
				{1010, 1993, 2, 2},
				{1, 1000000000, 30, 1543},
			},
		},
		{
			name:    "edge_cases",
			queries: []query{{1, 1, 0, 0}, {5, 5, 1, 0}, {10, 10, 30, 0}},
		},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	qcnt := rng.Intn(50) + 1
	queries := make([]query, qcnt)
	for i := 0; i < qcnt; i++ {
		l := rng.Int63n(1_000_000) + 1
		r := l + rng.Int63n(1_000_000)
		iVal := rng.Intn(31)
		var k int64
		if iVal > 0 {
			k = rng.Int63n(1 << uint(iVal))
		}
		queries[i] = query{l: l, r: r, i: iVal, k: k}
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx+1),
		queries: queries,
	}
}
