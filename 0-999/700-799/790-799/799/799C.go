package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fountain struct {
	b int
	p int
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) update(i, val int) {
	for ; i <= b.n; i += i & -i {
		if b.tree[i] < val {
			b.tree[i] = val
		}
	}
}

func (b *BIT) query(i int) int {
	if i <= 0 {
		return 0
	}
	if i > b.n {
		i = b.n
	}
	res := 0
	for ; i > 0; i -= i & -i {
		if b.tree[i] > res {
			res = b.tree[i]
		}
	}
	return res
}

func bestPair(arr []fountain, budget int) int {
	if len(arr) < 2 {
		return 0
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].p < arr[j].p })
	bit := NewBIT(budget)
	best := 0
	for _, f := range arr {
		if f.p > budget {
			break
		}
		remain := budget - f.p
		if v := bit.query(remain); v+f.b > best {
			best = v + f.b
		}
		bit.update(f.p, f.b)
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, c, d int
	if _, err := fmt.Fscan(in, &n, &c, &d); err != nil {
		return
	}
	coins := make([]fountain, 0, n)
	diams := make([]fountain, 0, n)
	bestCoin := 0
	bestDiam := 0
	for i := 0; i < n; i++ {
		var b, p int
		var t string
		fmt.Fscan(in, &b, &p, &t)
		if t == "C" {
			coins = append(coins, fountain{b, p})
			if p <= c && b > bestCoin {
				bestCoin = b
			}
		} else {
			diams = append(diams, fountain{b, p})
			if p <= d && b > bestDiam {
				bestDiam = b
			}
		}
	}

	ans := 0
	if bestCoin > 0 && bestDiam > 0 {
		ans = bestCoin + bestDiam
	}
	if v := bestPair(coins, c); v > ans {
		ans = v
	}
	if v := bestPair(diams, d); v > ans {
		ans = v
	}
	fmt.Fprintln(out, ans)
}
