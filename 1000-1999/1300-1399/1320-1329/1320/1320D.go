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

type hasher struct {
	pre  []int64
	pow  []int64
	ones []int64
	mod  int64
	base int64
}

func newHasher(bits []int, base, mod int64) *hasher {
	n := len(bits)
	h := &hasher{
		pre:  make([]int64, n+1),
		pow:  make([]int64, n+1),
		ones: make([]int64, n+1),
		mod:  mod,
		base: base,
	}
	h.pow[0] = 1
	for i := 0; i < n; i++ {
		h.pow[i+1] = (h.pow[i] * base) % mod
		h.ones[i+1] = (h.ones[i]*base + 1) % mod
		h.pre[i+1] = (h.pre[i]*base + int64(bits[i])) % mod
	}
	return h
}

func (h *hasher) getRange(l, r int) int64 {
	if l > r {
		return 0
	}
	res := (h.pre[r+1] - (h.pre[l]*h.pow[r-l+1])%h.mod) % h.mod
	if res < 0 {
		res += h.mod
	}
	return res
}

func (h *hasher) normalize(l, r int, parity int) int64 {
	if l > r {
		return 0
	}
	hash := h.getRange(l, r)
	if parity&1 == 1 {
		hash = (h.ones[r-l+1] - hash) % h.mod
		if hash < 0 {
			hash += h.mod
		}
	}
	return hash
}

func leadingParity(l, r int, nextZero []int) int {
	firstZero := nextZero[l]
	if firstZero > r {
		return (r - l + 1) & 1
	}
	return (firstZero - l) & 1
}

func trailingParity(l, r int, prevZero []int) int {
	lastZero := prevZero[r]
	if lastZero < l {
		return (r - l + 1) & 1
	}
	return (r - lastZero) & 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	prefOnes := make([]int, n+1)
	prefZeros := make([]int, n+1)
	prevZero := make([]int, n+1)
	nextZero := make([]int, n+2)
	zeroBits := make([]int, 0)

	for i := 1; i <= n; i++ {
		prefOnes[i] = prefOnes[i-1]
		prefZeros[i] = prefZeros[i-1]
		if s[i-1] == '1' {
			prefOnes[i]++
			prevZero[i] = prevZero[i-1]
		} else {
			prefZeros[i]++
			prevZero[i] = i
			zeroBits = append(zeroBits, prefOnes[i-1]&1)
		}
	}
	nextZero[n+1] = n + 1
	for i := n; i >= 1; i-- {
		if s[i-1] == '0' {
			nextZero[i] = i
		} else {
			nextZero[i] = nextZero[i+1]
		}
	}

	hasher1 := newHasher(zeroBits, base1, mod1)
	hasher2 := newHasher(zeroBits, base2, mod2)

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var l1, l2, length int
		fmt.Fscan(in, &l1, &l2, &length)
		r1 := l1 + length - 1
		r2 := l2 + length - 1

		ones1 := prefOnes[r1] - prefOnes[l1-1]
		ones2 := prefOnes[r2] - prefOnes[l2-1]
		if ones1 != ones2 {
			fmt.Fprintln(out, "NO")
			continue
		}
		if ones1 == length {
			fmt.Fprintln(out, "YES")
			continue
		}

		lead1 := leadingParity(l1, r1, nextZero)
		lead2 := leadingParity(l2, r2, nextZero)
		if lead1 != lead2 {
			fmt.Fprintln(out, "NO")
			continue
		}
		trail1 := trailingParity(l1, r1, prevZero)
		trail2 := trailingParity(l2, r2, prevZero)
		if trail1 != trail2 {
			fmt.Fprintln(out, "NO")
			continue
		}

		baseParity1 := prefOnes[l1-1] & 1
		baseParity2 := prefOnes[l2-1] & 1

		zeroStart1 := prefZeros[l1-1]
		zeroEnd1 := prefZeros[r1] - 1
		zeroStart2 := prefZeros[l2-1]
		zeroEnd2 := prefZeros[r2] - 1

		hash1a := hasher1.normalize(zeroStart1, zeroEnd1, baseParity1)
		hash1b := hasher2.normalize(zeroStart1, zeroEnd1, baseParity1)
		hash2a := hasher1.normalize(zeroStart2, zeroEnd2, baseParity2)
		hash2b := hasher2.normalize(zeroStart2, zeroEnd2, baseParity2)

		if hash1a == hash2a && hash1b == hash2b {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
