package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type RowInfo struct {
	row   int
	left  int
	right int
}

func distCollect(start, L, R int, endLeft bool) int {
	if L > R {
		L, R = R, L
	}
	if L == R {
		if start > L {
			return start - L
		}
		return L - start
	}
	if endLeft {
		// end at left
		option1 := abs(start-R) + (R - L)
		option2 := abs(start-L) + 2*(R-L)
		if option1 < option2 {
			return option1
		}
		return option2
	}
	option1 := abs(start-L) + (R - L)
	option2 := abs(start-R) + 2*(R-L)
	if option1 < option2 {
		return option1
	}
	return option2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gatherCandidates(safe []int, xs ...int) []int {
	m := make(map[int]struct{})
	for _, x := range xs {
		idx := sort.SearchInts(safe, x)
		if idx < len(safe) {
			m[safe[idx]] = struct{}{}
		}
		if idx > 0 {
			m[safe[idx-1]] = struct{}{}
		}
	}
	res := make([]int, 0, len(m))
	for v := range m {
		res = append(res, v)
	}
	return res
}

func transition(prevCost, prevPos, prevRow int, cur RowInfo, safe []int) (int, int) {
	dRow := cur.row - prevRow
	cand := gatherCandidates(safe, prevPos, cur.left, cur.right)
	bestL := int(1 << 60)
	bestR := int(1 << 60)
	for _, s := range cand {
		base := prevCost + abs(prevPos-s) + dRow
		cL := base + distCollect(s, cur.left, cur.right, true)
		cR := base + distCollect(s, cur.left, cur.right, false)
		if cL < bestL {
			bestL = cL
		}
		if cR < bestR {
			bestR = cR
		}
	}
	return bestL, bestR
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k, q int
	if _, err := fmt.Fscan(in, &n, &m, &k, &q); err != nil {
		return
	}
	rowMap := make(map[int][2]int)
	for i := 0; i < k; i++ {
		var r, c int
		fmt.Fscan(in, &r, &c)
		if val, ok := rowMap[r]; ok {
			if c < val[0] || val[0] == 0 {
				val[0] = c
			}
			if c > val[1] {
				val[1] = c
			}
			rowMap[r] = val
		} else {
			rowMap[r] = [2]int{c, c}
		}
	}
	safeCols := make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &safeCols[i])
	}
	sort.Ints(safeCols)

	// build rows list
	rows := make([]RowInfo, 0, len(rowMap)+1)
	for r, lr := range rowMap {
		rows = append(rows, RowInfo{r, lr[0], lr[1]})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].row < rows[j].row })
	if len(rows) == 0 || rows[0].row != 1 {
		// insert starting row
		rows = append([]RowInfo{{1, 1, 1}}, rows...)
	} else {
		// ensure we start at column 1 but keep treasure positions
		// nothing to do
	}

	dpL := make([]int, len(rows))
	dpR := make([]int, len(rows))

	// base case
	if rows[0].left == 1 && rows[0].right == 1 && rowMap[1] == [2]int{} {
		dpL[0], dpR[0] = 0, 0
	} else {
		dpL[0] = distCollect(1, rows[0].left, rows[0].right, true)
		dpR[0] = distCollect(1, rows[0].left, rows[0].right, false)
	}

	for i := 1; i < len(rows); i++ {
		cur := rows[i]
		bestL1, bestR1 := transition(dpL[i-1], rows[i-1].left, rows[i-1].row, cur, safeCols)
		bestL2, bestR2 := transition(dpR[i-1], rows[i-1].right, rows[i-1].row, cur, safeCols)
		if bestL1 > bestL2 {
			bestL1 = bestL2
		}
		if bestR1 > bestR2 {
			bestR1 = bestR2
		}
		dpL[i] = bestL1
		dpR[i] = bestR1
	}

	res := dpL[len(rows)-1]
	if dpR[len(rows)-1] < res {
		res = dpR[len(rows)-1]
	}
	fmt.Fprintln(out, res)
}
