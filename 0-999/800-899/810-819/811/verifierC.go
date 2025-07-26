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

const MaxA = 5000

func solveCase(a []int) int {
	n := len(a) - 1
	first := make([]int, MaxA+1)
	last := make([]int, MaxA+1)
	for i := 0; i <= MaxA; i++ {
		first[i] = n + 1
		last[i] = 0
	}
	for i := 1; i <= n; i++ {
		v := a[i]
		if first[v] == n+1 {
			first[v] = i
		}
		last[v] = i
	}
	dp := make([]int, n+1)
	visited := make([]bool, MaxA+1)
	for l := 1; l <= n; l++ {
		for i := 0; i <= MaxA; i++ {
			visited[i] = false
		}
		xorVal := 0
		r := l - 1
		for j := l; j <= n; j++ {
			v := a[j]
			if first[v] < l {
				break
			}
			if !visited[v] {
				visited[v] = true
				xorVal ^= v
			}
			if last[v] > r {
				r = last[v]
			}
			if j == r {
				if cand := dp[l-1] + xorVal; cand > dp[r] {
					dp[r] = cand
				}
			}
		}
	}
	for i := 1; i <= n; i++ {
		if dp[i] < dp[i-1] {
			dp[i] = dp[i-1]
		}
	}
	return dp[n]
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(15) + 1
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = rng.Intn(20)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d", solveCase(arr))
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
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
