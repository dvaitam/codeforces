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
		var s string
		fmt.Fscan(reader, &s)

		zeros := make([]int, 0, n)
		for i, ch := range s {
			if ch == '0' {
				zeros = append(zeros, i+1) // 1-indexed positions
			}
		}
		k := len(zeros)
		prefix := make([]int, k+1)
		for i := 0; i < k; i++ {
			prefix[i+1] = prefix[i] + zeros[i]
		}

		for i := 1; i <= n; i++ {
			if i > k {
				if i > 1 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, -1)
				continue
			}
			targetSum := (n-i+1)*i + i*(i-1)/2
			sumZeros := prefix[k] - prefix[k-i]
			cost := targetSum - sumZeros
			if i > 1 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, cost)
		}
		writer.WriteByte('\n')
	}
}
