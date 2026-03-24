package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 998244353

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

type testG struct {
	a int
	c int
}

func solveReference(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		val := 0
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			val = val*10 + int(data[idx]-'0')
			idx++
		}
		return val
	}

	t := nextInt()
	tests := make([]testG, t)
	maxN := 5
	maxK := 5

	for i := 0; i < t; i++ {
		a := nextInt()
		_ = nextInt()
		c := nextInt()
		tests[i] = testG{a: a, c: c}
		z := a + 5
		if z > maxK {
			maxK = z
		}
		if z+c > maxN {
			maxN = z + c
		}
		if 2*z-6 > maxN {
			maxN = 2*z - 6
		}
	}

	fac := make([]int64, maxN+1)
	ifac := make([]int64, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[maxN] = modPow(fac[maxN], MOD-2)
	for i := maxN; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}

	comb := func(n, r int) int64 {
		if n < 0 || r < 0 || r > n {
			return 0
		}
		return fac[n] * ifac[r] % MOD * ifac[n-r] % MOD
	}

	pref := make([]int64, maxK+1)
	choose5 := make([]int64, maxK+1)
	for k := 5; k <= maxK; k++ {
		v := comb(2*k-6, k-5) - comb(2*k-6, k)
		if v < 0 {
			v += MOD
		}
		pref[k] = v
		choose5[k] = comb(k, 5)
	}

	var out bytes.Buffer
	for i, tc := range tests {
		z := tc.a + 5
		c := tc.c

		invTotal := ifac[z+c] * fac[c] % MOD * fac[z] % MOD
		invZ5 := ifac[z] * fac[5] % MOD * fac[z-5] % MOD

		lim := z - 1
		if c+5 < lim {
			lim = c + 5
		}

		sub := int64(0)
		for k := 5; k <= lim; k++ {
			phit := pref[k] * comb(z+c-2*k+5, c-k+5) % MOD * invTotal % MOD
			loss := 1 - choose5[k]*invZ5%MOD
			if loss < 0 {
				loss += MOD
			}
			sub += phit * loss % MOD
			if sub >= MOD {
				sub -= MOD
			}
		}

		ans := int64(1) - sub
		if ans < 0 {
			ans += MOD
		}

		out.WriteString(strconv.FormatInt(ans, 10))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}

	return out.String()
}

// keep unused import for io
var _ = io.ReadAll

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) string {
	a := rng.Intn(3)
	b := rng.Intn(3)
	c := rng.Intn(3)
	return fmt.Sprintf("1\n%d %d %d\n", a, b, c)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	rng := rand.New(rand.NewSource(49))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected := strings.TrimSpace(solveReference(test))
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
