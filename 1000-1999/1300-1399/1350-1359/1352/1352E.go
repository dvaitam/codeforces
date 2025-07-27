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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		// prefix sums
		prefix := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + a[i]
		}

		seen := make([]bool, n+1)
		// mark sums of all subarrays with length >= 2
		for l := 0; l < n; l++ {
			for r := l + 2; r <= n; r++ {
				sum := prefix[r] - prefix[l]
				if sum > n {
					break
				}
				seen[sum] = true
			}
		}

		ans := 0
		for i := 0; i < n; i++ {
			if seen[a[i]] {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
