package main

import (
	"bufio"
	"fmt"
	"os"
)

func absDiffGap(minL, maxL, minR, maxR int) int {
	if maxL < minR {
		return minR - maxL
	}
	if maxR < minL {
		return minL - maxR
	}
	return 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		already := false
		for i := 0; i+1 < n; i++ {
			if abs(a[i]-a[i+1]) <= 1 {
				already = true
				break
			}
		}
		if already {
			fmt.Fprintln(out, 0)
			continue
		}

		ans := n + 5
		for s := 0; s+1 < n; s++ {
			sizeLeft := s + 1
			sizeRight := n - (s + 1)
			leftMin := make([]int, sizeLeft+1)
			leftMax := make([]int, sizeLeft+1)
			mn, mx := 0, 0
			for lenL := 1; lenL <= sizeLeft; lenL++ {
				idx := s - (lenL - 1)
				if lenL == 1 {
					mn, mx = a[idx], a[idx]
				} else {
					if a[idx] < mn {
						mn = a[idx]
					}
					if a[idx] > mx {
						mx = a[idx]
					}
				}
				leftMin[lenL] = mn
				leftMax[lenL] = mx
			}
			rightMin := make([]int, sizeRight+1)
			rightMax := make([]int, sizeRight+1)
			for lenR := 1; lenR <= sizeRight; lenR++ {
				idx := s + lenR
				if lenR == 1 {
					mn, mx = a[idx], a[idx]
				} else {
					if a[idx] < mn {
						mn = a[idx]
					}
					if a[idx] > mx {
						mx = a[idx]
					}
				}
				rightMin[lenR] = mn
				rightMax[lenR] = mx
			}
			for lenL := 1; lenL <= sizeLeft; lenL++ {
				low, high, best := 1, sizeRight, 0
				for low <= high {
					mid := (low + high) >> 1
					gap := absDiffGap(leftMin[lenL], leftMax[lenL], rightMin[mid], rightMax[mid])
					if gap <= 1 {
						best = mid
						high = mid - 1
					} else {
						low = mid + 1
					}
				}
				if best > 0 {
					cost := lenL + best - 2
					if cost < ans {
						ans = cost
					}
				}
			}
		}
		if ans > n {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
