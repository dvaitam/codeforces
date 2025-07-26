package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

// sumSubarrayMin computes the sum of minimums of all subarrays of arr.
func sumSubarrayMin(arr []int64) int64 {
	n := len(arr)
	left := make([]int, n)
	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && arr[stack[len(stack)-1]] > arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			left[i] = i + 1
		} else {
			left[i] = i - stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	right := make([]int, n)
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && arr[stack[len(stack)-1]] >= arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			right[i] = n - i
		} else {
			right[i] = stack[len(stack)-1] - i
		}
		stack = append(stack, i)
	}
	var res int64
	for i := 0; i < n; i++ {
		res += arr[i] * int64(left[i]) * int64(right[i])
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		maxVal := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}
		diff := make([]int64, n)
		for i := 0; i < n; i++ {
			diff[i] = maxVal - a[i]
		}
		var sumInc int64
		for i := 0; i < n; i++ {
			sumInc += diff[i]
		}
		for shift := 0; shift < n; shift++ {
			inc := make([]int64, n)
			for i := 0; i < n; i++ {
				inc[i] = diff[(i+shift)%n]
			}
			ops := inc[0]
			for i := 1; i < n; i++ {
				if inc[i] > inc[i-1] {
					ops += inc[i] - inc[i-1]
				}
			}
			smin := sumSubarrayMin(inc)
			cost := (2*smin - sumInc) % MOD
			if cost < 0 {
				cost += MOD
			}
			fmt.Fprintf(writer, "%d %d\n", ops, cost)
		}
	}
}
