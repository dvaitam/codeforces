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

const (
	mod  = 1000000007
	maxK = 1000000 + 5
)

var fac, ifac []int64

func modPow(base, exp int64) int64 {
	base %= mod
	if base < 0 {
		base += mod
	}
	res := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func initFactorials(limit int) {
	fac = make([]int64, limit+1)
	ifac = make([]int64, limit+1)
	fac[0] = 1
	for i := 1; i <= limit; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}
	ifac[limit] = modPow(fac[limit], mod-2)
	for i := limit; i >= 1; i-- {
		ifac[i-1] = ifac[i] * int64(i) % mod
	}
}

func fall(n, r int) int64 {
	if r < 0 || r > n {
		return 0
	}
	return fac[n] * ifac[n-r] % mod
}

func getPow(powVals []int64, idx int, k int, w int64) int64 {
	if idx <= 0 {
		return 1
	}
	if idx < len(powVals) && powVals[idx] != 0 {
		return powVals[idx]
	}
	exp := w - int64(idx)
	if exp < 0 {
		return 0
	}
	return modPow(int64(k), exp)
}

func countLen(k int, w int64, t int, powVals []int64) int64 {
	if t > k {
		return 0
	}
	if int64(t) <= w {
		ff := fall(k, t)
		powVal := getPow(powVals, t, k, w)
		return ff * ff % mod * powVal % mod
	}
	d := int(int64(t) - w)
	if d <= 0 || d > k {
		return 0
	}
	restSize := k - d
	if int64(restSize) < w {
		return 0
	}
	rest := fall(restSize, int(w))
	overlap := fall(k, d)
	return overlap * rest % mod * rest % mod
}

func countLenNext(k int, w int64, t int, powVals []int64) int64 {
	if t+1 > k && int64(t) <= w-2 {
		return 0
	}
	if int64(t) <= w-2 {
		ff := fall(k, t+1)
		powVal := getPow(powVals, t+2, k, w)
		return ff * ff % mod * powVal % mod
	}
	d := int(int64(t) - w + 2)
	if d <= 0 || d > k {
		return 0
	}
	restSize := k - d
	target := w - 1
	if target < 0 || int64(restSize) < target {
		return 0
	}
	rest := fall(restSize, int(target))
	overlap := fall(k, d)
	return overlap * rest % mod * rest % mod
}

func expected(k int, w int64) int64 {
	if k == 0 {
		return 0
	}
	powLimit := k + 2
	if int64(powLimit) > w {
		powLimit = int(w)
	}
	powVals := make([]int64, powLimit+2)
	if powLimit >= 1 {
		powVals[1] = modPow(int64(k), w-1)
		inv := modPow(int64(k), mod-2)
		for i := 2; i <= powLimit; i++ {
			powVals[i] = powVals[i-1] * inv % mod
		}
	}
	ans := int64(0)
	for t := 1; t <= k; t++ {
		ans += countLen(k, w, t, powVals)
		ans -= countLenNext(k, w, t, powVals)
	}
	ans %= mod
	if ans < 0 {
		ans += mod
	}
	return ans
}

type testCase struct {
	input  string
	expect string
}

func buildCase(k int, w int64) testCase {
	exp := expected(k, w)
	return testCase{
		input:  fmt.Sprintf("%d %d\n", k, w),
		expect: fmt.Sprint(exp),
	}
}

func deterministicCases() []testCase {
	spec := []struct {
		k int
		w int64
	}{
		{1, 2},
		{1, 10},
		{2, 2},
		{2, 5},
		{3, 5},
		{4, 8},
		{10, 25},
		{100, 200},
		{1000000, 2},
		{1000000, 1000},
		{999999, 1000000000},
	}
	cases := make([]testCase, 0, len(spec))
	for _, tc := range spec {
		cases = append(cases, buildCase(tc.k, tc.w))
	}
	return cases
}

func randomCase(rng *rand.Rand) testCase {
	switch rng.Intn(3) {
	case 0:
		k := rng.Intn(4) + 1
		w := int64(rng.Intn(6) + 2)
		return buildCase(k, w)
	case 1:
		k := rng.Intn(200) + 1
		w := int64(rng.Intn(400) + 2)
		return buildCase(k, w)
	default:
		k := rng.Intn(1000000) + 1
		w := rng.Int63n(1000000000-1) + 2
		return buildCase(k, w)
	}
}

func runCandidate(bin, input string) (string, error) {
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
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	initFactorials(maxK)
	bin := args[0]
	tests := deterministicCases()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		tests = append(tests, randomCase(rng))
	}
	for idx, tc := range tests {
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expect, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
