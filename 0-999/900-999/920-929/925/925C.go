package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	nums := make([]uint64, n)
	used := make([]bool, 60)
	for i := 0; i < n; i++ {
		var x uint64
		fmt.Fscan(in, &x)
		nums[i] = x
		bit := 63 - bits.LeadingZeros64(x)
		if used[bit] {
			fmt.Println("No")
			return
		}
		used[bit] = true
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	fmt.Println("Yes")
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(nums[i])
	}
	fmt.Println()
}
