package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sort.Ints(nums)

	var sum int
	for i := 0; i < n; i += 2 {
		sum += nums[i+1] - nums[i]
	}

	fmt.Println(sum)
}
