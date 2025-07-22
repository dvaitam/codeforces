package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func maxInt16(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func solveExpected(s, p string) []int {
	n := len(s)
	m := len(p)
	match := make([]bool, n)
	for i := 0; i+m <= n; i++ {
		ok := true
		for j := 0; j < m; j++ {
			if s[i+j] != p[j] {
				ok = false
				break
			}
		}
		match[i] = ok
	}
	const negInf int16 = -10000
	dp := make([][]int16, n+1)
	for i := range dp {
		dp[i] = make([]int16, n+1)
		for j := range dp[i] {
			dp[i][j] = negInf
		}
	}
	dp[0][0] = 0
	for i := 0; i <= n; i++ {
		for j := 0; j <= n; j++ {
			cur := dp[i][j]
			if cur < 0 {
				continue
			}
			if i < n {
				if j+1 <= n {
					dp[i+1][j+1] = maxInt16(dp[i+1][j+1], cur)
				}
				dp[i+1][j] = maxInt16(dp[i+1][j], cur)
			}
			if i+m <= n && match[i] {
				dp[i+m][j] = maxInt16(dp[i+m][j], cur+1)
			}
		}
	}
	ans := make([]int, n+1)
	best := int16(0)
	for x := 0; x <= n; x++ {
		if dp[n][x] > best {
			best = dp[n][x]
		}
		ans[x] = int(best)
	}
	return ans
}

type testCase struct {
	s string
	p string
}

func randStr(rng *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rng.Intn(3))
	}
	return string(b)
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	m := rng.Intn(min(n, 5)) + 1
	s := randStr(rng, n)
	p := randStr(rng, m)
	return testCase{s, p}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%s\n%s\n", tc.s, tc.p)
	exp := solveExpected(tc.s, tc.p)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		var val int
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad number %q", f)
		}
		if val != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{"aaaaa", "aa"}, {"abababa", "aba"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
