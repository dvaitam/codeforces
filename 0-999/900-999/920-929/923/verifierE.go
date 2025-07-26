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

const mod int64 = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func inv(a int64) int64 {
	return modPow(a, mod-2)
}

func step(dp []int64, invs []int64) []int64 {
	n := len(dp) - 1
	prefix := int64(0)
	res := make([]int64, n+1)
	for i := n; i >= 0; i-- {
		prefix = (prefix + dp[i]*invs[i]) % mod
		res[i] = prefix
	}
	return res
}

func solveE(n, m int, p []int64) []int64 {
	invs := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		invs[i] = inv(int64(i + 1))
	}
	dp := make([]int64, n+1)
	copy(dp, p)
	for stepNum := 0; stepNum < m; stepNum++ {
		dp = step(dp, invs)
	}
	for i := range dp {
		dp[i] = (dp[i]%mod + mod) % mod
	}
	return dp
}

func genTest(rng *rand.Rand) (int, int, []int64) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	p := make([]int64, n+1)
	for i := range p {
		p[i] = int64(rng.Intn(100)) % mod
	}
	return n, m, p
}

func formatInput(n, m int, p []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, p := genTest(rng)
		input := formatInput(n, m, p)
		expVals := solveE(n, m, p)
		var exp strings.Builder
		for j, v := range expVals {
			if j > 0 {
				exp.WriteByte(' ')
			}
			fmt.Fprintf(&exp, "%d", v)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp.String() {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp.String(), got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
