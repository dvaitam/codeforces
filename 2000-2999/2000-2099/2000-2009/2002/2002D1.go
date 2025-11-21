package main

import (
	"bufio"
	"fmt"
	"os"
)

type segTree struct {
	size int
	min  []int
	max  []int
	inf  int
}

func newSegTree(n int) *segTree {
	size := 1
	for size < n {
		size <<= 1
	}
	inf := int(1e9)
	minArr := make([]int, size*2)
	maxArr := make([]int, size*2)
	for i := range minArr {
		minArr[i] = inf
	}
	return &segTree{
		size: size,
		min:  minArr,
		max:  maxArr,
		inf:  inf,
	}
}

func (st *segTree) update(pos, val int) {
	idx := pos + st.size
	st.min[idx] = val
	st.max[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		left := idx << 1
		right := left | 1
		if st.min[left] < st.min[right] {
			st.min[idx] = st.min[left]
		} else {
			st.min[idx] = st.min[right]
		}
		if st.max[left] > st.max[right] {
			st.max[idx] = st.max[left]
		} else {
			st.max[idx] = st.max[right]
		}
	}
}

func (st *segTree) queryMin(l, r int) int {
	if l > r {
		return st.inf
	}
	res := st.inf
	l += st.size
	r += st.size
	for l <= r {
		if l&1 == 1 {
			if st.min[l] < res {
				res = st.min[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.min[r] < res {
				res = st.min[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func (st *segTree) queryMax(l, r int) int {
	if l > r {
		return 0
	}
	res := 0
	l += st.size
	r += st.size
	for l <= r {
		if l&1 == 1 {
			if st.max[l] > res {
				res = st.max[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.max[r] > res {
				res = st.max[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(in, &n, &q)
		for i := 2; i <= n; i++ {
			var tmp int
			fmt.Fscan(in, &tmp)
		}

		p := make([]int, n+1)
		pos := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &p[i])
			pos[p[i]] = i
		}

		tin := make([]int, n+1)
		tout := make([]int, n+1)
		timer := 0
		var dfs func(int)
		dfs = func(v int) {
			tin[v] = timer
			timer++
			left := v * 2
			if left <= n {
				dfs(left)
				if left+1 <= n {
					dfs(left + 1)
				}
			}
			tout[v] = timer - 1
		}
		dfs(1)

		seg := newSegTree(n)
		for v := 1; v <= n; v++ {
			seg.update(tin[v], pos[v])
		}

		isInternal := func(v int) bool {
			return v*2 <= n
		}

		ok := make([]bool, n+1)
		var compute func(int) bool
		compute = func(v int) bool {
			if !isInternal(v) {
				return true
			}
			left := v * 2
			right := left + 1
			minL := seg.queryMin(tin[left], tout[left])
			minR := seg.queryMin(tin[right], tout[right])
			ppos := pos[v]
			if !(ppos < minL && ppos < minR) {
				return false
			}
			maxL := seg.queryMax(tin[left], tout[left])
			maxR := seg.queryMax(tin[right], tout[right])
			if maxL < minR || maxR < minL {
				return true
			}
			return false
		}

		bad := 0
		for v := 1; v <= n; v++ {
			ok[v] = compute(v)
			if !ok[v] {
				bad++
			}
		}

		mark := make([]bool, n+1)
		updates := make([]int, 0, 64)
		addNode := func(v int) {
			for v > 0 {
				if !mark[v] {
					mark[v] = true
					updates = append(updates, v)
				}
				v >>= 1
			}
		}

		for i := 0; i < q; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if x != y {
				vx := p[x]
				vy := p[y]
				p[x], p[y] = p[y], p[x]
				pos[vx], pos[vy] = y, x
				seg.update(tin[vx], pos[vx])
				seg.update(tin[vy], pos[vy])

				addNode(vx)
				addNode(vy)

				for _, v := range updates {
					old := ok[v]
					now := compute(v)
					ok[v] = now
					if old && !now {
						bad++
					} else if !old && now {
						bad--
					}
					mark[v] = false
				}
				updates = updates[:0]
			}

			if bad == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

