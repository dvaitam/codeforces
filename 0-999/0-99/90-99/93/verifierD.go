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

const MOD int64 = 1000000007

func power(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % MOD
		}
		base = (base * base) % MOD
		exp /= 2
	}
	return res
}

func geoSum(a, n int64) int64 {
	if n <= 0 {
		return 0
	}
	num := (power(3, n) - 1 + MOD) % MOD
	ans := (num * 500000004) % MOD
	ans = (ans * a) % MOD
	return ans
}

func sumF(n int64) int64 {
	if n <= 0 {
		return 0
	}
	ans := int64(4)
	ans = (ans + geoSum(4, n/2)) % MOD
	ans = (ans + geoSum(7, (n-1)/2)) % MOD

	K := (n + 1) / 2
	ans = (ans + geoSum(4, K/2)) % MOD
	ans = (ans + geoSum(7, (K-1)/2)) % MOD

	return ans
}

func compute(L, R int64) int64 {
	return (sumF(R) - sumF(L-1) + MOD) % MOD
}

type testCase struct {
	L, R int64
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.L, tc.R)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val int64
	if _, err := fmt.Fscan(strings.NewReader(strings.TrimSpace(out.String())), &val); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	expected := compute(tc.L, tc.R)
	if val != expected {
		return fmt.Errorf("expected %d got %d, input: %d %d", expected, val, tc.L, tc.R)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	cases = append(cases, testCase{L: 1, R: 1})
	cases = append(cases, testCase{L: 1, R: 3})
	cases = append(cases, testCase{L: 123, R: 12345})
	for i := 0; i < 100; i++ {
		L := rng.Int63n(1_000_000_000) + 1
		R := L + rng.Int63n(1_000_000_000-L+1)
		cases = append(cases, testCase{L: L, R: R})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
