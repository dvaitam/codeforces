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
	arr []int64
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{arr: []int64{1, 1, 1}},
		{arr: []int64{5, 2, 5, 2}},
		{arr: []int64{2, 2, 5, 9, 5}},
		{arr: []int64{1, 2, 3, 4}},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 3
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(10) + 1)
	}
	return testCase{arr: arr}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(len(tc.arr)))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, count int) ([]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != count {
		return nil, fmt.Errorf("expected %d outputs, got %d", count, len(lines))
	}
	ans := make([]int64, count)
	for i, line := range lines {
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", line)
		}
		ans[i] = val
	}
	return ans, nil
}

func bruteForce(tc testCase) int64 {
	n := len(tc.arr)
	best := int64(1 << 62)
	var dfs func(int, []int64)
	dfs = func(idx int, cur []int64) {
		if idx == n {
			if len(cur) == 0 {
				return
			}
			ok := true
			for i := 0; i < len(cur); i++ {
				l := (i - 1 + len(cur)) % len(cur)
				r := (i + 1) % len(cur)
				if cur[i] != cur[l] && cur[i] != cur[r] {
					ok = false
					break
				}
			}
			if ok {
				cost := int64(0)
				for i := 0; i < n; i++ {
					cost += abs(cur[i] - tc.arr[i])
				}
				if cost < best {
					best = cost
				}
			}
			return
		}
		for delta := -3; delta <= 3; delta++ {
			cur[idx] = tc.arr[idx] + int64(delta)
			dfs(idx+1, cur)
		}
	}
	dfs(0, make([]int64, n))
	if best == int64(1<<62) {
		return -1
	}
	return best
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	input := buildInput(tests)
	out, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	results, err := parseOutput(out, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := bruteForce(tc)
		if results[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\nn=%d arr=%v\n", i+1, exp, results[i], len(tc.arr), tc.arr)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
