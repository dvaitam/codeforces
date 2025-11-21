package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 998244353
const maxN = 3000000

var (
	fac    []int
	invFac []int
	spf    []int
)

func modPow(a, e int) int {
	res := 1
	base := a
	exp := e
	for exp > 0 {
		if exp&1 == 1 {
			res = int(int64(res) * int64(base) % mod)
		}
		base = int(int64(base) * int64(base) % mod)
		exp >>= 1
	}
	return res
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func precompute() {
	fac = make([]int, maxN+1)
	invFac = make([]int, maxN+1)
	fac[0] = 1
	for i := 1; i <= maxN; i++ {
		fac[i] = int(int64(fac[i-1]) * int64(i) % mod)
	}
	invFac[maxN] = modPow(fac[maxN], mod-2)
	for i := maxN; i >= 1; i-- {
		invFac[i-1] = int(int64(invFac[i]) * int64(i) % mod)
	}
	spf = make([]int, maxN+1)
	for i := 2; i <= maxN; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= maxN/i {
				for j := i * i; j <= maxN; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	spf[0] = 0
	spf[1] = 1
}

func getDivisors(x int) []int {
	divs := []int{1}
	for x > 1 {
		p := spf[x]
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		size := len(divs)
		mul := 1
		for c := 1; c <= cnt; c++ {
			mul *= p
			for i := 0; i < size; i++ {
				divs = append(divs, divs[i]*mul)
			}
		}
	}
	return divs
}

func main() {
	precompute()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, c, k int
		fmt.Fscan(in, &a, &b, &c, &k)
		total := a * b * c
		vals := make([]int, 0)
		cntVals := make([]int, 0)
		g := 0
		prev := -1
		for i := 0; i < k; i++ {
			var d int
			fmt.Fscan(in, &d)
			if prev == -1 || d != prev {
				vals = append(vals, d)
				cntVals = append(cntVals, 1)
				prev = d
			} else {
				cntVals[len(cntVals)-1]++
			}
			if i == 0 {
				g = d
			} else {
				g = gcd(g, d)
			}
		}

	divs := getDivisors(g)
	sort.Ints(divs)
	m := len(divs)
	idxMap := make(map[int]int, m)
	for i, d := range divs {
		idxMap[d] = i
	}

	prod := make([]int, m)
	for i := range prod {
		prod[i] = 1
	}

	for id, x := range vals {
		freq := cntVals[id]
		divsX := getDivisors(x)
		for _, d := range divsX {
			if idx, ok := idxMap[d]; ok {
				y := x / d
				term := modPow(invFac[y], freq)
				prod[idx] = int(int64(prod[idx]) * int64(term) % mod)
			}
		}
	}

	count := make([]int64, m)
	for i := 0; i < m; i++ {
		d := divs[i]
		val := int64(gcd(a, d)) * int64(gcd(b, d)) * int64(gcd(c, d))
		count[i] = val
		for j := 0; j < i; j++ {
			if d%divs[j] == 0 {
				count[i] -= count[j]
			}
		}
	}

	invTotal := modPow(total%mod, mod-2)
	ans := 0
	for i := 0; i < m; i++ {
		L := divs[i]
		cycles := total / L
		term := int(int64(fac[cycles]) * int64(prod[i]) % mod)
		add := int((count[i] % int64(mod) + int64(mod)) % int64(mod))
		ans = (ans + int(int64(add)*int64(term)%int64(mod))) % mod
	}

	ans = int(int64(ans) * int64(invTotal) % mod)
	fmt.Fprintln(out, ans)
	}
}

