package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, delta int) {
	i++
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	i++
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		total := n * m
		arr := make([]int, total)
		for i := 0; i < total; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		pairs := make([]pair, total)
		for i := 0; i < total; i++ {
			pairs[i] = pair{val: arr[i], idx: i}
		}

		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].val == pairs[j].val {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].val < pairs[j].val
		})

		seat := make([]int, total)
		i := 0
		for i < total {
			j := i
			for j < total && pairs[j].val == pairs[i].val {
				j++
			}
			pos := i
			for k := j - 1; k >= i; k-- {
				seat[pairs[k].idx] = pos
				pos++
			}
			i = j
		}

		bit := NewBIT(total)
		ans := 0
		for idx := 0; idx < total; idx++ {
			p := seat[idx]
			ans += bit.Sum(p)
			bit.Add(p, 1)
		}
		fmt.Fprintln(writer, ans)
	}
}
