package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244853

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type testCase struct {
	n int
	m int
}

func modPow(a, e int64) int64 {
	a %= mod
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func solveE(n, m int) string {
	maxN := n + m
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % mod * invFact[n-k] % mod
	}

	total := comb(n+m, n)
	r := n - m
	ans := int64(0)
	for k := 1; k <= n; k++ {
		if r > k {
			ans += total
		} else {
			ans += comb(n+m, n-k)
		}
		if ans >= mod {
			ans %= mod
		}
	}
	return fmt.Sprintf("%d", ans%mod)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(6))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{{0, 0}, {1, 1}, {2, 0}, {0, 2}, {3, 3}}
	tests = append(tests, fixed...)
	for len(tests) < 100 {
		n := rng.Intn(20)
		m := rng.Intn(20)
		tests = append(tests, testCase{n: n, m: m})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.n, t.m)
		expect := solveE(t.n, t.m)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expect, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
