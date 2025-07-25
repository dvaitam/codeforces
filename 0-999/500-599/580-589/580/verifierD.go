package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveCase(a []int64, bonus [][]int64, m int) int64 {
	n := len(a)
	maxMask := 1 << n
	dp := make([][]int64, maxMask)
	for i := range dp {
		dp[i] = make([]int64, n)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = a[i]
	}
	var ans int64
	for mask := 0; mask < maxMask; mask++ {
		cnt := bits.OnesCount(uint(mask))
		if cnt > m {
			continue
		}
		for last := 0; last < n; last++ {
			val := dp[mask][last]
			if val < 0 {
				continue
			}
			if cnt == m {
				if val > ans {
					ans = val
				}
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) != 0 {
					continue
				}
				newMask := mask | 1<<next
				nv := val + a[next] + bonus[last][next]
				if nv > dp[newMask][next] {
					dp[newMask][next] = nv
				}
			}
		}
	}
	return ans
}

func buildCase(a []int64, bonus [][]int64, m int) testCase {
	n := len(a)
	var sb strings.Builder
	k := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if bonus[i][j] != 0 {
				k++
			}
		}
	}
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if bonus[i][j] != 0 {
				sb.WriteString(fmt.Sprintf("%d %d %d\n", i+1, j+1, bonus[i][j]))
			}
		}
	}
	ans := solveCase(a, bonus, m)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d\n", ans)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 2
	m := rng.Intn(n) + 1
	a := make([]int64, n)
	for i := range a {
		a[i] = int64(rng.Intn(10))
	}
	bonus := make([][]int64, n)
	for i := range bonus {
		bonus[i] = make([]int64, n)
		for j := range bonus[i] {
			if i != j && rng.Intn(4) == 0 {
				bonus[i][j] = int64(rng.Intn(5))
			}
		}
	}
	return buildCase(a, bonus, m)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	a := []int64{1, 1}
	bonus := [][]int64{{0, 1}, {0, 0}}
	cases := []testCase{buildCase(a, bonus, 2)}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
