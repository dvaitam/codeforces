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
const maxN = 100

var fact [maxN + 1]int64
var invFact [maxN + 1]int64
var bell [maxN + 1]int64

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func initPrecalc() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	invFact[maxN] = modPow(fact[maxN], MOD-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % MOD
	}

	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}

	bell[0] = 1
	for i := 0; i < maxN; i++ {
		sum := int64(0)
		for k := 0; k <= i; k++ {
			sum = (sum + comb(i, k)*bell[k]) % MOD
		}
		bell[i+1] = sum
	}
}

func solveCase(n int) string {
	comb := func(n, k int) int64 {
		if k < 0 || k > n {
			return 0
		}
		return fact[n] * invFact[k] % MOD * invFact[n-k] % MOD
	}
	ans := int64(0)
	for m := 0; m < n; m++ {
		ans = (ans + comb(n, m)*bell[m]) % MOD
	}
	return fmt.Sprintf("%d\n", ans%MOD)
}

type testCase struct {
	n        int
	expected string
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(maxN-1) + 1
	return testCase{n: n, expected: solveCase(n)}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d\n", tc.n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(tc.expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(tc.expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initPrecalc()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	edges := []int{1, 2, 3, 10, 20, 50, 100}
	for _, e := range edges {
		if e <= maxN {
			cases = append(cases, testCase{n: e, expected: solveCase(e)})
		}
	}

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
