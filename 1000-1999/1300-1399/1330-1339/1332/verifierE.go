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
	a %= mod
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveCase(n, m, L, R int64) int64 {
	N := n * m
	total := R - L + 1
	if N%2 == 1 {
		return modPow(total%mod, N)
	}
	even := R/2 - (L-1)/2
	odd := total - even
	diff := (even - odd) % mod
	if diff < 0 {
		diff += mod
	}
	a := modPow(total%mod, N)
	b := modPow(diff, N)
	ans := (a + b) % mod
	ans = ans * ((mod + 1) / 2) % mod
	return ans
}

func generateCase(rng *rand.Rand) (int64, int64, int64, int64) {
	n := int64(rng.Intn(4) + 1)
	m := int64(rng.Intn(4) + 1)
	L := int64(rng.Intn(10))
	R := L + int64(rng.Intn(10))
	return n, m, L, R
}

func runCase(bin string, n, m, L, R int64) error {
	input := fmt.Sprintf("%d %d %d %d\n", n, m, L, R)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := solveCase(n, m, L, R)
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, L, R := generateCase(rng)
		if err := runCase(bin, n, m, L, R); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
