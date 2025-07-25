package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const maxBits = 22
const maxNodes = 1 << (maxBits + 1)

// trie structure for fast max AND queries
type trie struct {
	ch   [][2]int32
	mx   []int32
	size int32
}

func newTrie() *trie {
	t := &trie{
		ch:   make([][2]int32, maxNodes),
		mx:   make([]int32, maxNodes),
		size: 1,
	}
	for i := range t.mx {
		t.mx[i] = -1
	}
	return t
}

func (t *trie) reset() {
	t.size = 1
	t.ch[0][0], t.ch[0][1] = 0, 0
	t.mx[0] = -1
}

func (t *trie) add(v int) {
	node := int32(0)
	if int32(v) > t.mx[node] {
		t.mx[node] = int32(v)
	}
	for i := maxBits - 1; i >= 0; i-- {
		b := (v >> i) & 1
		if t.ch[node][b] == 0 {
			t.ch[node][b] = t.size
			t.ch[t.size][0], t.ch[t.size][1] = 0, 0
			t.mx[t.size] = -1
			t.size++
		}
		node = t.ch[node][b]
		if int32(v) > t.mx[node] {
			t.mx[node] = int32(v)
		}
	}
}

func (t *trie) query(mask int) int {
	if t.size == 1 {
		return 0
	}
	node := int32(0)
	res := 0
	for i := maxBits - 1; i >= 0; i-- {
		bit := (mask >> i) & 1
		c0 := t.ch[node][0]
		c1 := t.ch[node][1]
		if bit == 1 {
			if c1 != 0 {
				node = c1
				res |= 1 << i
			} else if c0 != 0 {
				node = c0
			} else {
				break
			}
		} else {
			if c0 == 0 && c1 == 0 {
				break
			} else if c0 == 0 {
				node = c1
			} else if c1 == 0 {
				node = c0
			} else {
				rem := mask & ((1 << i) - 1)
				v0 := int(t.mx[c0]) & rem
				v1 := int(t.mx[c1]) & rem
				if v0 >= v1 {
					node = c0
				} else {
					node = c1
				}
			}
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	trieA := newTrie()
	trieB := newTrie()
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		mask := (1 << bits.Len(uint(n-1))) - 1
		trieA.reset()
		trieB.reset()
		ans := 0
		last := 0
		for i := 0; i < q; i++ {
			var e int
			fmt.Fscan(in, &e)
			v := (e + last) % n
			diff1 := trieA.query(mask ^ v)
			if diff1 > ans {
				ans = diff1
			}
			diff2 := trieB.query(v)
			if diff2 > ans {
				ans = diff2
			}
			trieA.add(v)
			trieB.add(mask ^ v)
			last = ans
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
