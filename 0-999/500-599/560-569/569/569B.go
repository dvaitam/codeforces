package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(reader, &nums[i])
	}
	// seenIndex[v] stores the last index of value v in nums
	seenIndex := make([]int, n+2)
	for i := range seenIndex {
		seenIndex[i] = -1
	}
	for i, v := range nums {
		if v >= 1 && v <= n {
			seenIndex[v] = i
		}
	}
	// stay[i] indicates if original nums[i] should be kept
	stay := make([]bool, n)
	for v := 1; v <= n; v++ {
		if idx := seenIndex[v]; idx >= 0 {
			stay[idx] = true
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	// find the first missing positive integer
	j := 1
	for j <= n+1 && seenIndex[j] >= 0 {
		j++
	}
	for i, v := range nums {
		if stay[i] {
			fmt.Fprint(writer, v, " ")
		} else {
			fmt.Fprint(writer, j, " ")
			j++
			for j <= n+1 && seenIndex[j] >= 0 {
				j++
			}
		}
	}
	writer.Flush()
}
