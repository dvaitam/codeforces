package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const mod int64 = 998244353
const root int64 = 3

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 {
	return modPow(a, mod-2)
}

func ntt(a []int64, invert bool) {
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
		wlen := modPow(root, (mod-1)/int64(length))
		if invert {
			wlen = modInv(wlen)
		}
		for i := 0; i < n; i += length {
			w := int64(1)
			half := length / 2
			for j := 0; j < half; j++ {
				u := a[i+j]
				v := a[i+j+half] * w % mod
				a[i+j] = (u + v) % mod
				a[i+j+half] = (u - v + mod) % mod
				w = w * wlen % mod
			}
		}
	}
	if invert {
		invN := modInv(int64(n))
		for i := 0; i < n; i++ {
			a[i] = a[i] * invN % mod
		}
	}
}

func convolution(a, b []int64, limit int) []int64 {
	n := 1
	for n < len(a)+len(b)-1 {
		n <<= 1
	}
	fa := make([]int64, n)
	fb := make([]int64, n)
	copy(fa, a)
	copy(fb, b)
	ntt(fa, false)
	ntt(fb, false)
	for i := 0; i < n; i++ {
		fa[i] = fa[i] * fb[i] % mod
	}
	ntt(fa, true)
	resLen := len(a) + len(b) - 1
	if resLen > limit+1 {
		resLen = limit + 1
	}
	return fa[:resLen]
}

type Poly struct{ coeff []int64 }

type PQ []Poly

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return len(pq[i].coeff) < len(pq[j].coeff) }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Poly)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		var p int
		fmt.Fscan(in, &p)
		counts[p]++
	}
	limit := n / 2
	pq := &PQ{}
	heap.Init(pq)
	for _, c := range counts {
		sz := c
		if sz > limit {
			sz = limit
		}
		poly := make([]int64, sz+1)
		for i := 0; i <= sz; i++ {
			poly[i] = 1
		}
		heap.Push(pq, Poly{poly})
	}
	for pq.Len() > 1 {
		a := heap.Pop(pq).(Poly)
		b := heap.Pop(pq).(Poly)
		c := convolution(a.coeff, b.coeff, limit)
		heap.Push(pq, Poly{c})
	}
	res := heap.Pop(pq).(Poly).coeff
	if limit < len(res) {
		fmt.Println(res[limit] % mod)
	} else {
		fmt.Println(0)
	}
}
