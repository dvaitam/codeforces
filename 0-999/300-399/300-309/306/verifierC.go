package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod int64 = 1000000009
const maxN = 4005

var fac [maxN]int64
var invFac [maxN]int64

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

func initFac() {
	fac[0] = 1
	for i := 1; i < maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	invFac[maxN-1] = modPow(fac[maxN-1], mod-2)
	for i := maxN - 1; i > 0; i-- {
		invFac[i-1] = invFac[i] * int64(i) % mod
	}
}

func comb(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fac[n] * (invFac[k] * invFac[n-k] % mod) % mod
}

func expected(n, w, b int) int64 {
	wf := fac[w]
	bf := fac[b]
	var sum int64
	for t := 2; t <= n-1; t++ {
		y := n - t
		if y < 1 {
			continue
		}
		waysStripes := int64(t - 1)
		c1 := comb(w-1, t-1)
		c2 := comb(b-1, y-1)
		if c1 == 0 || c2 == 0 {
			continue
		}
		sum = (sum + waysStripes*c1%mod*c2) % mod
	}
	return wf * bf % mod * sum % mod
}

type testCase struct{ n, w, b int }

func runCase(bin string, tc testCase) (int64, error) {
	input := fmt.Sprintf("%d %d %d\n", tc.n, tc.w, tc.b)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	s := strings.TrimSpace(out.String())
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer output: %q", s)
	}
	return val % mod, nil
}

func generateCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{3, 2, 1}, {4, 3, 1}, {5, 3, 2}, {10, 7, 3}}
	for len(cases) < 100 {
		n := rng.Intn(20) + 3
		w := rng.Intn(20) + 2
		b := rng.Intn(20) + 1
		if w+b < n {
			continue
		}
		cases = append(cases, testCase{n, w, b})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	initFac()
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := expected(tc.n, tc.w, tc.b)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
