package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesE = `5 3 1 4
9 3 2 1
2 2 1 2
10 9 6 3
4 1 3 2
2 2 1 2
6 3 6 3
3 3 2 4
10 4 3 2
9 5 2 5
6 1 3 5
6 5 2 4
8 5 7 4
4 2 3 3
2 1 2 3
10 9 8 3
4 2 1 4
5 4 3 2
7 4 6 5
7 6 5 2
7 1 7 1
5 3 5 5
5 1 3 2
6 4 1 1
7 6 1 3
7 1 3 3
7 2 7 4
3 2 3 2
9 5 3 3
8 3 6 5
2 2 1 4
4 3 1 4
5 4 2 1
2 1 2 5
5 3 1 1
10 5 7 2
9 4 8 4
9 1 4 4
9 4 7 2
9 4 1 1
6 3 2 5
5 2 4 3
4 3 1 3
3 3 2 1
9 7 2 4
5 5 2 3
6 6 4 3
8 4 5 3
8 8 2 3
5 1 4 5
4 3 1 2
9 8 7 4
5 1 2 2
2 2 1 4
8 4 1 2
4 3 4 5
9 1 2 1
3 2 3 3
4 1 3 1
10 1 5 3
3 1 3 4
8 4 5 4
5 4 1 1
3 3 2 5
8 7 8 1
5 3 4 4
3 3 1 3
4 2 3 4
7 3 5 1
4 1 3 1
10 2 8 5
9 9 2 5
5 4 3 3
5 2 1 1
7 5 4 5
6 5 4 5
9 7 3 3
7 6 3 2
8 2 3 5
4 3 2 5
7 6 5 1
3 2 3 2
7 6 3 3
4 3 1 5
2 1 2 1
4 2 4 5
3 1 3 4
5 5 3 4
5 4 2 3
7 2 6 3
10 9 8 4
10 7 6 3
9 7 1 3
4 4 3 4
8 1 3 5
3 2 3 3
6 1 3 4
5 4 5 1
4 2 1 3
4 1 2 4`

const mod = 1000000007

// Embedded solver from 479E.go.
func countWays(n, a, b, k int) int {
	dp := make([]int, n+1)
	dp[a] = 1
	prefix := make([]int, n+1)
	newDp := make([]int, n+1)

	for step := 0; step < k; step++ {
		current := 0
		for i := 1; i <= n; i++ {
			current += dp[i]
			if current >= mod {
				current -= mod
			}
			prefix[i] = current
		}
		for i := range newDp {
			newDp[i] = 0
		}

		if a < b {
			for y := 1; y < b; y++ {
				limit := (b + y - 1) / 2
				val := prefix[limit]
				if y <= limit {
					val -= dp[y]
					if val < 0 {
						val += mod
					}
				}
				newDp[y] = val
			}
		} else {
			for y := b + 1; y <= n; y++ {
				start := (b+y)/2 + 1
				val := prefix[n] - prefix[start-1]
				if val < 0 {
					val += mod
				}
				if y >= start {
					val -= dp[y]
					if val < 0 {
						val += mod
					}
				}
				newDp[y] = val
			}
		}
		dp, newDp = newDp, dp
	}

	ans := 0
	for _, v := range dp {
		ans += v
		if ans >= mod {
			ans -= mod
		}
	}
	return ans
}

type testCase struct {
	n int
	a int
	b int
	k int
}

func parseCases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesE), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 integers", idx+1)
		}
		n, err1 := strconv.Atoi(fields[0])
		a, err2 := strconv.Atoi(fields[1])
		b, err3 := strconv.Atoi(fields[2])
		k, err4 := strconv.Atoi(fields[3])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return nil, fmt.Errorf("line %d: parse error", idx+1)
		}
		cases = append(cases, testCase{n: n, a: a, b: b, k: k})
	}
	return cases, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d %d %d %d\n", tc.n, tc.a, tc.b, tc.k)
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := countWays(tc.n, tc.a, tc.b, tc.k)
		input := buildInput(tc)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output\n", idx+1)
			os.Exit(1)
		}
		if got%mod != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
