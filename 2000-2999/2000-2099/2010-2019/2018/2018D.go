package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const negInf int64 = -1 << 60

// Node keeps DP for a segment:
// dp[lSel][rSel][mask] where lSel/rSel mark whether the first/last position
// in the segment is selected, and mask has bits:
// 1 - a position with maximal value is selected inside the segment
// 2 - a position with current minimum value (the sweeping value) is selected
type Node struct {
	dp [2][2][4]int64
}

func makeForbidden() Node {
	var res Node
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 4; k++ {
				res.dp[i][j][k] = negInf
			}
		}
	}
	res.dp[0][0][0] = 0
	return res
}

func makeLeaf(isMax, isSpecial, allowed bool) Node {
	if !allowed {
		return makeForbidden()
	}
	var res Node
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 4; k++ {
				res.dp[i][j][k] = negInf
			}
		}
	}
	res.dp[0][0][0] = 0 // not selected
	mask := 0
	if isMax {
		mask |= 1
	}
	if isSpecial {
		mask |= 2
	}
	res.dp[1][1][mask] = 1 // selected
	return res
}

func merge(a, b Node) Node {
	var res Node
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for m := 0; m < 4; m++ {
				res.dp[i][j][m] = negInf
			}
		}
	}
	for la := 0; la < 2; la++ {
		for lb := 0; lb < 2; lb++ {
			for lm := 0; lm < 4; lm++ {
				if a.dp[la][lb][lm] == negInf {
					continue
				}
				for rf := 0; rf < 2; rf++ { // right first
					for rr := 0; rr < 2; rr++ { // right last
						if lb == 1 && rf == 1 {
							continue // adjacent selected, invalid
						}
						for rm := 0; rm < 4; rm++ {
							if b.dp[rf][rr][rm] == negInf {
								continue
							}
							nm := lm | rm
							val := a.dp[la][lb][lm] + b.dp[rf][rr][rm]
							if val > res.dp[la][rr][nm] {
								res.dp[la][rr][nm] = val
							}
						}
					}
				}
			}
		}
	}
	return res
}

type SegTree struct {
	n    int
	size int
	data []Node
	isMx []bool
}

func NewSegTree(values []int, maxVal int) *SegTree {
	n := len(values)
	size := 1
	for size < n {
		size <<= 1
	}
	data := make([]Node, size<<1)
	isMx := make([]bool, n)
	for i, v := range values {
		isMx[i] = v == maxVal
	}
	st := &SegTree{n: n, size: size, data: data, isMx: isMx}
	for i := 0; i < size; i++ {
		if i < n {
			st.data[size+i] = makeForbidden()
		} else {
			// empty leaves behave as forbidden
			st.data[size+i] = makeForbidden()
		}
	}
	for i := size - 1; i >= 1; i-- {
		st.data[i] = merge(st.data[i<<1], st.data[i<<1|1])
	}
	return st
}

// update leaf i (0-indexed) according to whether it is allowed and special at current sweep
func (st *SegTree) update(idx int, allowed, special bool) {
	p := st.size + idx
	st.data[p] = makeLeaf(st.isMx[idx], special, allowed)
	for p >>= 1; p > 0; p >>= 1 {
		st.data[p] = merge(st.data[p<<1], st.data[p<<1|1])
	}
}

func (st *SegTree) root() Node {
	return st.data[1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] > maxVal {
				maxVal = arr[i]
			}
		}

		pos := make(map[int][]int)
		for i, v := range arr {
			pos[v] = append(pos[v], i)
		}
		uniq := make([]int, 0, len(pos))
		for v := range pos {
			uniq = append(uniq, v)
		}
		sort.Ints(uniq)

		st := NewSegTree(arr, maxVal)

		// initially all nodes are forbidden (already set)
		ans := int64(0)
		// sweep values descending
		for i := len(uniq) - 1; i >= 0; i-- {
			v := uniq[i]
			indices := pos[v]

			// add as special
			for _, idx := range indices {
				st.update(idx, true, true)
			}

			root := st.root()
			best := int64(negInf)
			for a := 0; a < 2; a++ {
				for b := 0; b < 2; b++ {
					for mask := 0; mask < 4; mask++ {
						if mask&1 == 1 && mask&2 == 2 {
							if root.dp[a][b][mask] > best {
								best = root.dp[a][b][mask]
							}
						}
					}
				}
			}
			if best > negInf {
				score := int64(maxVal+v) + best
				if score > ans {
					ans = score
				}
			}

			// demote to normal (allowed but not special)
			for _, idx := range indices {
				st.update(idx, true, false)
			}
		}

		fmt.Fprintln(out, ans)
	}
}
