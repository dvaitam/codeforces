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

func expected(n, d int, c []int) string {
	sort.Ints(c)
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + c[i-1]
	}
	const INF = int(1e9)
	dp := make([]int, n+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] < INF && 2*prefix[j]+d >= prefix[i] {
				if dp[j]+1 < dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}
	}
	bestSum, bestDays := 0, 0
	for i := 0; i <= n; i++ {
		if dp[i] < INF {
			if prefix[i] > bestSum {
				bestSum = prefix[i]
				bestDays = dp[i]
			} else if prefix[i] == bestSum && dp[i] < bestDays {
				bestDays = dp[i]
			}
		}
	}
	return fmt.Sprintf("%d %d", bestSum, bestDays)
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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	d := rng.Intn(5) + 1
	c := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, d)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c[i])
	}
	sb.WriteByte('\n')
	exp := expected(n, d, c)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
