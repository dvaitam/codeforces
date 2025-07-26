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
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}

		// prefix[i][v] holds the count of indices <= i with value < v
		prefix := make([][]uint16, n+1)
		for i := range prefix {
			prefix[i] = make([]uint16, n+1)
		}
		for i := 1; i <= n; i++ {
			copy(prefix[i], prefix[i-1])
			pi := p[i-1]
			for v := pi + 1; v <= n; v++ {
				prefix[i][v]++
			}
		}

		// suffix[i][v] holds the count of indices >= i with value < v
		suffix := make([][]uint16, n+2)
		for i := range suffix {
			suffix[i] = make([]uint16, n+1)
		}
		for i := n; i >= 1; i-- {
			copy(suffix[i], suffix[i+1])
			pi := p[i-1]
			for v := pi + 1; v <= n; v++ {
				suffix[i][v]++
			}
		}

		var ans int64
		for b := 2; b <= n-1; b++ {
			for c := b + 1; c <= n-1; c++ {
				aCount := int(prefix[b-1][p[c-1]])
				dCount := int(suffix[c+1][p[b-1]])
				ans += int64(aCount * dCount)
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
