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
		var n, s, k int
		fmt.Fscan(reader, &n, &s, &k)
		closed := make(map[int]bool, k)
		for i := 0; i < k; i++ {
			var x int
			fmt.Fscan(reader, &x)
			closed[x] = true
		}
		// since k <= 1000, checking up to k+100 is enough
		ans := 0
		for {
			lower := s - ans
			upper := s + ans
			okLower := lower >= 1 && !closed[lower]
			okUpper := upper <= n && !closed[upper]
			if okLower || okUpper {
				fmt.Fprintln(writer, ans)
				break
			}
			ans++
		}
	}
}
