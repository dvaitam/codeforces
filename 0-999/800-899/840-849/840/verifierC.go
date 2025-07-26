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

const MOD int = 1_000_000_007

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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func squareFree(x int64) int64 {
	var res int64 = 1
	for p := int64(2); p*p <= x; p++ {
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		if cnt%2 == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func expected(a []int64) int {
	n := len(a)
	sf := make([]int64, n)
	for i := range a {
		sf[i] = squareFree(a[i])
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			if i != j && sf[i] != sf[j] {
				adj[i][j] = true
			}
		}
	}
	full := 1 << n
	dp := make([][]int, full)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		dp[1<<i][i] = 1
	}
	for mask := 0; mask < full; mask++ {
		for last := 0; last < n; last++ {
			val := dp[mask][last]
			if val == 0 {
				continue
			}
			for next := 0; next < n; next++ {
				if mask&(1<<next) == 0 && adj[last][next] {
					nm := mask | (1 << next)
					dp[nm][next] = (dp[nm][next] + val) % MOD
				}
			}
		}
	}
	ans := 0
	finalMask := full - 1
	for i := 0; i < n; i++ {
		ans = (ans + dp[finalMask][i]) % MOD
	}
	return ans
}

func genCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(6) + 2 // up to 7-8 maybe
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(50) + 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(a[i], 10))
	}
	sb.WriteByte('\n')
	return sb.String(), a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, arr := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "output: %s\n", out)
			os.Exit(1)
		}
		exp := expected(arr)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
