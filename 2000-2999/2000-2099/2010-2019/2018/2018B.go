package main

import (
	"bufio"
	"fmt"
	"os"
)

func extreme(a []int, keepSmall bool) int {
	n := len(a)
	L, R := 0, n-1
	for cur := n; cur >= 2; cur-- {
		canL := a[L] >= cur
		canR := a[R] >= cur
		if !canL && !canR {
			return -1
		}
		if canL && canR {
			if keepSmall {
				R--
			} else {
				L++
			}
		} else if canL {
			L++
		} else {
			R--
		}
	}
	if a[L] < 1 {
		return -1
	}
	return L + 1
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		minPos := extreme(a, true)
		if minPos == -1 {
			fmt.Fprintln(out, 0)
			continue
		}
		maxPos := extreme(a, false)
		if maxPos == -1 {
			fmt.Fprintln(out, 0)
			continue
		}
		ans := maxPos - minPos + 1
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(out, ans)
	}
}
