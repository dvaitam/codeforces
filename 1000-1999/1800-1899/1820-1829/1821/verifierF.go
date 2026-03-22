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

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= MOD
	for b > 0 {
		if b%2 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		b /= 2
	}
	return res
}

// Embedded solver for 1821F
func solveF(n, m, k int64) string {
	if m*(k+1) > n {
		return "0"
	}

	D := n - m*(k+1)

	fact := make([]int64, n+1)
	invFact := make([]int64, n+1)
	fact[0] = 1
	invFact[0] = 1
	for i := int64(1); i <= n; i++ {
		fact[i] = (fact[i-1] * i) % MOD
	}
	invFact[n] = powMod(fact[n], MOD-2)
	for i := n - 1; i >= 1; i-- {
		invFact[i] = (invFact[i+1] * (i + 1)) % MOD
	}

	nCr := func(nn, r int64) int64 {
		if r < 0 || r > nn {
			return 0
		}
		return fact[nn] * invFact[r] % MOD * invFact[nn-r] % MOD
	}

	ans := int64(0)
	limit := D / k
	if m < limit {
		limit = m
	}

	for j := int64(0); j <= limit; j++ {
		term := nCr(m, j)
		term = (term * powMod(2, m-j)) % MOD
		term = (term * nCr(D-j*k+m, m)) % MOD

		if j%2 == 1 {
			ans = (ans - term + MOD) % MOD
		} else {
			ans = (ans + term) % MOD
		}
	}

	return fmt.Sprintf("%d", ans)
}

func genCaseF(rng *rand.Rand) (int, int, int) {
	n := rng.Intn(20) + 1
	m := rng.Intn(n) + 1
	k := rng.Intn(n) + 1
	return n, m, k
}

func runBin(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, k := genCaseF(rng)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		expect := solveF(int64(n), int64(m), int64(k))
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput: %s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
