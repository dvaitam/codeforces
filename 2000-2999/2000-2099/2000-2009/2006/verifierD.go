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

type query struct {
	l int
	r int
}

type testCase struct {
	n       int
	q       int
	k       int
	b       []int
	queries []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go -- /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	expected, err := solveAllBrute(tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to compute expected answers: %v\n", err)
		os.Exit(1)
	}

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(expected))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}
	mismatch := 0
	for i := range expected {
		if expected[i] != refAns[i] {
			mismatch++
		}
	}
	if mismatch > 0 {
		fmt.Fprintf(os.Stderr, "warning: reference solution disagrees with brute answers on %d outputs, continuing with brute expectations\n", mismatch)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(expected))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on output %d: expected %d, got %d\nInput:\n%s", i+1, expected[i], candAns[i], input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2006D-ref-*")
	if err != nil {
		return "", err
	}
	_ = tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "2006D.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(target, input string) (string, error) {
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
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(k int, b []int, qs []query) {
		tc := testCase{n: len(b), q: len(qs), k: k, b: append([]int(nil), b...), queries: append([]query(nil), qs...)}
		tests = append(tests, tc)
	}

	// deterministic edge cases
	add(1, []int{1, 1}, []query{{0, 1}})
	add(5, []int{5, 5}, []query{{0, 1}})
	add(10, []int{10, 9, 9, 9, 9}, []query{{0, 4}})
	add(25, []int{12, 15, 2, 18}, []query{{0, 3}})
	add(20, []int{2, 3, 5}, []query{{0, 2}})
	add(30, []int{30, 1, 30, 1}, []query{{0, 3}, {1, 2}})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 120 {
		n := rng.Intn(6) + 2  // 2..7
		k := rng.Intn(25) + 5 // 5..29
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = rng.Intn(k) + 1
		}
		q := rng.Intn(4) + 2 // 2..5 queries
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n - 1)
			r := rng.Intn(n-l-1) + l + 1
			qs[i] = query{l: l, r: r}
		}
		add(k, b, qs)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.q, tc.k))
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", qu.l+1, qu.r+1))
		}
	}
	return sb.String()
}

func parseOutput(out string, expectedCount int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expectedCount {
		return nil, fmt.Errorf("expected %d integers, got %d", expectedCount, len(fields))
	}
	res := make([]int, expectedCount)
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}

func solveAllBrute(tests []testCase) ([]int, error) {
	ans := make([]int, 0, 512)
	for _, tc := range tests {
		for _, qu := range tc.queries {
			sub := tc.b[qu.l : qu.r+1]
			val := minChanges(sub, tc.k)
			ans = append(ans, val)
		}
	}
	return ans, nil
}

func minChanges(arr []int, k int) int {
	n := len(arr)
	if n == 0 {
		return 0
	}
	full := 1 << n
	const inf = int(1e9)
	dp := make([][]int, full)
	for i := range dp {
		dp[i] = make([]int, 2*n)
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}

	for i := 0; i < n; i++ {
		dp[1<<i][i<<1] = 0   // keep as is
		dp[1<<i][i<<1|1] = 1 // change to 1
	}

	for mask := 0; mask < full; mask++ {
		for state := 0; state < 2*n; state++ {
			cur := dp[mask][state]
			if cur == inf {
				continue
			}
			last := state >> 1
			changedLast := (state & 1) == 1
			valLast := 1
			if !changedLast {
				valLast = arr[last]
			}
			for nxt := 0; nxt < n; nxt++ {
				if mask&(1<<nxt) != 0 {
					continue
				}
				// try keeping nxt
				valNext := arr[nxt]
				if int64(valLast)*int64(valNext) <= int64(k) {
					nm := mask | (1 << nxt)
					idx := nxt << 1
					if cur < dp[nm][idx] {
						dp[nm][idx] = cur
					}
				}
				// try changing nxt to 1
				if int64(valLast) <= int64(k) { // valNext == 1
					nm := mask | (1 << nxt)
					idx := nxt<<1 | 1
					if cur+1 < dp[nm][idx] {
						dp[nm][idx] = cur + 1
					}
				}
			}
		}
	}

	fullMask := full - 1
	ans := inf
	for state := 0; state < 2*n; state++ {
		if dp[fullMask][state] < ans {
			ans = dp[fullMask][state]
		}
	}
	return ans
}
