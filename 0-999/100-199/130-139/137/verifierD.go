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

func compute(s string, k int) (int, string) {
	n := len(s)
	cost := make([][]int, n)
	for i := range cost {
		cost[i] = make([]int, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			c := 0
			if s[i] != s[j] {
				c = 1
			}
			if i+1 <= j-1 {
				cost[i][j] = cost[i+1][j-1] + c
			} else {
				cost[i][j] = c
			}
		}
	}
	const inf = 1_000_000_000
	dp := make([][]int, k+1)
	prev := make([][]int, k+1)
	for p := 0; p <= k; p++ {
		dp[p] = make([]int, n)
		prev[p] = make([]int, n)
		for i := 0; i < n; i++ {
			dp[p][i] = inf
			prev[p][i] = -1
		}
	}
	for i := 0; i < n; i++ {
		dp[1][i] = cost[0][i]
		prev[1][i] = -1
	}
	for p := 2; p <= k; p++ {
		for i := p - 1; i < n; i++ {
			for t := p - 2; t < i; t++ {
				cur := dp[p-1][t] + cost[t+1][i]
				if cur < dp[p][i] {
					dp[p][i] = cur
					prev[p][i] = t
				}
			}
		}
	}
	best := inf
	bp := 1
	for p := 1; p <= k; p++ {
		if dp[p][n-1] < best {
			best = dp[p][n-1]
			bp = p
		}
	}
	segments := make([][2]int, 0, bp)
	p := bp
	i := n - 1
	for p > 0 {
		t := prev[p][i]
		l := 0
		if t >= 0 {
			l = t + 1
		}
		segments = append(segments, [2]int{l, i})
		i = t
		p--
	}
	for l, r := 0, len(segments)-1; l < r; l, r = l+1, r-1 {
		segments[l], segments[r] = segments[r], segments[l]
	}
	runes := []rune(s)
	for _, seg := range segments {
		l, r := seg[0], seg[1]
		for a, b := l, r; a < b; a, b = a+1, b-1 {
			if runes[a] != runes[b] {
				runes[b] = runes[a]
			}
		}
	}
	out := make([]rune, 0, n+len(segments)-1)
	for idx, seg := range segments {
		if idx > 0 {
			out = append(out, '+')
		}
		for j := seg[0]; j <= seg[1]; j++ {
			out = append(out, runes[j])
		}
	}
	return best, string(out)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	k := rng.Intn(n) + 1
	expCost, expStr := compute(s, k)
	var sb strings.Builder
	sb.WriteString(s)
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", k)
	input := sb.String()
	exp := fmt.Sprintf("%d\n%s", expCost, expStr)
	return input, exp
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
