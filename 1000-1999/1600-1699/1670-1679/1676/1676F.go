package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func solve() {
	var n, k int
	fmt.Fscan(reader, &n, &k)
	freq := make(map[int]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		freq[x]++
	}
	nums := make([]int, 0)
	for v, c := range freq {
		if c >= k {
			nums = append(nums, v)
		}
	}
	if len(nums) == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	sort.Ints(nums)
	bestL, bestR := nums[0], nums[0]
	curL, curR := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1]+1 {
			curR = nums[i]
		} else {
			if curR-curL > bestR-bestL {
				bestL, bestR = curL, curR
			}
			curL, curR = nums[i], nums[i]
		}
	}
	if curR-curL > bestR-bestL {
		bestL, bestR = curL, curR
	}
	fmt.Fprintln(writer, bestL, bestR)
}

func main() {
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve()
	}
}
