package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const MOD int = 998244353
const ROOT int = 3

func modAdd(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}
func modSub(a, b int) int {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}
func modMul(a, b int) int { return int(int64(a) * int64(b) % int64(MOD)) }

func modPow(a, e int) int {
	res := 1
	base := a % MOD
	for e > 0 {
		if e&1 == 1 {
			res = modMul(res, base)
		}
		base = modMul(base, base)
		e >>= 1
	}
	return res
}

func modInv(a int) int { return modPow(a, MOD-2) }

func ntt(a []int, invert bool) {
	n := len(a)
	for i, j := 1, 0; i < n; i++ {
		bit := n >> 1
		for ; j&bit != 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit
		if i < j {
			a[i], a[j] = a[j], a[i]
		}
	}
	for length := 2; length <= n; length <<= 1 {
		wlen := modPow(ROOT, (MOD-1)/length)
		if invert {
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := 1
			half := length >> 1
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := modMul(a[i+j+half], w)
				a[i+j] = modAdd(u, v)
				a[i+j+half] = modSub(u, v)
				w = modMul(w, wlen)
			}
		}
	}
	if invert {
		invN := modInv(n)
		for i := 0; i < n; i++ {
			a[i] = modMul(a[i], invN)
		}
	}
}

func polyMul(a, b []int, limit int) []int {
	need := len(a) + len(b) - 1
	if need > limit {
		need = limit
	}
	size := 1
	for size < need {
		size <<= 1
	}
	fa := make([]int, size)
	fb := make([]int, size)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < size; i++ {
		fa[i] = modMul(fa[i], fb[i])
	}
	ntt(fa, true)
	res := make([]int, need)
	copy(res, fa[:need])
	return res
}

type polyItem struct {
	coef []int
}

type polyHeap []polyItem

func (h polyHeap) Len() int            { return len(h) }
func (h polyHeap) Less(i, j int) bool  { return len(h[i].coef) < len(h[j].coef) }
func (h polyHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *polyHeap) Push(x interface{}) { *h = append(*h, x.(polyItem)) }
func (h *polyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		adj[x] = append(adj[x], y)
		adj[y] = append(adj[y], x)
	}

	parent := make([]int, n+1)
	childCnt := make([]int, n+1)
	q := make([]int, 0, n)
	q = append(q, 1)
	parent[1] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, u := range adj[v] {
			if u == parent[v] {
				continue
			}
			parent[u] = v
			childCnt[v]++
			q = append(q, u)
		}
	}

	freq := make(map[int]int)
	for v := 1; v <= n; v++ {
		if childCnt[v] > 0 {
			freq[childCnt[v]]++
		}
	}

	fac := make([]int, n+1)
	invf := make([]int, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = modMul(fac[i-1], i)
	}
	invf[n] = modInv(fac[n])
	for i := n; i > 0; i-- {
		invf[i-1] = modMul(invf[i], i)
	}

	var pq polyHeap
	for r, cnt := range freq {
		poly := make([]int, cnt+1)
		poly[0] = 1
		pow := 1
		for i := 1; i <= cnt; i++ {
			pow = modMul(pow, r)
			comb := modMul(fac[cnt], modMul(invf[i], invf[cnt-i]))
			poly[i] = modMul(comb, pow)
		}
		heap.Push(&pq, polyItem{poly})
	}
	if pq.Len() == 0 {
		fmt.Fprintln(out, fac[n])
		return
	}
	limit := n + 1
	for pq.Len() > 1 {
		a := heap.Pop(&pq).(polyItem).coef
		b := heap.Pop(&pq).(polyItem).coef
		c := polyMul(a, b, limit)
		heap.Push(&pq, polyItem{c})
	}
	F := heap.Pop(&pq).(polyItem).coef

	ans := 0
	for k, val := range F {
		if k > n {
			break
		}
		term := modMul(val, fac[n-k])
		if k%2 == 1 {
			ans = modSub(ans, term)
		} else {
			ans = modAdd(ans, term)
		}
	}
	fmt.Fprintln(out, ans)
}
