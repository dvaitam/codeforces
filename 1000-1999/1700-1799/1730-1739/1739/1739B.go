package main

import (
	"bufio"
	"fmt"
	"os"
)

func add(dp map[int]int, tr map[int]int, v, prev, cnt int) {
	prevCount := dp[v]
	newCount := prevCount + cnt
	if newCount > 2 {
		newCount = 2
	}
	dp[v] = newCount
	if cnt == 1 {
		if prevCount == 0 {
			if _, ok := tr[v]; !ok {
				tr[v] = prev
			} else {
				delete(tr, v)
			}
		} else {
			delete(tr, v)
		}
	} else {
		delete(tr, v)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		d := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &d[i])
		}
		dp := map[int]int{d[0]: 1}
		trace := make([]map[int]int, n)
		trace[0] = map[int]int{d[0]: -1}
		for i := 1; i < n; i++ {
			dp2 := make(map[int]int)
			tr2 := make(map[int]int)
			for val, cnt := range dp {
				add(dp2, tr2, val+d[i], val, cnt)
				if d[i] > 0 && val >= d[i] {
					add(dp2, tr2, val-d[i], val, cnt)
				}
			}
			dp = dp2
			trace[i] = tr2
		}
		total := 0
		finalVal := 0
		for val, cnt := range dp {
			if total+cnt > 1 {
				total = 2
				break
			}
			total += cnt
			if cnt == 1 {
				finalVal = val
			}
		}
		if total != 1 {
			fmt.Fprintln(out, -1)
		} else {
			ans := make([]int, n)
			cur := finalVal
			for i := n - 1; i >= 0; i-- {
				ans[i] = cur
				cur = trace[i][cur]
			}
			for i := 0; i < n; i++ {
				if i > 0 {
					out.WriteByte(' ')
				}
				fmt.Fprint(out, ans[i])
			}
			out.WriteByte('\n')
		}
	}
}
