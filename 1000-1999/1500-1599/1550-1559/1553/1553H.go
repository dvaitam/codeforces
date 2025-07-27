package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int = 1 << 60

func solve(nums []int, b int) ([]int, []int, []int) {
	size := 1 << b
	ans := make([]int, size)
	mn := make([]int, size)
	mx := make([]int, size)
	if b == 0 {
		ans[0] = inf
		if len(nums) > 0 {
			mn[0] = 0
			mx[0] = 0
		} else {
			mn[0] = inf
			mx[0] = -inf
		}
		return ans, mn, mx
	}
	mid := 1 << (b - 1)
	leftNums := make([]int, 0)
	rightNums := make([]int, 0)
	for _, v := range nums {
		if v&mid == 0 {
			leftNums = append(leftNums, v)
		} else {
			rightNums = append(rightNums, v-mid)
		}
	}
	ansL, mnL, mxL := solve(leftNums, b-1)
	ansR, mnR, mxR := solve(rightNums, b-1)
	for mask := 0; mask < size; mask++ {
		hi := mask >> (b - 1)
		low := mask & (mid - 1)
		if hi == 0 {
			cross := inf
			if len(leftNums) > 0 && len(rightNums) > 0 {
				cross = (mnR[low] + mid) - mxL[low]
			}
			a := ansL[low]
			if ansR[low] < a {
				a = ansR[low]
			}
			if cross < a {
				a = cross
			}
			ans[mask] = a
			minVal := inf
			if len(leftNums) > 0 && mnL[low] < minVal {
				minVal = mnL[low]
			}
			if len(rightNums) > 0 && mnR[low]+mid < minVal {
				minVal = mnR[low] + mid
			}
			mn[mask] = minVal
			maxVal := -inf
			if len(leftNums) > 0 && mxL[low] > maxVal {
				maxVal = mxL[low]
			}
			if len(rightNums) > 0 && mxR[low]+mid > maxVal {
				maxVal = mxR[low] + mid
			}
			mx[mask] = maxVal
		} else {
			cross := inf
			if len(leftNums) > 0 && len(rightNums) > 0 {
				cross = (mnL[low] + mid) - mxR[low]
			}
			a := ansL[low]
			if ansR[low] < a {
				a = ansR[low]
			}
			if cross < a {
				a = cross
			}
			ans[mask] = a
			minVal := inf
			if len(rightNums) > 0 && mnR[low] < minVal {
				minVal = mnR[low]
			}
			if len(leftNums) > 0 && mnL[low]+mid < minVal {
				minVal = mnL[low] + mid
			}
			mn[mask] = minVal
			maxVal := -inf
			if len(rightNums) > 0 && mxR[low] > maxVal {
				maxVal = mxR[low]
			}
			if len(leftNums) > 0 && mxL[low]+mid > maxVal {
				maxVal = mxL[low] + mid
			}
			mx[mask] = maxVal
		}
	}
	return ans, mn, mx
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &nums[i])
	}
	ans, _, _ := solve(nums, k)
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
}
