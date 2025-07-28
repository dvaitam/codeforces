package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const MOD int64 = 1_000_000_007

func solveCase(s, t string) (int, int64) {
	L := len(t)
	positions := make([]int, 0)
	for i := 0; i+L <= len(s); i++ {
		if s[i:i+L] == t {
			positions = append(positions, i)
		}
	}
	m := len(positions)
	if m == 0 {
		return 0, 1
	}
	right := make([]int, m)
	j := 0
	for i := 0; i < m; i++ {
		if j < i {
			j = i
		}
		for j+1 < m && positions[j+1] <= positions[i]+L-1 {
			j++
		}
		right[i] = j
	}
	next := make([]int, m)
	for i := 0; i < m; i++ {
		x := positions[i] + L
		next[i] = sort.Search(len(positions), func(p int) bool { return positions[p] >= x })
	}
	const INF = int(1 << 30)
	dp := make([]int, m+1)
	cnt := make([]int64, m+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[m] = 0
	cnt[m] = 1
	for i := m - 1; i >= 0; i-- {
		r := right[i]
		best := INF
		var bestCnt int64
		for k := i; k <= r; k++ {
			j := next[k]
			cand := 1 + dp[j]
			if cand < best {
				best = cand
				bestCnt = cnt[j]
			} else if cand == best {
				bestCnt += cnt[j]
				if bestCnt >= MOD {
					bestCnt %= MOD
				}
			}
		}
		dp[i] = best
		cnt[i] = bestCnt % MOD
	}
	return dp[0], cnt[0]
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	type test struct{ s, t string }
	tests := []test{{s: "abc", t: "a"}, {s: "aaaa", t: "aa"}}
	for len(tests) < 120 {
		n := rng.Intn(8) + 2
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('a' + rng.Intn(3)))
		}
		s := sb.String()
		m := rng.Intn(3) + 1
		if m > n {
			m = n
		}
		var sb2 strings.Builder
		for i := 0; i < m; i++ {
			sb2.WriteByte(byte('a' + rng.Intn(3)))
		}
		t := sb2.String()
		tests = append(tests, test{s: s, t: t})
	}

	for i, tc := range tests {
		moves, cnt := solveCase(tc.s, tc.t)
		input := fmt.Sprintf("1\n%s\n%s\n", tc.s, tc.t)
		expected := fmt.Sprintf("%d %d", moves, cnt)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
