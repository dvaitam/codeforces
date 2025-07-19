package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// Bitset for fixed size bit operations
type Bitset struct {
	bits []uint64
	size int
}

// NewBitset creates a bitset of given size
func NewBitset(size int) Bitset {
	return Bitset{bits: make([]uint64, (size+63)/64), size: size}
}

// FlipRange flips bits in [l, r] (1-based inclusive)
func (b *Bitset) FlipRange(l, r int) {
	for i := l - 1; i <= r-1; i++ {
		w := i / 64
		o := uint(i % 64)
		b.bits[w] ^= 1 << o
	}
}

// Count returns number of set bits
func (b *Bitset) Count() int {
	cnt := 0
	for _, w := range b.bits {
		cnt += bits.OnesCount64(w)
	}
	return cnt
}

// And returns a new Bitset as bitwise AND
func (b *Bitset) And(o *Bitset) Bitset {
	res := NewBitset(b.size)
	for i := range b.bits {
		res.bits[i] = b.bits[i] & o.bits[i]
	}
	return res
}

// Xor performs in-place bitwise XOR
func (b *Bitset) Xor(o *Bitset) {
	for i := range b.bits {
		b.bits[i] ^= o.bits[i]
	}
}

// Clone returns a copy of the bitset
func (b *Bitset) Clone() Bitset {
	res := NewBitset(b.size)
	copy(res.bits, b.bits)
	return res
}

// FindFirst returns zero-based index of first set bit, or -1
func (b *Bitset) FindFirst() int {
	for wi, w := range b.bits {
		if w != 0 {
			t := bits.TrailingZeros64(w)
			idx := wi*64 + t
			if idx < b.size {
				return idx
			}
		}
	}
	return -1
}

// pair holds size and id
type pair struct{ sz, id int }

// qi holds queued pairs for checking
type qi struct{ x, y int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m, Q int
	fmt.Fscan(in, &n, &m, &Q)
	// initialize data
	foo := make([]Bitset, n+1)
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		foo[i] = NewBitset(m)
		sz[i] = 0
	}
	// set and queue
	S := make([]pair, 0, n+3)
	// sentinel values
	S = append(S, pair{sz: int(1e9), id: -1})
	S = append(S, pair{sz: -1, id: -1})
	// queue
	q := make([]qi, 0, Q*2+10)
	head := 0
	// helpers for queue
	push := func(a, b int) {
		q = append(q, qi{a, b})
	}
	// chk function
	chk := func(x, y int) bool {
		if x < 1 || y < 1 {
			return false
		}
		if sz[x] > sz[y] {
			x, y = y, x
		}
		// check if foo[x] is not subset of foo[y]
		and := foo[x].And(&foo[y])
		// subset if and == foo[x]
		// compare bits
		for i := range and.bits {
			if and.bits[i] != foo[x].bits[i] {
				return true
			}
		}
		return false
	}
	// add element x to S
	add := func(x int) {
		// find lower_bound
		idx := 0
		for idx < len(S) && (S[idx].sz < sz[x] || (S[idx].sz == sz[x] && S[idx].id < x)) {
			idx++
		}
		// neighbors
		pre := S[idx-1]
		cur := S[idx]
		// insert
		S = append(S[:idx], append([]pair{{sz[x], x}}, S[idx:]...)...)
		push(x, pre.id)
		push(cur.id, x)
	}
	// delete element x from S
	del := func(x int) {
		// find element
		idx := 0
		for idx < len(S) && !(S[idx].sz == sz[x] && S[idx].id == x) {
			idx++
		}
		pre := S[idx-1]
		nxt := S[idx+1]
		// remove
		S = append(S[:idx], S[idx+1:]...)
		push(pre.id, nxt.id)
	}
	// initialize S with 1..n
	for i := 1; i <= n; i++ {
		add(i)
	}
	// process queries
	for qi := 0; qi < Q; qi++ {
		var x, l, r int
		fmt.Fscan(in, &x, &l, &r)
		del(x)
		foo[x].FlipRange(l, r)
		sz[x] = foo[x].Count()
		add(x)
		// find valid in queue
		for head < len(q) && !chk(q[head].x, q[head].y) {
			head++
		}
		if head >= len(q) {
			fmt.Fprintln(out, -1)
		} else {
			a := q[head].x
			b := q[head].y
			if a > b {
				a, b = b, a
			}
			inter := foo[a].And(&foo[b])
			t1 := foo[a].Clone()
			t1.Xor(&inter)
			id1 := t1.FindFirst() + 1
			t2 := foo[b].Clone()
			t2.Xor(&inter)
			id2 := t2.FindFirst() + 1
			if id1 > id2 {
				id1, id2 = id2, id1
			}
			fmt.Fprintf(out, "%d %d %d %d\n", a, id1, b, id2)
		}
	}
}
