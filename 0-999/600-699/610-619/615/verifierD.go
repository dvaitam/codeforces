package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const mod int64 = 1000000007

func modPow(base, exp, m int64) int64 {
	res := int64(1)
	base %= m
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % m
		}
		base = base * base % m
		exp >>= 1
	}
	return res
}

func expected(primes []int) int64 {
	counts := make(map[int]int)
	for _, p := range primes {
		counts[p]++
	}
	type pair struct{ p, c int }
	arr := make([]pair, 0, len(counts))
	for p, c := range counts {
		arr = append(arr, pair{p, c})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].p < arr[j].p })

	k := len(arr)
	modMinus1 := mod - 1
	pref := make([]int64, k+1)
	suff := make([]int64, k+1)
	pref[0] = 1
	for i := 0; i < k; i++ {
		pref[i+1] = pref[i] * int64(arr[i].c+1) % modMinus1
	}
	suff[k] = 1
	for i := k - 1; i >= 0; i-- {
		suff[i] = suff[i+1] * int64(arr[i].c+1) % modMinus1
	}

	ans := int64(1)
	for i := 0; i < k; i++ {
		a := int64(arr[i].c)
		expPart := a * (a + 1) / 2 % modMinus1
		other := pref[i] * suff[i+1] % modMinus1
		exp := expPart * other % modMinus1
		ans = ans * modPow(int64(arr[i].p), exp, mod) % mod
	}
	return ans
}

var primesList = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}

func genCase(rng *rand.Rand) (string, int64) {
	m := rng.Intn(6) + 1
	arr := make([]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		p := primesList[rng.Intn(len(primesList))]
		arr[i] = p
		sb.WriteString(fmt.Sprintf("%d ", p))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(arr)
}

func runCase(bin, input string, exp int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []string{
		"1\n2\n",
	}
	exps := []int64{2}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}

	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
