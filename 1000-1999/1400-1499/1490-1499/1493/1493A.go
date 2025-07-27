package main

import (
	"bufio"
	"fmt"
	"os"
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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		nums := make([]int, 0)
		// numbers greater than k can always be taken
		for i := k + 1; i <= n; i++ {
			nums = append(nums, i)
		}
		// choose numbers in (k/2, k)
		for i := k/2 + 1; i < k; i++ {
			nums = append(nums, i)
		}
		fmt.Fprintln(writer, len(nums))
		for i, v := range nums {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
