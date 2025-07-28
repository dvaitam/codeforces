package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

const mod = 998244353

func modPow(a, b int) int {
	res := 1
	base := a % mod
	for b > 0 {
		if b&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		b >>= 1
	}
	return res
}

func expectedAnswers(n int, q int, a, barr []int, queries [][2]int) []int {
	prefA := make([]int, n+1)
	prefB := make([]int, n+1)
	totalA, totalB := 0, 0
	for i := 1; i <= n; i++ {
		totalA += a[i]
		totalB += barr[i]
		prefA[i] = prefA[i-1] + a[i]
		prefB[i] = prefB[i-1] + barr[i]
	}
	fac := make([]int, totalB+1)
	invfac := make([]int, totalB+1)
	fac[0] = 1
	for i := 1; i <= totalB; i++ {
		fac[i] = fac[i-1] * i % mod
	}
	invfac[totalB] = modPow(fac[totalB], mod-2)
	for i := totalB; i >= 1; i-- {
		invfac[i-1] = invfac[i] * i % mod
	}
	combPrefix := make([]int, totalB+1)
	prefix := 0
	for i := 0; i <= totalB; i++ {
		c := fac[totalB] * invfac[i] % mod * invfac[totalB-i] % mod
		prefix += c
		if prefix >= mod {
			prefix -= mod
		}
		combPrefix[i] = prefix
	}
	pow2B := modPow(2, totalB)
	invPow2B := modPow(pow2B, mod-2)
	ans := make([]int, q)
	for idx, qv := range queries {
		l, r := qv[0], qv[1]
		ARange := prefA[r] - prefA[l-1]
		BRange := prefB[r] - prefB[l-1]
		nOut := totalB - BRange
		diff := 2*ARange - totalA
		threshold := nOut - diff
		if threshold < 0 {
			ans[idx] = 1
		} else if threshold >= totalB {
			ans[idx] = 0
		} else {
			val := pow2B - combPrefix[threshold]
			if val < 0 {
				val += mod
			}
			ans[idx] = val * invPow2B % mod
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(42)
	n := 10
	q := 100
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rand.Intn(3)
		b[i] = rand.Intn(3)
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, q)
	for i := 1; i <= n; i++ {
		if i > 1 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", a[i])
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", b[i])
	}
	input.WriteByte('\n')
	for i := 0; i < q; i++ {
		fmt.Fprintf(&input, "%d %d\n", queries[i][0], queries[i][1])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	expected := expectedAnswers(n, q, a, b, queries)
	for i := 0; i < q; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing output for query %d\n", i+1)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(scanner.Text(), &got)
		if got != expected[i] {
			fmt.Printf("query %d: expected %d, got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
