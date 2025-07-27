package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353

var fact, inv []int64

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}
func initComb(n int) {
	fact = make([]int64, n+1)
	inv = make([]int64, n+1)
	fact[0] = 1
	for i := 1; i <= n; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv[n] = modPow(fact[n], MOD-2)
	for i := n; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
}
func C(n, k int) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * inv[k] % MOD * inv[n-k] % MOD
}

type BIT struct {
	n    int
	tree []int
}

func (b *BIT) init(n int) { b.n = n; b.tree = make([]int, n+2) }
func (b *BIT) add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}
func (b *BIT) kth(k int) int {
	idx := 0
	bit := 1
	for bit<<1 <= b.n {
		bit <<= 1
	}
	for bit > 0 {
		next := idx + bit
		if next <= b.n && b.tree[next] < k {
			k -= b.tree[next]
			idx = next
		}
		bit >>= 1
	}
	return idx + 1
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	initComb(400000)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		j := make([]int, n+1)
		for i := 1; i <= n; i++ {
			j[i] = i
		}
		for k := 0; k < m; k++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			j[x] = y
		}
		bit := BIT{}
		bit.init(n)
		for i := 1; i <= n; i++ {
			bit.add(i, 1)
		}
		p := make([]int, n+1)
		for i := n; i >= 1; i-- {
			pos := bit.kth(j[i])
			p[pos] = i
			bit.add(pos, -1)
		}
		r := 0
		for i := 1; i < n; i++ {
			if p[i] > p[i+1] {
				r++
			}
		}
		ans := C(2*n-1-r, n)
		fmt.Fprintln(out, ans)
	}
}
