package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type BIT struct {
	bit  []int
	used []int
}

func (b *BIT) init(n int) {
	if len(b.bit) < n+2 {
		b.bit = make([]int, n+2)
	} else {
		for i := range b.bit {
			b.bit[i] = 0
		}
		b.bit = b.bit[:n+2]
	}
	b.used = b.used[:0]
}

func (b *BIT) clear() {
	for _, idx := range b.used {
		b.bit[idx] = 0
	}
	b.used = b.used[:0]
}

func (b *BIT) update(idx, val int) {
	for i := idx; i < len(b.bit); i += i & -i {
		if b.bit[i] < val {
			b.bit[i] = val
		}
		b.used = append(b.used, i)
	}
}

func (b *BIT) query(idx int) int {
	res := 0
	for i := idx; i > 0; i -= i & -i {
		if b.bit[i] > res {
			res = b.bit[i]
		}
	}
	return res
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	vals := make([]int, n)
	copy(vals, arr)
	sort.Ints(vals)
	uniq := make([]int, 0, n)
	for _, v := range vals {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	m := len(uniq)
	rank := make(map[int]int, m)
	for i, v := range uniq {
		rank[v] = i
	}
	pos := make([][]int, m)
	for i, v := range arr {
		pos[rank[v]] = append(pos[rank[v]], i+1) // use 1-based index
	}

	bitPrev := &BIT{}
	bitCurr := &BIT{}
	bitPrev.init(n)
	bitCurr.init(n)

	ans := 0

	for r := 0; r < m; r++ {
		bitCurr.clear()
		for _, idx := range pos[r] {
			best := 1
			if val := bitCurr.query(idx - 1); val+1 > best {
				best = val + 1
			}
			if val := bitPrev.query(idx - 1); val+1 > best {
				best = val + 1
			}
			bitCurr.update(idx, best)
			if best > ans {
				ans = best
			}
		}
		bitPrev.clear()
		bitPrev, bitCurr = bitCurr, bitPrev
	}

	fmt.Fprintln(writer, n-ans)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}
