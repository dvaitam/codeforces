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

func modPow(a, b int64) int64 {
	a %= MOD
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func solveC(n, k int64) int64 {
	if n%2 == 1 {
		base := (modPow(2, n-1) + 1) % MOD
		return modPow(base, k)
	}
	a := (modPow(2, n-1) - 1 + MOD) % MOD
	pow2n := modPow(2, n)
	ans := int64(1)
	cur := int64(1)
	for i := int64(1); i <= k; i++ {
		ans = (a*ans + cur) % MOD
		cur = cur * pow2n % MOD
	}
	return ans
}

func generateC(rng *rand.Rand) (int64, int64) {
	n := int64(rng.Intn(30) + 1)
	k := int64(rng.Intn(30) + 1)
	return n, k
}

func runCase(bin string, n, k int64) (string, error) {
	input := fmt.Sprintf("1\n%d %d\n", n, k)
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	type pair struct{ n, k int64 }
	cases := make([]pair, 0, 102)
	cases = append(cases, pair{1, 1})
	cases = append(cases, pair{2, 2})
	for i := 0; i < 100; i++ {
		n, k := generateC(rng)
		cases = append(cases, pair{n, k})
	}
	for i, tc := range cases {
		expect := solveC(tc.n, tc.k)
		out, err := runCase(bin, tc.n, tc.k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
