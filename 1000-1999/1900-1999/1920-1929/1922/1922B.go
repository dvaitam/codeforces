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
		var n int
		fmt.Fscan(reader, &n)
		freq := make([]int64, n+1)
		maxA := 0
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(reader, &a)
			freq[a]++
			if a > maxA {
				maxA = a
			}
		}
		prefix := make([]int64, maxA+1)
		var sum int64
		for i := 0; i <= maxA; i++ {
			sum += freq[i]
			prefix[i] = sum
		}
		var ans int64
		for v := 0; v <= maxA; v++ {
			if freq[v] >= 2 {
				c2 := freq[v] * (freq[v] - 1) / 2
				if v > 0 {
					ans += c2 * prefix[v-1]
				}
				if freq[v] >= 3 {
					ans += freq[v] * (freq[v] - 1) * (freq[v] - 2) / 6
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
