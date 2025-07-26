package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		candies := make([]int, n)
		minVal := int(1<<31 - 1)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &candies[i])
			if candies[i] < minVal {
				minVal = candies[i]
			}
			sum += candies[i]
		}
		fmt.Fprintln(out, sum-minVal*n)
	}
}
