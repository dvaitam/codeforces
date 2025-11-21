package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

type node struct {
	key         int
	priority    uint32
	add         int
	size        int
	left, right *node
}

var rng = rand.New(rand.NewSource(1))

func newNode(key int) *node {
	return &node{
		key:      key,
		priority: rng.Uint32(),
		size:     1,
	}
}

func size(nd *node) int {
	if nd == nil {
		return 0
	}
	return nd.size
}

func pull(nd *node) {
	if nd != nil {
		nd.size = 1 + size(nd.left) + size(nd.right)
	}
}

func apply(nd *node, delta int) {
	if nd != nil {
		nd.key += delta
		nd.add += delta
	}
}

func push(nd *node) {
	if nd != nil && nd.add != 0 {
		if nd.left != nil {
			apply(nd.left, nd.add)
		}
		if nd.right != nil {
			apply(nd.right, nd.add)
		}
		nd.add = 0
	}
}

func split(nd *node, key int) (*node, *node) {
	if nd == nil {
		return nil, nil
	}
	push(nd)
	if nd.key >= key {
		left, rest := split(nd.left, key)
		nd.left = rest
		pull(nd)
		return left, nd
	}
	rest, right := split(nd.right, key)
	nd.right = rest
	pull(nd)
	return nd, right
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority < b.priority {
		push(a)
		a.right = merge(a.right, b)
		pull(a)
		return a
	}
	push(b)
	b.left = merge(a, b.left)
	pull(b)
	return b
}

func contains(nd *node, key int) bool {
	for nd != nil {
		push(nd)
		if key < nd.key {
			nd = nd.left
		} else if key > nd.key {
			nd = nd.right
		} else {
			return true
		}
	}
	return false
}

func insert(nd *node, key int) *node {
	node := newNode(key)
	left, right := split(nd, key)
	return merge(merge(left, node), right)
}

func countLE(nd *node, key int) int {
	if nd == nil {
		return 0
	}
	push(nd)
	if key < nd.key {
		return countLE(nd.left, key)
	}
	return size(nd.left) + 1 + countLE(nd.right, key)
}

func dValue(d0 int, root *node, x int) int {
	if x == 0 {
		return d0
	}
	return d0 - countLE(root, x-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefix := make([]int, n+1)
		for i := 0; i < n; i++ {
			cur := prefix[i]
			if a[i] > cur {
				cur++
			} else if a[i] < cur {
				cur--
			}
			prefix[i+1] = cur
		}

		prefixMaxBefore := make([]int, n+1)
		best := prefix[0]
		for r := 1; r <= n; r++ {
			if prefix[r-1] > best {
				best = prefix[r-1]
			}
			prefixMaxBefore[r] = best
		}

		var root *node
		d0 := 0
		ans := 0
		for r := n; r >= 1; r-- {
			xq := prefixMaxBefore[r]
			val := xq + dValue(d0, root, xq)
			if val > ans {
				ans = val
			}

			ai := a[r-1]
			if !contains(root, 0) {
				d0++
			}

			_, root = split(root, 1)

			var lower *node
			lower, root = split(root, ai)
			if lower != nil {
				apply(lower, -1)
			}

			var mid *node
			mid, root = split(root, n-1)
			if mid != nil {
				apply(mid, 1)
			}

			root = merge(lower, mid)

			if ai-1 >= 0 && ai-1 <= n-1 && !contains(root, ai-1) {
				root = insert(root, ai-1)
			}
			if ai <= n-1 && !contains(root, ai) {
				root = insert(root, ai)
			}
		}

		fmt.Fprintln(out, ans)
	}
}
