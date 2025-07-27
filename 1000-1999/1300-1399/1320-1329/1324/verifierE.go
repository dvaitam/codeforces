package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, h, l, r int
	fmt.Fscan(in, &n, &h, &l, &r)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	const negInf = -1 << 60
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, h)
		for j := range dp[i] {
			dp[i][j] = negInf
		}
	}
	dp[0][0] = 0
	for i := 1; i <= n; i++ {
		ai := a[i-1]
		for t := 0; t < h; t++ {
			if dp[i-1][t] == negInf {
				continue
			}
			nt := (t + ai) % h
			val := dp[i-1][t]
			if l <= nt && nt <= r {
				val++
			}
			if val > dp[i][nt] {
				dp[i][nt] = val
			}
			nt = (t + ai - 1) % h
			val = dp[i-1][t]
			if l <= nt && nt <= r {
				val++
			}
			if val > dp[i][nt] {
				dp[i][nt] = val
			}
		}
	}
	ans := 0
	for t := 0; t < h; t++ {
		if dp[n][t] > ans {
			ans = dp[n][t]
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	h := rng.Intn(10) + 3
	l := rng.Intn(h)
	r := l + rng.Intn(h-l)
	if r < l {
		r = l
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, h, l, r)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(rng.Intn(h-1) + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solve(in)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
