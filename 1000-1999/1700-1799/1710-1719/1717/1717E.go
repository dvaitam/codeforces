package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func sievePhi(n int) []int {
	phi := make([]int, n+1)
	for i := 0; i <= n; i++ {
		phi[i] = i
	}
	for i := 2; i <= n; i++ {
		if phi[i] == i {
			for j := i; j <= n; j += i {
				phi[j] = phi[j] / i * (i - 1)
			}
		}
	}
	return phi
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	phi := sievePhi(n)

	var ans int64
	for g := 1; g < n; g++ {
		maxS := (n - 1) / g
		for s := 2; s <= maxS; s++ {
			c := n - g*s
			if c <= 0 {
				continue
			}
			gcdVal := gcd(int64(c), int64(g))
			lcm := (int64(c) / gcdVal % mod) * (int64(g) % mod) % mod
			ans = (ans + (int64(phi[s])%mod)*lcm) % mod
		}
	}
	fmt.Fprintln(out, ans%mod)
}
