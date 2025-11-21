package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod1  int64 = 1000000007
	mod2  int64 = 1000000009
	base1 int64 = 911382323
	base2 int64 = 972663749
)

var (
	pow1 = []int64{1}
	pow2 = []int64{1}
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

func ensurePow(n int) {
	for len(pow1) <= n {
		pow1 = append(pow1, pow1[len(pow1)-1]*base1%mod1)
		pow2 = append(pow2, pow2[len(pow2)-1]*base2%mod2)
	}
}

type Hasher struct {
	p1 []int64
	p2 []int64
}

func NewHasher(arr []int64) *Hasher {
	ensurePow(len(arr))
	p1 := make([]int64, len(arr)+1)
	p2 := make([]int64, len(arr)+1)
	for i, v := range arr {
		val1 := v % mod1
		val2 := v % mod2
		p1[i+1] = (p1[i]*base1 + val1) % mod1
		p2[i+1] = (p2[i]*base2 + val2) % mod2
	}
	return &Hasher{p1: p1, p2: p2}
}

func (h *Hasher) hash(l, r int) (int64, int64) {
	lenSeg := r - l
	h1 := (h.p1[r] - h.p1[l]*pow1[lenSeg]) % mod1
	if h1 < 0 {
		h1 += mod1
	}
	h2 := (h.p2[r] - h.p2[l]*pow2[lenSeg]) % mod2
	if h2 < 0 {
		h2 += mod2
	}
	return h1, h2
}

func (h *Hasher) equal(l1, l2, length int) bool {
	if length == 0 {
		return true
	}
	a1, a2 := h.hash(l1, l1+length)
	b1, b2 := h.hash(l2, l2+length)
	return a1 == b1 && a2 == b2
}

func lcp(h *Hasher, i, j, maxLen int) int {
	low, high := 0, maxLen
	for low < high {
		mid := (low + high + 1) >> 1
		if h.equal(i, j, mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return low
}

func better(h *Hasher, arr []int64, i, j, length int) bool {
	if i == j {
		return false
	}
	common := lcp(h, i, j, length)
	if common >= length {
		return false
	}
	return arr[i+common] < arr[j+common]
}

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(in.NextInt())
		}
		b := make([]int64, 2*n)
		copy(b, a)
		for i := 0; i < n; i++ {
			b[n+i] = a[i] + 1
		}
		hasher := NewHasher(b)

		best := 0
		for start := 1; start <= n; start++ {
			if better(hasher, b, start, best, n) {
				best = start
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, b[best+i])
		}
		fmt.Fprintln(out)
	}
}
