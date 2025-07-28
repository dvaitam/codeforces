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

var fact [2001]int64
var invfact [2001]int64
var pow2 [2001]int64
var invPow2 [2001]int64

func powmod(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initComb() {
	fact[0] = 1
	for i := 1; i <= 2000; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invfact[2000] = powmod(fact[2000], MOD-2)
	for i := 2000; i >= 1; i-- {
		invfact[i-1] = invfact[i] * int64(i) % MOD
	}
	pow2[0] = 1
	for i := 1; i <= 2000; i++ {
		pow2[i] = pow2[i-1] * 2 % MOD
	}
	invPow2[2000] = powmod(pow2[2000], MOD-2)
	for i := 2000; i >= 1; i-- {
		invPow2[i-1] = invPow2[i] * 2 % MOD
	}
}

func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invfact[k] % MOD * invfact[n-k] % MOD
}

func solve(n, m int, k int64) int64 {
	sum := int64(0)
	for x := 0; x < m; x++ {
		sum = (sum + int64(m-x)*C(n, x)) % MOD
	}
	ans := k % MOD * sum % MOD * invPow2[n-1] % MOD
	return ans
}

func buildInput(n, m int, k int64) string {
	return fmt.Sprintf("1\n%d %d %d\n", n, m, k)
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
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	initComb()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(n) + 1
		k := int64(rng.Intn(5) + 1)
		input := buildInput(n, m, k)
		exp := fmt.Sprintf("%d", solve(n, m, k))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
