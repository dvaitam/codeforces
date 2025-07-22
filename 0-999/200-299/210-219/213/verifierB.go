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

const mod = 1000000007

func modPow(a, e int64) int64 {
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

type caseB struct {
	n int
	a [10]int
}

func genCaseB(rng *rand.Rand) caseB {
	n := rng.Intn(20) + 1
	var a [10]int
	for i := 0; i < 10; i++ {
		a[i] = rng.Intn(4) // small to keep compute quick
	}
	return caseB{n: n, a: a}
}

func solveB(tc caseB) int64 {
	n := tc.n
	a := tc.a
	sumA := 0
	maxN := n
	for _, v := range a {
		sumA += v
		if v > maxN {
			maxN = v
		}
	}
	fact := make([]int64, maxN+1)
	invfact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invfact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invfact[i-1] = invfact[i] * int64(i) % mod
	}
	pow10 := make([]int64, n+1)
	pow10[0] = 1
	for i := 1; i <= n; i++ {
		pow10[i] = pow10[i-1] * 10 % mod
	}
	invAFact := int64(1)
	for i := 0; i < 10; i++ {
		if a[i] > maxN {
			return 0
		}
		invAFact = invAFact * invfact[a[i]] % mod
	}
	a0p := a[0]
	if a0p > 0 {
		a0p--
	}
	sumAp := sumA
	if a[0] > 0 {
		sumAp--
	}
	invAPrime := invfact[a0p]
	for i := 1; i < 10; i++ {
		invAPrime = invAPrime * invfact[a[i]] % mod
	}
	var ans int64
	for L := 1; L <= n; L++ {
		var tot int64
		if L >= sumA {
			tot = fact[L] * pow10[L-sumA] % mod * invfact[L-sumA] % mod * invAFact % mod
		}
		var inv int64
		if L-1 >= sumAp {
			inv = fact[L-1] * pow10[L-1-sumAp] % mod * invfact[L-1-sumAp] % mod * invAPrime % mod
		}
		ans = (ans + tot - inv + mod) % mod
	}
	return ans
}

func runB(bin string, tc caseB) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i := 0; i < 10; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, tc.a[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveB(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCaseB(rng)
		if err := runB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
