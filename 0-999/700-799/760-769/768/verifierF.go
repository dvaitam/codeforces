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

const MOD int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

var fac, ifac []int64

func comb(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expectedProbability(f, w, h int) int64 {
	if w == 0 {
		return 1
	}
	if f == 0 {
		if w > h {
			return 1
		}
		return 0
	}
	limit := f + w + 5
	fac = make([]int64, limit)
	ifac = make([]int64, limit)
	fac[0] = 1
	for i := 1; i < limit; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}
	ifac[limit-1] = modInv(fac[limit-1])
	for i := limit - 1; i > 0; i-- {
		ifac[i-1] = ifac[i] * int64(i) % MOD
	}
	var total, liked int64
	for x := 1; x <= f; x++ {
		for _, dy := range []int{-1, 0, 1} {
			y := x + dy
			if y < 1 || y > w {
				continue
			}
			if abs(x-y) > 1 {
				continue
			}
			patterns := int64(1)
			if x == y {
				patterns = 2
			}
			waysF := comb(f-1, x-1)
			waysW := comb(w-1, y-1)
			total = (total + patterns*waysF%MOD*waysW) % MOD
			if w >= (h+1)*y {
				waysWL := comb(w-(h+1)*y+y-1, y-1)
				liked = (liked + patterns*waysF%MOD*waysWL) % MOD
			}
		}
	}
	invTotal := modInv(total)
	ans := liked % MOD * invTotal % MOD
	return ans
}

type testCase struct {
	f      int
	w      int
	h      int
	expect int64
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	deterministic := []testCase{
		{f: 1, w: 1, h: 1},
		{f: 1, w: 2, h: 1},
		{f: 2, w: 0, h: 5},
	}
	for _, tc := range deterministic {
		tc.expect = expectedProbability(tc.f, tc.w, tc.h)
		tests = append(tests, tc)
	}
	for len(tests) < 100 {
		f := rng.Intn(5)
		w := rng.Intn(5)
		if f == 0 && w == 0 {
			f = 1
		}
		h := rng.Intn(5)
		tc := testCase{f: f, w: w, h: h}
		tc.expect = expectedProbability(f, w, h)
		tests = append(tests, tc)
	}
	return tests
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d %d\n", tc.f, tc.w, tc.h)
		expected := strconv.FormatInt(tc.expect, 10)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
