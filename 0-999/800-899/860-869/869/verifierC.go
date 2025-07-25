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

const modC = 998244353

func modPow(x, y int) int {
	res := 1
	x %= modC
	for y > 0 {
		if y&1 == 1 {
			res = res * x % modC
		}
		x = x * x % modC
		y >>= 1
	}
	return res
}

func solveC(a, b, c int) int {
	maxN := a
	if b > maxN {
		maxN = b
	}
	if c > maxN {
		maxN = c
	}
	fact := make([]int, maxN+1)
	invFact := make([]int, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * i % modC
	}
	invFact[maxN] = modPow(fact[maxN], modC-2)
	for i := maxN; i >= 1; i-- {
		invFact[i-1] = invFact[i] * i % modC
	}
	comb := func(n, k int) int {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % modC * invFact[n-k] % modC
	}
	F := func(x, y int) int {
		limit := x
		if y < limit {
			limit = y
		}
		res := 0
		for i := 0; i <= limit; i++ {
			add := comb(x, i)
			add = add * comb(y, i) % modC
			add = add * fact[i] % modC
			res += add
			if res >= modC {
				res -= modC
			}
		}
		return res
	}
	ans := F(a, b)
	ans = ans * F(a, c) % modC
	ans = ans * F(b, c) % modC
	return ans
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a := rng.Intn(20) + 1
		b := rng.Intn(20) + 1
		c := rng.Intn(20) + 1
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
		expected := fmt.Sprintf("%d", solveC(a, b, c))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
