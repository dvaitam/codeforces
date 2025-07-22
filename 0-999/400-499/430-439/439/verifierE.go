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
const MAXN = 100000

var mu [MAXN + 1]int
var fac, invfac [MAXN + 1]int64
var divisors [MAXN + 1][]int

func modexp(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = (res * a) % MOD
		}
		a = (a * a) % MOD
		e >>= 1
	}
	return res
}

func initPre() {
	mu[1] = 1
	primes := []int{}
	isComp := make([]bool, MAXN+1)
	for i := 2; i <= MAXN; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			mu[i] = -1
		}
		for _, p := range primes {
			if i*p > MAXN {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				mu[i*p] = 0
				break
			} else {
				mu[i*p] = -mu[i]
			}
		}
	}
	for i := 1; i <= MAXN; i++ {
		for j := i; j <= MAXN; j += i {
			divisors[j] = append(divisors[j], i)
		}
	}
	fac[0] = 1
	for i := 1; i <= MAXN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	invfac[MAXN] = modexp(fac[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		invfac[i-1] = invfac[i] * int64(i) % MOD
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * invfac[k] % MOD * invfac[n-k] % MOD
}

func solveE(n, f int) int64 {
	var ans int64
	for _, d := range divisors[n] {
		nd := n / d
		if nd < f {
			continue
		}
		c := comb(nd-1, f-1)
		if mu[d] == 0 || c == 0 {
			continue
		}
		ans = (ans + int64(mu[d])*c) % MOD
	}
	if ans < 0 {
		ans += MOD
	}
	return ans
}

type testE struct {
	q     int
	pairs [][2]int
}

func genE(rng *rand.Rand) testE {
	q := rng.Intn(3) + 1
	pairs := make([][2]int, q)
	for i := 0; i < q; i++ {
		n := rng.Intn(50) + 1
		f := rng.Intn(n) + 1
		pairs[i] = [2]int{n, f}
	}
	return testE{q: q, pairs: pairs}
}

func runCase(bin string, tc testE) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.q)
	for _, p := range tc.pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	initPre()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testE{
		{q: 1, pairs: [][2]int{{6, 2}}},
		{q: 1, pairs: [][2]int{{7, 2}}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, genE(rng))
	}
	for i, tc := range cases {
		expectLines := make([]string, tc.q)
		for j, p := range tc.pairs {
			expectLines[j] = fmt.Sprint(solveE(p[0], p[1]))
		}
		expect := strings.Join(expectLines, "\n")
		out, err := runCase(exe, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
