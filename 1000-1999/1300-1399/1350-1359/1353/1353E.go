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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		totalOnes := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				totalOnes++
			}
		}

		bestGain := 0
		for r := 0; r < k; r++ {
			cur := 0
			maxCur := 0
			for i := r; i < n; i += k {
				if s[i] == '1' {
					cur++
				} else {
					cur--
				}
				if cur < 0 {
					cur = 0
				}
				if cur > maxCur {
					maxCur = cur
				}
			}
			if maxCur > bestGain {
				bestGain = maxCur
			}
		}

		ans := totalOnes - bestGain
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
