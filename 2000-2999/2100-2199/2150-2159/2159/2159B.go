package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	l    int32
	span int32 // width is span + 1
}

type rowPair struct {
	u         int32
	height    int32
	intervals []interval
}

type pqNode struct {
	area    int
	pairIdx int32
	intIdx  int32
}

type minHeap []pqNode

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].area < h[j].area }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(pqNode)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

func transpose(grid [][]byte) [][]byte {
	n := len(grid)
	m := len(grid[0])
	res := make([][]byte, m)
	for i := 0; i < m; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = grid[j][i]
		}
		res[i] = row
	}
	return res
}

func buildRowOnes(grid [][]byte) [][]int32 {
	n := len(grid)
	rowOnes := make([][]int32, n)
	for i := 0; i < n; i++ {
		row := grid[i]
		tmp := make([]int32, 0, len(row))
		for j := 0; j < len(row); j++ {
			if row[j] == '1' {
				tmp = append(tmp, int32(j+1))
			}
		}
		rowOnes[i] = tmp
	}
	return rowOnes
}

func makePairs(rowOnes [][]int32) []rowPair {
	n := len(rowOnes)
	pairs := make([]rowPair, 0, n*(n-1)/2)
	for u := 0; u < n; u++ {
		if len(rowOnes[u]) == 0 {
			continue
		}
		for d := u + 1; d < n; d++ {
			if len(rowOnes[d]) == 0 {
				continue
			}
			i, j := 0, 0
			prev := int32(-1)
			capHint := len(rowOnes[u])
			if len(rowOnes[d]) < capHint {
				capHint = len(rowOnes[d])
			}
			intervals := make([]interval, 0, capHint)
			for i < len(rowOnes[u]) && j < len(rowOnes[d]) {
				if rowOnes[u][i] == rowOnes[d][j] {
					cur := rowOnes[u][i]
					if prev != -1 {
						intervals = append(intervals, interval{l: prev, span: cur - prev})
					}
					prev = cur
					i++
					j++
				} else if rowOnes[u][i] < rowOnes[d][j] {
					i++
				} else {
					j++
				}
			}
			if len(intervals) == 0 {
				continue
			}
			if len(intervals) > 1 {
				sort.Slice(intervals, func(a, b int) bool {
					if intervals[a].span == intervals[b].span {
						return intervals[a].l < intervals[b].l
					}
					return intervals[a].span < intervals[b].span
				})
			}
			pairs = append(pairs, rowPair{
				u:         int32(u),
				height:    int32(d - u + 1),
				intervals: intervals,
			})
		}
	}
	return pairs
}

func initParents(n, m int) [][]int32 {
	parents := make([][]int32, n)
	limit := int32(m + 1)
	for i := 0; i < n; i++ {
		arr := make([]int32, m+2)
		for j := int32(0); j <= limit+1; j++ {
			arr[j] = j
		}
		parents[i] = arr
	}
	return parents
}

func findNext(parent []int32, x int32) int32 {
	if x >= int32(len(parent)) {
		return int32(len(parent) - 1)
	}
	y := x
	for parent[y] != y {
		y = parent[y]
	}
	for parent[x] != x {
		nxt := parent[x]
		parent[x] = y
		x = nxt
	}
	return y
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = []byte(s)
		}

		origN, origM := n, m
		transposed := false
		if n > m {
			grid = transpose(grid)
			n, m = m, n
			transposed = true
		}

		rowOnes := buildRowOnes(grid)
		pairs := makePairs(rowOnes)

		parents := initParents(n, m)
		ans := make([][]int, n)
		for i := 0; i < n; i++ {
			ans[i] = make([]int, m)
		}

		h := make(minHeap, 0, len(pairs))
		for idx, p := range pairs {
			width := int(p.intervals[0].span + 1)
			area := width * int(p.height)
			heap.Push(&h, pqNode{area: area, pairIdx: int32(idx), intIdx: 0})
		}

		totalCells := n * m
		assigned := 0

		for h.Len() > 0 && assigned < totalCells {
			node := heap.Pop(&h).(pqNode)
			p := &pairs[node.pairIdx]
			inter := p.intervals[node.intIdx]
			l := inter.l
			r := inter.l + inter.span
			area := node.area
			rowStart := int(p.u)
			rowEnd := rowStart + int(p.height)
			for row := rowStart; row < rowEnd; row++ {
				parent := parents[row]
				col := findNext(parent, l)
				for col <= r {
					ans[row][col-1] = area
					assigned++
					parent[col] = findNext(parent, col+1)
					col = parent[col]
				}
			}
			nextIdx := node.intIdx + 1
			if int(nextIdx) < len(p.intervals) {
				nextInter := p.intervals[nextIdx]
				width := int(nextInter.span + 1)
				heap.Push(&h, pqNode{area: width * int(p.height), pairIdx: node.pairIdx, intIdx: nextIdx})
			}
		}

		if transposed {
			newAns := make([][]int, origN)
			for i := 0; i < origN; i++ {
				newAns[i] = make([]int, origM)
			}
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					newAns[j][i] = ans[i][j]
				}
			}
			ans = newAns
			n, m = origN, origM
		}

		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, ans[i][j])
			}
			fmt.Fprintln(out)
		}
	}
}
