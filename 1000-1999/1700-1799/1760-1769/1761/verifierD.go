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

const mod = 1_000_000_007

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func solve(n, k int) int {
	if k > n {
		return 0
	}
	dp0 := make([]int, k+2)
	dp1 := make([]int, k+2)
	dp0[0] = 1
	for i := 0; i < n; i++ {
		ndp0 := make([]int, k+2)
		ndp1 := make([]int, k+2)
		for j := 0; j <= k && j <= i; j++ {
			if dp0[j] != 0 {
				ndp0[j] = (ndp0[j] + 3*dp0[j]) % mod
				ndp1[j+1] = (ndp1[j+1] + dp0[j]) % mod
			}
			if dp1[j] != 0 {
				ndp0[j] = (ndp0[j] + dp1[j]) % mod
				ndp1[j+1] = (ndp1[j+1] + 3*dp1[j]) % mod
			}
		}
		dp0, dp1 = ndp0, ndp1
	}
	return (dp0[k] + dp1[k]) % mod
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n)
	input := fmt.Sprintf("%d %d\n", n, k)
	exp := fmt.Sprintf("%d", solve(n, k))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
