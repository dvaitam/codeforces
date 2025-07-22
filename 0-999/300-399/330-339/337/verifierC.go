package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000009

func modPow(x, e int64) int64 {
	res := int64(1)
	x %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = (res * x) % MOD
		}
		x = (x * x) % MOD
		e >>= 1
	}
	return res
}

func solve(n, m, k int64) int64 {
	w := n - m
	safeCap := (w + 1) * (k - 1)
	if m <= safeCap {
		return m % MOD
	}
	t := m - safeCap
	b := (t + k - 1) / k
	pow2b := modPow(2, b)
	blocksScore := (2 * k) % MOD * ((pow2b - 1 + MOD) % MOD) % MOD
	rem := (m - b*k) % MOD
	if rem < 0 {
		rem += MOD
	}
	ans := (blocksScore + rem) % MOD
	return ans
}

func generateTest(rng *rand.Rand) (string, int64) {
	n := int64(rng.Intn(1000000) + 2)
	k := int64(rng.Intn(int(n-1)) + 2)
	m := int64(rng.Intn(int(n + 1)))
	inp := fmt.Sprintf("%d %d %d\n", n, m, k)
	return inp, solve(n, m, k)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest(rng)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inp)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n%s", t, err, out.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got int64
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: failed to parse output: %v\nOutput: %s\n", t, err, gotStr)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %d\nGot: %d\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
