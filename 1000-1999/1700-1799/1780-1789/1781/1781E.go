package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Rect struct {
	u, l, d, r int64
	idx        int
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func max3(a, b, c int64) int64 {
	if a < b {
		a = b
	}
	if a < c {
		a = c
	}
	return a
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		rects := make([]Rect, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &rects[i].u, &rects[i].l, &rects[i].d, &rects[i].r)
			rects[i].idx = i
		}
		sort.Slice(rects, func(i, j int) bool {
			if rects[i].l == rects[j].l {
				return rects[i].r < rects[j].r
			}
			return rects[i].l < rects[j].l
		})
		res := make([][4]int64, n)
		right := [3]int64{0, 0, 0}
		total := int64(0)
		for _, re := range rects {
			if re.u == re.d { // one row
				row := re.u
				start := maxInt64(re.l, right[row]+1)
				if start <= re.r {
					res[re.idx] = [4]int64{row, start, row, re.r}
					total += re.r - start + 1
					right[row] = re.r
				}
				continue
			}
			// try to place as two rows
			startBoth := max3(re.l, right[1]+1, right[2]+1)
			if startBoth <= re.r {
				res[re.idx] = [4]int64{1, startBoth, 2, re.r}
				total += (re.r - startBoth + 1) * 2
				right[1] = re.r
				right[2] = re.r
				continue
			}
			// try row1 only
			start1 := maxInt64(re.l, right[1]+1)
			start2 := maxInt64(re.l, right[2]+1)
			width1 := int64(-1)
			if start1 <= re.r {
				width1 = re.r - start1 + 1
			}
			width2 := int64(-1)
			if start2 <= re.r {
				width2 = re.r - start2 + 1
			}
			if width1 >= width2 && width1 > 0 {
				res[re.idx] = [4]int64{1, start1, 1, re.r}
				total += width1
				right[1] = re.r
			} else if width2 > 0 {
				res[re.idx] = [4]int64{2, start2, 2, re.r}
				total += width2
				right[2] = re.r
			}
		}
		fmt.Fprintln(out, total)
		for i := 0; i < n; i++ {
			ans := res[i]
			fmt.Fprintln(out, ans[0], ans[1], ans[2], ans[3])
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	solve(in, out)
}
