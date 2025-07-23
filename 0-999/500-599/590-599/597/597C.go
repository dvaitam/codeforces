package main

import (
	"bufio"
	"fmt"
	"os"
)

// BIT is a Fenwick tree for prefix sums of int64 values.
type BIT struct {
	n    int
	tree []int64
}

// newBIT creates a BIT with size n.
func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+2)}
}

// add adds val at position idx (1-based).
func (b *BIT) add(idx int, val int64) {
	for idx <= b.n {
		b.tree[idx] += val
		idx += idx & -idx
	}
}

// sum returns prefix sum of [1..idx].
func (b *BIT) sum(idx int) int64 {
	var s int64
	for idx > 0 {
		s += b.tree[idx]
		idx -= idx & -idx
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	if k == 0 {
		fmt.Fprintln(writer, n)
		return
	}
	bits := make([]*BIT, k+2)
	for i := range bits {
		bits[i] = newBIT(n + 2)
	}
	for _, x := range arr {
		bits[1].add(x, 1)
		for j := k + 1; j >= 2; j-- {
			cnt := bits[j-1].sum(x - 1)
			if cnt != 0 {
				bits[j].add(x, cnt)
			}
		}
	}
	ans := bits[k+1].sum(n + 1)
	fmt.Fprintln(writer, ans)
}
