package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	in  string
	out string
}

func solveD(a, b, k, t int) int64 {
    mod := int64(1e9 + 7)
    d := b - a
    // single player's sum distribution over [-k*t, k*t]
    offset := k * t
    maxDim := 2*offset + 1
    dp := make([]int64, maxDim)
    dp[offset] = 1
    for step := 1; step <= t; step++ {
        // prefix sums of dp
        ps := make([]int64, maxDim+1)
        for i := 0; i < maxDim; i++ { ps[i+1] = (ps[i] + dp[i]) % mod }
        next := make([]int64, maxDim)
        maxS := step * k
        for s := -maxS; s <= maxS; s++ {
            idx := s + offset
            L := s - k + offset
            if L < 0 { L = 0 }
            R1 := s + k + offset + 1
            if R1 > maxDim { R1 = maxDim }
            next[idx] = (ps[R1] - ps[L] + mod) % mod
        }
        dp = next
    }
    // prefix of final dp for quick sums
    ps := make([]int64, maxDim+1)
    for i := 0; i < maxDim; i++ { ps[i+1] = (ps[i] + dp[i]) % mod }
    // sum over s_M and s_L with s_M - s_L > d  => s_L <= s_M - d - 1
    var total int64
    maxS := k * t
    for sM := -maxS; sM <= maxS; sM++ {
        idxM := sM + offset
        if dp[idxM] == 0 { continue }
        limit := sM - d - 1
        if limit < -maxS { continue }
        if limit > maxS { limit = maxS }
        sumL := ps[limit+offset+1]
        total = (total + dp[idxM]*sumL) % mod
    }
    return total
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(4))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		a := rng.Intn(20) + 1
		b := rng.Intn(20) + 1
		k := rng.Intn(5) + 1
		t := rng.Intn(3) + 1
		expect := solveD(a, b, k, t)
		in := fmt.Sprintf("%d %d %d %d\n", a, b, k, t)
		out := fmt.Sprintf("%d\n", expect)
		tests[i] = testCase{in: in, out: out}
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := strings.TrimSpace(tc.out)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
