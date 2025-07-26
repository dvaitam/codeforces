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

const MOD = 998244353

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

func solveCase(n, k int) int {
	dp0 := make([][]int, n+1)
	dp1 := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp0[i] = make([]int, k+1)
		dp1[i] = make([]int, k+1)
	}
	if k >= 2 {
		dp0[1][2] = 1
	}
	if k >= 1 {
		dp1[1][1] = 1
	}
	for i := 2; i <= n; i++ {
		for j := 1; j <= k; j++ {
			v0 := dp0[i-1][j]
			if j-1 >= 0 {
				v0 += 2 * dp1[i-1][j-1]
			}
			if j-2 >= 0 {
				v0 += dp0[i-1][j-2]
			}
			dp0[i][j] = v0 % MOD
			v1 := dp1[i-1][j]
			if j-1 >= 0 {
				v1 += dp1[i-1][j-1]
			}
			v1 += 2 * dp0[i-1][j]
			dp1[i][j] = v1 % MOD
		}
	}
	res := (dp0[n][k] + dp1[n][k]) * 2 % MOD
	return res
}

func genCase(rng *rand.Rand) (int, int) {
	n := rng.Intn(6) + 1
	k := rng.Intn(2*n) + 1
	return n, k
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k := genCase(rng)
		input := fmt.Sprintf("%d %d\n", n, k)
		expect := solveCase(n, k)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(out, &got); err != nil || got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
