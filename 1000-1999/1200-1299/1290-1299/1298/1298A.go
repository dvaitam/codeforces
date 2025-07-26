package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	nums := make([]int, 4)
	for i := 0; i < 4; i++ {
		if _, err := fmt.Fscan(in, &nums[i]); err != nil {
			return
		}
	}
	sort.Ints(nums)
	total := nums[3]
	fmt.Printf("%d %d %d\n", total-nums[0], total-nums[1], total-nums[2])
}
