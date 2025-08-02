package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	N      = 200007
	MaxBit = 400014
	Words  = (MaxBit + 63) / 64
)

type Bitset [Words]uint64

func (bs *Bitset) SetAll() {
	for i := range bs {
		bs[i] = ^uint64(0)
	}
}

func (bs *Bitset) ResetBit(i int) {
	bs[i>>6] &= ^(1 << (i & 63))
}

func (bs *Bitset) Set(i int) {
	bs[i>>6] |= 1 << (i & 63)
}

func (bs *Bitset) Test(i int) bool {
	if i >= MaxBit {
		return false
	}
	return bs[i>>6]&(1<<(i&63)) != 0
}

func (bs *Bitset) Or(other *Bitset) {
	for i := range bs {
		bs[i] |= other[i]
	}
}

func And(dest, a, b *Bitset) {
	for i := range dest {
		dest[i] = a[i] & b[i]
	}
}

func ShiftLeft(dest *Bitset, src *Bitset, sh int) {
	for i := range dest {
		dest[i] = 0
	}
	if sh == 0 {
		*dest = *src
		return
	}
	q := sh >> 6
	r := sh & 63
	for i := 0; i < Words-q; i++ {
		dest[i+q] |= src[i] << r
		if r != 0 && i+q+1 < Words {
			dest[i+q+1] |= src[i] >> (64 - r)
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var p, f, g, tmp, shiftTmp Bitset
	p.SetAll()
	p.ResetBit(0)
	p.ResetBit(1)
	for i := 2; i < N; i++ {
		if p.Test(i) {
			for j := i * 2; j < N; j += i {
				p.ResetBit(j)
			}
		}
	}
	for i := 3; i < N; i += 2 {
		if p.Test(i-2) && !f.Test(i-2) {
			f.Set(i)
		}
		if p.Test(i) && !f.Test(i) {
			tmp.Set(i)
		}
	}
	for i := 3; i < N; i += 2 {
		if tmp.Test(i) {
			ShiftLeft(&shiftTmp, &tmp, i)
			f.Or(&shiftTmp)
		}
	}
	f.Set(4)
	And(&tmp, &f, &p)
	for i := 3; i < N; i += 2 {
		if tmp.Test(i) {
			ShiftLeft(&shiftTmp, &tmp, i)
			g.Or(&shiftTmp)
		}
	}
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		s0, s1 := 0, 0
		for i := 1; i <= n; i++ {
			if f.Test(a[i]) {
				s0++
			}
			if g.Test(a[i]) {
				s1++
			}
		}
		if s0 == 0 || (s0 == n && n%2 == 1 && (s1 == 0 || s1 == n)) {
			fmt.Println("Bob")
		} else {
			fmt.Println("Alice")
		}
	}
}
