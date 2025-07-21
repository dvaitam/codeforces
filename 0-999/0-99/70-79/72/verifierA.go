package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n int) string {
	// generate primes up to n
	isPrime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if isPrime[i] {
			for j := i * i; j <= n; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := make([]int, 0)
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	vals := make([]int, 0, len(primes)+1)
	for i := len(primes) - 1; i >= 0; i-- {
		vals = append(vals, primes[i])
	}
	if n >= 1 {
		vals = append(vals, 1)
	}
	m := len(vals)
	dp := make([][]bool, m+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[m][0] = true
	for i := m - 1; i >= 0; i-- {
		vi := vals[i]
		for s := 0; s <= n; s++ {
			if dp[i+1][s] {
				dp[i][s] = true
			} else if s >= vi && dp[i+1][s-vi] {
				dp[i][s] = true
			}
		}
	}
	if !dp[0][n] {
		return "0"
	}
	res := make([]int, 0)
	sum := n
	idx := 0
	for sum > 0 {
		picked := false
		for i := idx; i < m; i++ {
			v := vals[i]
			if v > sum {
				continue
			}
			if dp[i+1][sum-v] {
				res = append(res, v)
				sum -= v
				idx = i + 1
				picked = true
				break
			}
		}
		if !picked {
			return "0"
		}
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10000) + 1
	input := fmt.Sprintf("%d\n", n)
	return input, expected(n)
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
		return fmt.Errorf("expected '%s' got '%s'", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
