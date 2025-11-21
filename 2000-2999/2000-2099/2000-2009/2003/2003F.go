package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type fenwick struct {
	bit []int
}

func newFenwick(n int) fenwick {
	return fenwick{bit: make([]int, n+1)}
}

func (f *fenwick) reset() {
	for i := range f.bit {
		f.bit[i] = 0
	}
}

func (f *fenwick) update(idx, val int) {
	for idx < len(f.bit) {
		if val > f.bit[idx] {
			f.bit[idx] = val
		}
		idx += idx & -idx
	}
}

func (f *fenwick) query(idx int) int {
	res := 0
	for idx > 0 {
		if f.bit[idx] > res {
			res = f.bit[idx]
		}
		idx -= idx & -idx
	}
	return res
}

func triesFor(m int) int {
	switch m {
	case 1:
		return 1
	case 2:
		return 40
	case 3:
		return 60
	case 4:
		return 120
	default:
		return 256 // m == 5
	}
}

// runOnce calculates the best colorful subsequence sum for a fixed coloring of b values.
func runOnce(colors []int, trees []fenwick, masks [][]int, a, b, c []int, n, targetMask int) int {
	// clear Fenwick trees for all non-zero masks
	for i := 1; i < len(trees); i++ {
		trees[i].reset()
	}

	for idx := 0; idx < n; idx++ {
		ai := a[idx]
		ci := c[idx]
		col := colors[b[idx]]
		bit := 1 << col

		for _, mask := range masks[col] {
			var prev int
			if mask != 0 {
				prev = trees[mask].query(ai)
				if prev == 0 {
					continue
				}
			}

			nv := prev + ci
			nm := mask | bit
			trees[nm].update(ai, nv)
		}
	}

	return trees[targetMask].query(len(trees[0].bit) - 1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}

	a := make([]int, n)
	b := make([]int, n)
	c := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &b[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &c[i])
	}

	seen := make([]bool, n+1)
	uniq := 0
	for _, v := range b {
		if !seen[v] {
			seen[v] = true
			uniq++
		}
	}
	if uniq < m {
		fmt.Println(-1)
		return
	}

	limit := 1 << m
	trees := make([]fenwick, limit)
	for i := 0; i < limit; i++ {
		trees[i] = newFenwick(n)
	}

	masks := make([][]int, m)
	for col := 0; col < m; col++ {
		bit := 1 << col
		cur := make([]int, 0, limit)
		for mask := 0; mask < limit; mask++ {
			if mask&bit == 0 {
				cur = append(cur, mask)
			}
		}
		masks[col] = cur
	}

	colors := make([]int, n+1)
	targetMask := limit - 1

	ans := 0

	// deterministic first coloring to slightly reduce failure risk
	for i := 1; i <= n; i++ {
		colors[i] = i % m
	}
	if v := runOnce(colors, trees, masks, a, b, c, n, targetMask); v > ans {
		ans = v
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tries := triesFor(m)
	for t := 0; t < tries; t++ {
		for i := 1; i <= n; i++ {
			colors[i] = rnd.Intn(m)
		}
		if v := runOnce(colors, trees, masks, a, b, c, n, targetMask); v > ans {
			ans = v
		}
	}

	if ans == 0 {
		fmt.Println(-1)
	} else {
		fmt.Println(ans)
	}
}
