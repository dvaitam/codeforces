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

func runCandidate(bin string, input string) (string, error) {
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

func solve(names []string) int64 {
	const K = 26
	const INF int64 = 1e18
	dp := make([][K]int64, K)
	for i := 0; i < K; i++ {
		for j := 0; j < K; j++ {
			dp[i][j] = -INF
		}
	}
	var ans int64
	for _, s := range names {
		if len(s) == 0 {
			continue
		}
		u := int(s[0] - 'a')
		v := int(s[len(s)-1] - 'a')
		w := int64(len(s))
		for st := 0; st < K; st++ {
			if dp[st][u] > -INF {
				if dp[st][u]+w > dp[st][v] {
					dp[st][v] = dp[st][u] + w
				}
			}
		}
		if w > dp[u][v] {
			dp[u][v] = w
		}
		if dp[u][u] > ans {
			ans = dp[u][u]
		}
		if dp[v][v] > ans {
			ans = dp[v][v]
		}
	}
	if ans < 0 {
		ans = 0
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	names := make([]string, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(5) + 1
		b := make([]byte, l)
		for j := 0; j < l; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		names[i] = string(b)
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d\n", n)
	for _, s := range names {
		fmt.Fprintf(&in, "%s\n", s)
	}
	exp := fmt.Sprintf("%d", solve(names))
	return in.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
