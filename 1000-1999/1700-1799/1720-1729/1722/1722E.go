package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxHW = 1000

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		// accumulate area by (h,w)
		rect := make([][]int64, maxHW+1)
		for i := range rect {
			rect[i] = make([]int64, maxHW+1)
		}
		for i := 0; i < n; i++ {
			var h, w int
			fmt.Fscan(reader, &h, &w)
			if h <= maxHW && w <= maxHW {
				rect[h][w] += int64(h * w)
			}
		}
		// build prefix sums
		prefix := make([][]int64, maxHW+1)
		for i := range prefix {
			prefix[i] = make([]int64, maxHW+1)
		}
		for i := 1; i <= maxHW; i++ {
			var rowSum int64
			for j := 1; j <= maxHW; j++ {
				rowSum += rect[i][j]
				prefix[i][j] = prefix[i-1][j] + rowSum
			}
		}
		for ; q > 0; q-- {
			var hs, ws, hb, wb int
			fmt.Fscan(reader, &hs, &ws, &hb, &wb)
			if hs+1 > hb-1 || ws+1 > wb-1 {
				fmt.Fprintln(writer, 0)
				continue
			}
			h2, w2 := hb-1, wb-1
			ans := prefix[h2][w2] - prefix[hs][w2] - prefix[h2][ws] + prefix[hs][ws]
			fmt.Fprintln(writer, ans)
		}
	}
}
