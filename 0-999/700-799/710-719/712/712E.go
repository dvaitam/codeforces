package main

import (
	"bufio"
	"fmt"
	"os"
)

// node represents the transition probabilities for an interval
type node struct {
	r1, r2 float64
}

// merge combines two nodes according to the problem's formula
func merge(a, b node) node {
	if a.r1 < 0 {
		return b
	}
	if b.r1 < 0 {
		return a
	}
	r1 := a.r1 * b.r1 / (1 - a.r2*(1-b.r1))
	r2 := b.r2 + ((1 - b.r2) * a.r2 * b.r1 / (1 - (1-b.r1)*a.r2))
	return node{r1, r2}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	size := n
	t := make([]node, 2*size)

	// Read initial probabilities
	for i := 1; i <= n; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		p := float64(a) / float64(b)
		t[size+i-1] = node{p, p}
	}

	// Build segment tree
	for i := size - 1; i >= 1; i-- {
		t[i] = merge(t[2*i], t[2*i+1])
	}

	// Process queries
	for qi := 0; qi < q; qi++ {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var idx, a, b int
			fmt.Fscan(reader, &idx, &a, &b)
			p := float64(a) / float64(b)
			pos := size + idx - 1
			t[pos] = node{p, p}
			for pos >>= 1; pos > 0; pos >>= 1 {
				t[pos] = merge(t[2*pos], t[2*pos+1])
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l = l + size - 1
			r = r + size
			left := node{-1, 0}
			right := node{-1, 0}
			for l < r {
				if l&1 == 1 {
					left = merge(left, t[l])
					l++
				}
				if r&1 == 1 {
					r--
					right = merge(t[r], right)
				}
				l >>= 1
				r >>= 1
			}
			res := merge(left, right)
			fmt.Fprintf(writer, "%.9f\n", res.r1)
		}
	}
}
