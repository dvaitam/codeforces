package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type state struct {
	val string
	len int
	ok  bool
}

// compareNumbers returns -1 if a<b, 0 if equal, 1 if a>b (decimal strings without sign)
func compareNumbers(a, b string) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return 1
	}
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// solve computes the minimal digits to remove to achieve minimal possible cost.
func solve(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	maxSum := 9 * n
	dp := make([]state, maxSum+1)
	dp[0] = state{val: "", len: 0, ok: true}

	for _, ch := range s {
		d := int(ch - '0')
		next := make([]state, maxSum+1)
		copy(next, dp) // skipping current digit keeps existing best
		for sum, st := range dp {
			if !st.ok {
				continue
			}
			ns := sum + d
			if ns > maxSum {
				continue
			}
			candVal := st.val + string(ch)
			candLen := st.len + 1
			if !next[ns].ok {
				next[ns] = state{val: candVal, len: candLen, ok: true}
				continue
			}
			comp := compareNumbers(candVal, next[ns].val)
			if comp < 0 || (comp == 0 && candLen > next[ns].len) {
				next[ns] = state{val: candVal, len: candLen, ok: true}
			}
		}
		dp = next
	}

	var bestVal string
	bestSum := 0
	bestKeep := 0
	for sum := 1; sum <= maxSum; sum++ {
		if !dp[sum].ok {
			continue
		}
		valStr := dp[sum].val
		if bestSum == 0 {
			bestSum = sum
			bestVal = valStr
			bestKeep = dp[sum].len
			continue
		}
		// compare val/sum with bestVal/bestSum via cross multiplication
		var left, right big.Int
		left.SetString(valStr, 10)
		right.SetString(bestVal, 10)
		left.Mul(&left, big.NewInt(int64(bestSum)))
		right.Mul(&right, big.NewInt(int64(sum)))
		cmp := left.Cmp(&right)
		if cmp < 0 {
			bestSum = sum
			bestVal = valStr
			bestKeep = dp[sum].len
		} else if cmp == 0 && dp[sum].len > bestKeep {
			bestSum = sum
			bestVal = valStr
			bestKeep = dp[sum].len
		}
	}
	return n - bestKeep
}

type testCase struct {
	val string
}

func genCase(rng *rand.Rand) testCase {
	// length between 1 and 50
	len := rng.Intn(50) + 1
	b := make([]byte, len)
	for i := 0; i < len; i++ {
		if i == 0 {
			b[i] = byte(rng.Intn(9)+1) + '0'
		} else {
			b[i] = byte(rng.Intn(10)) + '0'
		}
	}
	return testCase{val: string(b)}
}

func buildInput(cases []testCase) (string, []int) {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cases))
	exp := make([]int, len(cases))
	for i, tc := range cases {
		fmt.Fprintln(&sb, tc.val)
		exp[i] = solve(tc.val)
	}
	return sb.String(), exp
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/2093B_binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	for i := 0; i < 200; i++ {
		cases = append(cases, genCase(rng))
	}
	input, expected := buildInput(cases)
	output, err := runCandidate(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run candidate: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Fields(output)
	if len(lines) < len(expected) {
		fmt.Fprintf(os.Stderr, "not enough outputs: got %d expected %d\n", len(lines), len(expected))
		os.Exit(1)
	}
	for i, exp := range expected {
		got := strings.TrimSpace(lines[i])
		if got != fmt.Sprintf("%d", exp) {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %d got %s\ninput: %s\n", i+1, exp, got, cases[i].val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
