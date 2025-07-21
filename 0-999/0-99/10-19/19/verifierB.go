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

type item struct {
	t int
	c int64
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(items []item) int64 {
	n := len(items)
	const inf int64 = 1 << 60
	dp := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dp[i] = inf
	}
	for _, it := range items {
		w := it.t + 1
		for j := n; j >= 0; j-- {
			if dp[j] == inf {
				continue
			}
			nj := j + w
			if nj > n {
				nj = n
			}
			if dp[j]+it.c < dp[nj] {
				dp[nj] = dp[j] + it.c
			}
		}
	}
	return dp[n]
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	items := make([]item, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		items[i] = item{t: rng.Intn(6), c: int64(rng.Intn(50) + 1)}
		fmt.Fprintf(&sb, "%d %d\n", items[i].t, items[i].c)
	}
	expect := fmt.Sprintf("%d", solveCase(items))
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
