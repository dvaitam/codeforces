package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// BIT implements a Fenwick tree for sum queries
type BIT struct {
	n   int
	bit []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, bit: make([]int, n+1)}
}

func (b *BIT) Add(i, v int) {
	for x := i; x <= b.n; x += x & -x {
		b.bit[x] += v
	}
}

func (b *BIT) Sum(i int) int {
	s := 0
	for x := i; x > 0; x -= x & -x {
		s += b.bit[x]
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	// polygon edge directions flattened
	arr := make([]int, 0, 300000)
	polyStart := make([]int, n+1)
	polyEnd := make([]int, n+1)
	// map direction to id
	type Dir struct{ x, y int64 }
	dirMap := make(map[Dir]int)
	nextID := 1
	var cum int
	// read polygons
	for i := 1; i <= n; i++ {
		var k int
		fmt.Fscan(in, &k)
		xs := make([]int64, k)
		ys := make([]int64, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &xs[j], &ys[j])
		}
		polyStart[i] = cum + 1
		for j := 0; j < k; j++ {
			x1, y1 := xs[j], ys[j]
			x2, y2 := xs[(j+1)%k], ys[(j+1)%k]
			dx := x2 - x1
			dy := y2 - y1
			// normalize direction
			if dx == 0 {
				dy = 1
			} else if dy == 0 {
				dx = 1
			} else {
				if dx < 0 {
					dx, dy = -dx, -dy
				}
				// gcd
				a, b := dx, dy
				if a < 0 {
					a = -a
				}
				if b < 0 {
					b = -b
				}
				for b != 0 {
					a, b = b, a%b
				}
				g := a
				dx /= g
				dy /= g
			}
			d := Dir{dx, dy}
			id, ok := dirMap[d]
			if !ok {
				id = nextID
				dirMap[d] = id
				nextID++
			}
			arr = append(arr, id)
			cum++
		}
		polyEnd[i] = cum
	}
	// read queries
	var q int
	fmt.Fscan(in, &q)
	type Query struct{ L, R, idx int }
	qs := make([]Query, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		qs[i] = Query{polyStart[l], polyEnd[r], i}
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].R < qs[j].R })

	bit := NewBIT(len(arr))
	last := make([]int, nextID)
	ans := make([]int, q)
	curr := 0
	for _, qu := range qs {
		for curr < qu.R {
			// process position curr+1
			id := arr[curr]
			if last[id] != 0 {
				bit.Add(last[id], -1)
			}
			bit.Add(curr+1, 1)
			last[id] = curr + 1
			curr++
		}
		ans[qu.idx] = bit.Sum(qu.R) - bit.Sum(qu.L-1)
	}
	// output
	for i := 0; i < q; i++ {
		fmt.Fprintln(out, ans[i])
	}
}
