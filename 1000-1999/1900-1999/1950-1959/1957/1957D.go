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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// prefix xor
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			p[i] = p[i-1] ^ a[i-1]
		}
		const maxB = 30
		count := make([][]int, maxB)
		for b := 0; b < maxB; b++ {
			count[b] = make([]int, n+1)
		}
		for i := 1; i <= n; i++ {
			x := p[i]
			for b := 0; b < maxB; b++ {
				count[b][i] = count[b][i-1]
				if (x>>b)&1 == 1 {
					count[b][i]++
				}
			}
		}
		var ans int64
		for y := 1; y <= n; y++ {
			ay := a[y-1]
			k := 0
			for b := maxB - 1; b >= 0; b-- {
				if (ay>>b)&1 == 1 {
					k = b
					break
				}
			}
			leftOne := count[k][y-1]
			leftZero := y - leftOne
			rightOne := count[k][n] - count[k][y-1]
			rightZero := (n - y + 1) - rightOne
			ans += int64(leftZero*rightZero + leftOne*rightOne)
		}
		fmt.Fprintln(writer, ans)
	}
}
