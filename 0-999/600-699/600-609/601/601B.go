package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumMax(arr []int64) int64 {
	n := len(arr)
	if n == 0 {
		return 0
	}
	left := make([]int64, n)
	right := make([]int64, n)
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		for len(stack) > 0 && arr[stack[len(stack)-1]] < arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			left[i] = int64(i + 1)
		} else {
			left[i] = int64(i - stack[len(stack)-1])
		}
		stack = append(stack, i)
	}
	stack = stack[:0]
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && arr[stack[len(stack)-1]] <= arr[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			right[i] = int64(n - i)
		} else {
			right[i] = int64(stack[len(stack)-1] - i)
		}
		stack = append(stack, i)
	}
	var res int64
	for i := 0; i < n; i++ {
		res += arr[i] * left[i] * right[i]
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	diff := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		v := a[i+1] - a[i]
		if v < 0 {
			v = -v
		}
		diff[i] = v
	}
	out := bufio.NewWriter(os.Stdout)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		if r-l <= 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		arr := diff[l-1 : r-1]
		ans := sumMax(arr)
		fmt.Fprintln(out, ans)
	}
	out.Flush()
}
