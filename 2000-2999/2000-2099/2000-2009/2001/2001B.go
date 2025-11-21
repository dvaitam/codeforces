package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n int) []int {
	if n == 1 {
		return []int{1}
	}
	if n == 2 {
		return nil
	}

	ans := make([]int, n)
	left := 0
	right := n - 1
	cur := 1
	for left <= right {
		ans[left] = cur
		left++
		cur++
		if left > right {
			break
		}
		ans[right] = cur
		right--
		cur++
	}
	return ans
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
		res := solve(n)
		if res == nil {
			fmt.Fprintln(out, -1)
		} else {
			for i, v := range res {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
