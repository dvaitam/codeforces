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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expectedC(arr []int) int {
	n := len(arr)
	const V = 20
	offset := V
	const INF = int(1e9)
	dp := make([]int, 2*V+1)
	for i := range dp {
		dp[i] = abs(arr[0] - (i - offset))
	}
	for i := 1; i < n; i++ {
		ndp := make([]int, 2*V+1)
		for j := range ndp {
			ndp[j] = INF
		}
		best := INF
		for v := -V; v <= V; v++ {
			idx := v + offset
			if dp[idx] < best {
				best = dp[idx]
			}
			cost := best + abs(arr[i]-v)
			ndp[idx] = cost
		}
		copy(dp, ndp)
	}
	ans := INF
	for _, v := range dp {
		if v < ans {
			ans = v
		}
	}
	return ans
}

func generateCaseC(rng *rand.Rand) (string, int) {
	n := rng.Intn(6) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(21) - 10
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), expectedC(arr)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
