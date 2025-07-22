package main

import (
	"bufio"
	"fmt"
	"os"
)

const MaxVal = 100000

// Node represents a node in binary trie
type Node struct {
	child [2]*Node
	min   int
}

func newNode() *Node {
	return &Node{min: 1 << 30}
}

// Trie structure storing numbers and minimal value per prefix
type Trie struct {
	root *Node
}

func (t *Trie) insert(x int) {
	if t.root == nil {
		t.root = newNode()
	}
	node := t.root
	if x < node.min {
		node.min = x
	}
	for i := 16; i >= 0; i-- {
		b := (x >> i) & 1
		if node.child[b] == nil {
			node.child[b] = newNode()
		}
		node = node.child[b]
		if x < node.min {
			node.min = x
		}
	}
}

func (t *Trie) query(x, limit int) int {
	if t == nil || t.root == nil || t.root.min > limit {
		return -1
	}
	node := t.root
	res := 0
	for i := 16; i >= 0; i-- {
		bit := (x >> i) & 1
		prefer := bit ^ 1
		if node.child[prefer] != nil && node.child[prefer].min <= limit {
			node = node.child[prefer]
			res |= prefer << i
		} else if node.child[bit] != nil && node.child[bit].min <= limit {
			node = node.child[bit]
			res |= bit << i
		} else {
			return -1
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}

	// Precompute divisors
	divisors := make([][]int, MaxVal+1)
	for i := 1; i <= MaxVal; i++ {
		for j := i; j <= MaxVal; j += i {
			divisors[j] = append(divisors[j], i)
		}
	}

	tries := make(map[int]*Trie)

	for ; q > 0; q-- {
		var t int
		if _, err := fmt.Fscan(reader, &t); err != nil {
			return
		}
		if t == 1 {
			var u int
			fmt.Fscan(reader, &u)
			for _, d := range divisors[u] {
				tr := tries[d]
				if tr == nil {
					tr = &Trie{}
					tries[d] = tr
				}
				tr.insert(u)
			}
		} else if t == 2 {
			var x, k, s int
			fmt.Fscan(reader, &x, &k, &s)
			if x%k != 0 {
				fmt.Fprintln(writer, -1)
				continue
			}
			limit := s - x
			if limit < 0 {
				fmt.Fprintln(writer, -1)
				continue
			}
			tr := tries[k]
			if tr == nil {
				fmt.Fprintln(writer, -1)
				continue
			}
			ans := tr.query(x, limit)
			fmt.Fprintln(writer, ans)
		}
	}
}
