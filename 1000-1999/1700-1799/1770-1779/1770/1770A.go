package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		nums := make([]int, n+m)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &nums[i])
		}
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &nums[n+i])
		}
		sort.Slice(nums, func(i, j int) bool { return nums[i] > nums[j] })
		sum := 0
		for i := 0; i < n && i < len(nums); i++ {
			sum += nums[i]
		}
		fmt.Fprintln(writer, sum)
	}
}
