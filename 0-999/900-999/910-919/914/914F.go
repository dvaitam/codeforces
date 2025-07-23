package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const threshold = 320

var (
	s       []byte
	n       int
	wordLen int
	bitsArr [26][]uint64
	tmp     []uint64
)

func setBit(arr []uint64, idx int, val bool) {
	w := idx / 64
	b := uint(idx % 64)
	if val {
		arr[w] |= 1 << b
	} else {
		arr[w] &^= 1 << b
	}
}

func andShiftInto(dst []uint64, src []uint64, shift int) {
	ws := shift / 64
	bs := uint(shift % 64)
	n := len(dst)
	for i := 0; i < n; i++ {
		idx := i + ws
		var v uint64
		if idx < len(src) {
			v = src[idx] >> bs
		}
		if bs != 0 && idx+1 < len(src) {
			v |= src[idx+1] << (64 - bs)
		}
		dst[i] &= v
	}
}

func countRange(arr []uint64, l, r int) int {
	if r < l {
		return 0
	}
	lw := l / 64
	rw := r / 64
	lb := uint(l % 64)
	rb := uint(r % 64)
	if lw == rw {
		mask := ((uint64(1) << (rb - lb + 1)) - 1) << lb
		return bits.OnesCount64(arr[lw] & mask)
	}
	cnt := bits.OnesCount64(arr[lw] & (^uint64(0) << lb))
	for i := lw + 1; i < rw; i++ {
		cnt += bits.OnesCount64(arr[i])
	}
	mask := (uint64(1) << (rb + 1)) - 1
	cnt += bits.OnesCount64(arr[rw] & mask)
	return cnt
}

func querySmall(l, r int, pat []byte) int {
	m := len(pat)
	if m == 0 {
		return 0
	}
	copy(tmp, bitsArr[pat[0]-'a'])
	for j := 1; j < m; j++ {
		andShiftInto(tmp, bitsArr[pat[j]-'a'], j)
	}
	start := l
	end := r - m + 1
	if end < start {
		return 0
	}
	return countRange(tmp, start, end)
}

func queryLarge(l, r int, pat []byte) int {
	m := len(pat)
	end := r - m + 1
	if end < l {
		return 0
	}
	cnt := 0
	for i := l; i <= end; i++ {
		match := true
		for j := 0; j < m; j++ {
			if s[i+j] != pat[j] {
				match = false
				break
			}
		}
		if match {
			cnt++
		}
	}
	return cnt
}

func processQuery(l, r int, pat []byte) int {
	m := len(pat)
	if m > r-l+1 {
		return 0
	}
	if m > threshold {
		return queryLarge(l, r, pat)
	}
	return querySmall(l, r, pat)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var str string
	fmt.Fscan(in, &str)
	s = []byte(str)
	n = len(s)
	wordLen = (n + 63) / 64
	for i := 0; i < 26; i++ {
		bitsArr[i] = make([]uint64, wordLen)
	}
	tmp = make([]uint64, wordLen)
	for i, ch := range s {
		setBit(bitsArr[ch-'a'], i, true)
	}
	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			var pos int
			var c string
			fmt.Fscan(in, &pos, &c)
			pos--
			old := s[pos]
			newc := c[0]
			if old != newc {
				setBit(bitsArr[old-'a'], pos, false)
				setBit(bitsArr[newc-'a'], pos, true)
				s[pos] = newc
			}
		} else if t == 2 {
			var l, r int
			var y string
			fmt.Fscan(in, &l, &r, &y)
			pat := []byte(y)
			ans := processQuery(l-1, r-1, pat)
			fmt.Fprintln(out, ans)
		}
	}
}
