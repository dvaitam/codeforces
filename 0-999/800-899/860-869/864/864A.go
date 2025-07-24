package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	counts := make(map[int]int)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		counts[x]++
	}

	if len(counts) != 2 {
		fmt.Fprintln(out, "NO")
		return
	}
	nums := make([]int, 0, 2)
	for k := range counts {
		nums = append(nums, k)
	}
	if counts[nums[0]] != counts[nums[1]] {
		fmt.Fprintln(out, "NO")
		return
	}
	sort.Ints(nums)
	fmt.Fprintln(out, "YES")
	fmt.Fprintln(out, nums[0], nums[1])
}
