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

const MOD int64 = 1e9 + 7

type edge struct {
	to  int
	val int
}

func solveD(s string) int64 {
	n := len(s)
	edges := make([][]edge, n)
	for i := 0; i < n; i++ {
		val := 0
		for j := i; j < n; j++ {
			val = val*2 + int(s[j]-'0')
			if val > 20 {
				break
			}
			if val > 0 {
				edges[i] = append(edges[i], edge{j + 1, val})
			}
		}
	}
	dp := make([]map[int]int64, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make(map[int]int64)
	}
	for i := 0; i < n; i++ {
		dp[i][0] = (dp[i][0] + 1) % MOD
	}
	for i := 0; i < n; i++ {
		for mask, cnt := range dp[i] {
			for _, e := range edges[i] {
				nm := mask | (1 << (e.val - 1))
				dp[e.to][nm] = (dp[e.to][nm] + cnt) % MOD
			}
		}
	}
	var ans int64
	for i := 0; i <= n; i++ {
		for mask, cnt := range dp[i] {
			if mask != 0 && (mask&(mask+1)) == 0 {
				ans = (ans + cnt) % MOD
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(12) + 1
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	s := sb.String()
	expected := solveD(s)
	return fmt.Sprintf("%d\n%s\n", n, s), fmt.Sprintf("%d", expected)
}

func runCase(exe, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := runCase(exe, in)
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
