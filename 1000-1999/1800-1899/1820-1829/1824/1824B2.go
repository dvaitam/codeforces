package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1_000_000_007

var fact, invFact []int64

func modPow(a, e int64) int64 {
	r := int64(1)
	for e > 0 {
		if e&1 == 1 {
			r = r * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return r
}

func initComb(n int) {
	fact = make([]int64, n+1)
	invFact = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[n] = modPow(fact[n], mod-2)
	for i := n; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	initComb(n)

	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	// compute subtree sizes with root 0
	parent := make([]int, n)
	order := make([]int, 0, n)
	stack := []int{0}
	parent[0] = -1
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			stack = append(stack, to)
		}
	}

	size := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		v := order[i]
		size[v] = 1
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			size[v] += size[to]
		}
	}

	half := k / 2
	combNK := comb(int64(n), int64(k))
	invCombNK := modPow(combNK, mod-2)

	cache := make(map[int]int64)

	calc := func(s int) int64 {
		if s <= half {
			return 0
		}
		if v, ok := cache[s]; ok {
			return v
		}
		var sum int64
		upper := k
		if s < upper {
			upper = s
		}
		for j := half + 1; j <= upper; j++ {
			sum = (sum + comb(int64(s), int64(j))*comb(int64(n-s), int64(k-j))) % mod
		}
		cache[s] = sum
		return sum
	}

	var total int64
	for v := 1; v < n; v++ {
		s := size[v]
		total = (total + calc(s)) % mod
		total = (total + calc(n-s)) % mod
	}

	ans := (int64(n)%mod*combNK%mod - total + mod) % mod
	ans = ans * invCombNK % mod
	fmt.Fprintln(out, ans)
}
