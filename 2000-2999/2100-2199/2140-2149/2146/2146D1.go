package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// binaryTrie stores unused numbers from [0, r] and supports extracting the one
// that maximizes the current position's bitwise OR value.
type binaryTrie struct {
	next   [][2]int
	cnt    []int
	maxBit int
	ptr    int
	path   []int
}

func newBinaryTrie(maxBit int) *binaryTrie {
	if maxBit < 0 {
		maxBit = 0
	}
	maxNodes := 1 << (maxBit + 2)
	t := &binaryTrie{
		next:   make([][2]int, maxNodes),
		cnt:    make([]int, maxNodes),
		maxBit: maxBit,
		ptr:    1,
		path:   make([]int, maxBit+2),
	}
	for i := range t.next {
		t.next[i][0] = -1
		t.next[i][1] = -1
	}
	return t
}

func (t *binaryTrie) add(x int) {
	node := 0
	t.cnt[node]++
	for b := t.maxBit; b >= 0; b-- {
		bit := (x >> b) & 1
		nxt := t.next[node][bit]
		if nxt == -1 {
			nxt = t.ptr
			t.ptr++
			t.next[node][bit] = nxt
		}
		node = nxt
		t.cnt[node]++
	}
}

func (t *binaryTrie) popBest(val int) int {
	node := 0
	idx := 0
	res := 0
	for b := t.maxBit; b >= 0; b-- {
		t.path[idx] = node
		idx++
		bit := (val >> b) & 1
		first := bit ^ 1
		choice := -1
		if nxt := t.next[node][first]; nxt != -1 && t.cnt[nxt] > 0 {
			choice = first
		} else if nxt := t.next[node][bit]; nxt != -1 && t.cnt[nxt] > 0 {
			choice = bit
		}
		if choice == -1 {
			panic("no available number for current prefix")
		}
		res |= choice << b
		node = t.next[node][choice]
	}
	t.path[idx] = node
	idx++
	for i := 0; i < idx; i++ {
		t.cnt[t.path[i]]--
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int
		fmt.Fscan(in, &l, &r) // easy version guarantees l = 0
		n := r + 1
		maxBit := 0
		if r > 0 {
			maxBit = bits.Len(uint(r)) - 1
		}
		trie := newBinaryTrie(maxBit)
		for x := 0; x <= r; x++ {
			trie.add(x)
		}
		perm := make([]int, n)
		for i := r; i >= 0; i-- {
			perm[i] = trie.popBest(i)
		}
		fmt.Fprintln(out, int64(n)*int64(r))
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, perm[i])
		}
		fmt.Fprintln(out)
	}
}
