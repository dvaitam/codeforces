package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxReachableSum(arr []int) int {
	n := len(arr)
	best := -1
	for _, val := range arr {
		if val > best {
			best = val
		}
	}
	for i := 0; i < n; i++ {
		if arr[i] == 1 {
			best++
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, s int
		fmt.Fscan(in, &n, &s)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		cnt := map[int]int{}
		for _, v := range arr {
			cnt[v]++
		}
		arranged := make([]int, 0, n)
		for i := 0; i < cnt[0]; i++ {
			arranged = append(arranged, 0)
		}
		arranged = append(arranged, 1)
		arranged = append(arranged, 2)
		for len(arranged) < n {
			arranged = append(arranged, 0)
		}

		maxSum := 0
		cur := arranged[0]
		for i := 0; i < n-1; i++ {
			cur += arranged[i+1]
			if cur > maxSum {
				maxSum = cur
			}
		}

		if s > maxSum {
		}
		fmt.Fprintln(out, "-1")
	}
}
