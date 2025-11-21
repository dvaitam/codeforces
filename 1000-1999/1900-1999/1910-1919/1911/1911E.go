package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	nums := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	powers := make([]int64, 0, 61)
	for p := int64(1); p <= (int64(1) << 60); p <<= 1 {
		powers = append(powers, p)
	}

	counts := make(map[int64]int, n)
	for _, x := range nums {
		for _, p := range powers {
			if c := counts[p-x]; c > 0 {
				fmt.Println("YES")
				return
			}
		}
		counts[x]++
	}

	fmt.Println("NO")
}
