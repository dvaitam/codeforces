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

	const limit = 1 << 18
	squares := make([]int, 0, 512)
	for i := 0; i*i < limit; i++ {
		squares = append(squares, i*i)
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		freq := make([]int64, limit)
		prefix := 0
		freq[0] = 1
		var bad int64

		for _, v := range arr {
			prefix ^= v
			for _, s := range squares {
				x := prefix ^ s
				if x < limit {
					bad += freq[x]
				}
			}
			freq[prefix]++
		}

		total := int64(n) * int64(n+1) / 2
		ans := total - bad
		fmt.Fprintln(out, ans)
	}
}
