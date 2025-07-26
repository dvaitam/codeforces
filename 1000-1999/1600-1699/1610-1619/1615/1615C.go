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
		var a, b string
		fmt.Fscan(reader, &a, &b)

		c00, c01, c10, c11 := 0, 0, 0, 0
		for i := 0; i < n; i++ {
			aBit := a[i]
			bBit := b[i]
			switch {
			case aBit == '0' && bBit == '0':
				c00++
			case aBit == '0' && bBit == '1':
				c01++
			case aBit == '1' && bBit == '0':
				c10++
			case aBit == '1' && bBit == '1':
				c11++
			}
		}

		ans := n + 5
		if c01 == c10 {
			cand := c01 + c10
			if cand < ans {
				ans = cand
			}
		}
		if c11 == c00+1 {
			cand := c11 + c00
			if cand < ans {
				ans = cand
			}
		}
		if ans == n+5 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}
