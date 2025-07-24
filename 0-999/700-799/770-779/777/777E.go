package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// BIT implements a Fenwick tree supporting prefix maximum queries.
type BIT struct {
	n    int
	tree []int64
}

// NewBIT creates a BIT of size n (1-based indexing).
func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

// Update sets tree[i] = max(tree[i], val) for all relevant nodes.
func (b *BIT) Update(i int, val int64) {
	for i <= b.n {
		if val > b.tree[i] {
			b.tree[i] = val
		}
		i += i & -i
	}
}

// Query returns the maximum value in prefix [1..i].
func (b *BIT) Query(i int) int64 {
	var res int64
	for i > 0 {
		if b.tree[i] > res {
			res = b.tree[i]
		}
		i -= i & -i
	}
	return res
}

// Ring represents a ring with inner radius a, outer radius b and height h.
type Ring struct {
	a, b int
	h    int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	rings := make([]Ring, n)
	bVals := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &rings[i].a, &rings[i].b, &rings[i].h)
		bVals[i] = rings[i].b
	}

	// Sort rings by descending outer radius, then by descending inner radius.
	sort.Slice(rings, func(i, j int) bool {
		if rings[i].b == rings[j].b {
			return rings[i].a > rings[j].a
		}
		return rings[i].b > rings[j].b
	})

	// Prepare unique sorted list of outer radii in descending order.
	sort.Ints(bVals)
	uniq := make([]int, 0, len(bVals))
	for i := len(bVals) - 1; i >= 0; i-- {
		if i == len(bVals)-1 || bVals[i] != bVals[i+1] {
			uniq = append(uniq, bVals[i])
		}
	}
	idxMap := make(map[int]int)
	for i, v := range uniq {
		idxMap[v] = i
	}

	bit := NewBIT(len(uniq))
	var ans int64
	for _, r := range rings {
		// Find prefix length of radii strictly greater than r.a.
		pos := sort.Search(len(uniq), func(i int) bool { return uniq[i] <= r.a })
		var best int64
		if pos > 0 {
			best = bit.Query(pos)
		}
		cur := int64(r.h) + best
		bit.Update(idxMap[r.b]+1, cur)
		if cur > ans {
			ans = cur
		}
	}

	fmt.Fprintln(writer, ans)
}
