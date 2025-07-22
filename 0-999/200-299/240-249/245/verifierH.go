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

func precomputePal(s string) [][]bool {
	n := len(s)
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		dp[i][i] = true
	}
	for l := 2; l <= n; l++ {
		for i := 0; i+l <= n; i++ {
			j := i + l - 1
			if s[i] == s[j] {
				if l == 2 || dp[i+1][j-1] {
					dp[i][j] = true
				}
			}
		}
	}
	return dp
}

func countPal(dp [][]bool, l, r int) int {
	cnt := 0
	for i := l; i <= r; i++ {
		for j := i; j <= r; j++ {
			if dp[i][j] {
				cnt++
			}
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(3))
	}
	s := string(b)
	q := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(s)
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	dp := precomputePal(s)
	var exp strings.Builder
	for i, qr := range queries {
		if i > 0 {
			exp.WriteByte(' ')
		}
		val := countPal(dp, qr[0]-1, qr[1]-1)
		exp.WriteString(fmt.Sprint(val))
	}
	return sb.String(), exp.String()
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
