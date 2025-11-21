package main

import (
	"bufio"
	"fmt"
	"os"
)

type interval struct {
	l, r int
}

func buildPermutation(n int, fixed [][2]int) []int {
	perm := make([]int, n)
	for i := range perm {
		perm[i] = -1
	}
	used := make([]bool, n)
	for _, fv := range fixed {
		pos := fv[0] - 1
		val := fv[1]
		perm[pos] = val
		used[val] = true
	}
	nextVal := 0
	for i := 0; i < n; i++ {
		if perm[i] != -1 {
			continue
		}
		for nextVal < n && used[nextVal] {
			nextVal++
		}
		perm[i] = nextVal
		used[nextVal] = true
	}
	return perm
}

func printPermutation(out *bufio.Writer, perm []int) {
	for i, v := range perm {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, v)
	}
	out.WriteByte('\n')
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
		segs := make([]interval, m)
		globalL, globalR := 1, n
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &segs[i].l, &segs[i].r)
			if segs[i].l > globalL {
				globalL = segs[i].l
			}
			if segs[i].r < globalR {
				globalR = segs[i].r
			}
		}

		if globalL <= globalR {
			perm := buildPermutation(n, [][2]int{{globalL, 0}})
			printPermutation(out, perm)
			continue
		}

		cover := make([]bool, n)
		leftMax := make([]int, n)
		rightMin := make([]int, n)
		for pos := 1; pos <= n; pos++ {
			cov := false
			lmax := 1
			rmin := n
			for _, seg := range segs {
				if seg.l <= pos && pos <= seg.r {
					cov = true
					if seg.l > lmax {
						lmax = seg.l
					}
					if seg.r < rmin {
						rmin = seg.r
					}
				}
			}
			idx := pos - 1
			cover[idx] = cov
			leftMax[idx] = lmax
			rightMin[idx] = rmin
		}

		found := false
		pos0, pos1 := 0, 0
		for pos := 1; pos <= n; pos++ {
			idx := pos - 1
			if !cover[idx] {
				pos0 = pos
				if pos != 1 {
					pos1 = 1
				} else {
					pos1 = 2
				}
				found = true
				break
			}
			if rightMin[idx]-leftMax[idx]+1 >= 2 {
				pos0 = pos
				if leftMax[idx] < pos {
					pos1 = leftMax[idx]
				} else {
					pos1 = rightMin[idx]
				}
				found = true
				break
			}
		}

		if found {
			perm := buildPermutation(n, [][2]int{{pos0, 0}, {pos1, 1}})
			printPermutation(out, perm)
			continue
		}

		assign := [][2]int{{1, 0}, {n, 1}}
		mid := 2
		if mid == n {
			mid = n - 1
		}
		assign = append(assign, [2]int{mid, 2})
		perm := buildPermutation(n, assign)
		printPermutation(out, perm)
	}
}
