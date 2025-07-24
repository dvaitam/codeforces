package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Trie structure for storing numbers in binary form
// Each node has two children indexed by bit value 0 or 1
// Stored as an array of pairs of child indices

type Trie struct {
	next [][2]int
}

func NewTrie() *Trie {
	t := &Trie{next: make([][2]int, 1)}
	t.next[0] = [2]int{-1, -1}
	return t
}

func (t *Trie) Insert(x int) {
	node := 0
	for b := 29; b >= 0; b-- {
		bit := (x >> b) & 1
		if t.next[node][bit] == -1 {
			t.next[node][bit] = len(t.next)
			t.next = append(t.next, [2]int{-1, -1})
		}
		node = t.next[node][bit]
	}
}

func (t *Trie) MinXor(x int) int {
	node := 0
	res := 0
	for b := 29; b >= 0; b-- {
		bit := (x >> b) & 1
		if t.next[node][bit] != -1 {
			node = t.next[node][bit]
		} else {
			res |= 1 << b
			node = t.next[node][1-bit]
		}
	}
	return res
}

func solve(arr []int, bit int) int64 {
	if len(arr) <= 1 || bit < 0 {
		return 0
	}
	mask := 1 << bit
	i := 0
	for i < len(arr) && arr[i]&mask == 0 {
		i++
	}
	left := arr[:i]
	right := arr[i:]
	if len(left) == 0 {
		return solve(right, bit-1)
	}
	if len(right) == 0 {
		return solve(left, bit-1)
	}
	trie := NewTrie()
	for _, v := range left {
		trie.Insert(v)
	}
	best := int64(1 << 60)
	for _, v := range right {
		cur := trie.MinXor(v)
		if int64(cur) < best {
			best = int64(cur)
		}
	}
	return best + solve(left, bit-1) + solve(right, bit-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	sort.Ints(arr)
	ans := solve(arr, 29)
	fmt.Println(ans)
}
