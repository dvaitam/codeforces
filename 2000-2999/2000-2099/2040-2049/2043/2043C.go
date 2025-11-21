package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Interval struct {
	l, r int
}

func subarrayInterval(arr []int) (int, int) {
	if len(arr) == 0 {
		return 0, 0
	}
	curMax, bestMax := arr[0], arr[0]
	curMin, bestMin := arr[0], arr[0]
	for i := 1; i < len(arr); i++ {
		v := arr[i]
		if curMax+v > v {
			curMax += v
		} else {
			curMax = v
		}
		if curMax > bestMax {
			bestMax = curMax
		}

		if curMin+v < v {
			curMin += v
		} else {
			curMin = v
		}
		if curMin < bestMin {
			bestMin = curMin
		}
	}
	return bestMin, bestMax
}

func suffixInterval(arr []int) (int, int) {
	sum := 0
	minV, maxV := 0, 0
	for i := len(arr) - 1; i >= 0; i-- {
		sum += arr[i]
		if sum < minV {
			minV = sum
		}
		if sum > maxV {
			maxV = sum
		}
	}
	return minV, maxV
}

func prefixInterval(arr []int) (int, int) {
	sum := 0
	minV, maxV := 0, 0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
		if sum < minV {
			minV = sum
		}
		if sum > maxV {
			maxV = sum
		}
	}
	return minV, maxV
}

func addInterval(intervals *[]Interval, l, r int) {
	if l > r {
		return
	}
	*intervals = append(*intervals, Interval{l, r})
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		pos := -1
		val := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] != -1 && arr[i] != 1 && pos == -1 {
				pos = i
				val = arr[i]
			}
		}

		intervals := make([]Interval, 0)
		if pos == -1 { // all +/-1
			if n > 0 {
				minV, maxV := subarrayInterval(arr)
				addInterval(&intervals, minV, maxV)
			}
		} else {
			left := arr[:pos]
			right := arr[pos+1:]
			if len(left) > 0 {
				minL, maxL := subarrayInterval(left)
				addInterval(&intervals, minL, maxL)
			}
			if len(right) > 0 {
				minR, maxR := subarrayInterval(right)
				addInterval(&intervals, minR, maxR)
			}
			lsMin, lsMax := suffixInterval(left)
			rsMin, rsMax := prefixInterval(right)
			addInterval(&intervals, val+lsMin+rsMin, val+lsMax+rsMax)
		}
		addInterval(&intervals, 0, 0)

		sort.Slice(intervals, func(i, j int) bool {
			if intervals[i].l == intervals[j].l {
				return intervals[i].r < intervals[j].r
			}
			return intervals[i].l < intervals[j].l
		})

		merged := make([]Interval, 0)
		for _, iv := range intervals {
			if len(merged) == 0 || iv.l > merged[len(merged)-1].r+1 {
				merged = append(merged, iv)
			} else {
				if iv.r > merged[len(merged)-1].r {
					merged[len(merged)-1].r = iv.r
				}
			}
		}

		total := 0
		for _, iv := range merged {
			total += iv.r - iv.l + 1
		}
		fmt.Fprintln(out, total)
		for idx, iv := range merged {
			for x := iv.l; x <= iv.r; x++ {
				if idx != 0 || x != merged[0].l {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, x)
			}
		}
		fmt.Fprintln(out)
	}
}
