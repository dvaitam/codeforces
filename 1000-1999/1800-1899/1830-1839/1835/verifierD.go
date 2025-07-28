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

type edge struct{ from, to int }

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

func bruteD(n int, edges []edge, k int) int64 {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	for _, e := range edges {
		adj[e.from][e.to] = true
	}
	dp := make([][][]bool, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([][]bool, n)
		for j := 0; j < n; j++ {
			dp[i][j] = make([]bool, n)
		}
	}
	for i := 0; i < n; i++ {
		dp[0][i][i] = true
	}
	for step := 1; step <= k; step++ {
		for i := 0; i < n; i++ {
			for mid := 0; mid < n; mid++ {
				if dp[step-1][i][mid] {
					for j := 0; j < n; j++ {
						if adj[mid][j] {
							dp[step][i][j] = true
						}
					}
				}
			}
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		if dp[k][i][i] {
			ans++
		}
		for j := i + 1; j < n; j++ {
			if dp[k][i][j] && dp[k][j][i] {
				ans++
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges + 1)
	k := rng.Intn(4) + 1
	seen := make(map[[2]int]struct{})
	var es []edge
	for len(es) < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		es = append(es, edge{u, v})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for _, e := range es {
		fmt.Fprintf(&sb, "%d %d\n", e.from+1, e.to+1)
	}
	input := sb.String()
	val := bruteD(n, es, k)
	expected := fmt.Sprintf("%d", val)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
