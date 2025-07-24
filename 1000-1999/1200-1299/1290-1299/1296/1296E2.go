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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	dp := make([]int, 26)
	colors := make([]int, n)
	maxColor := 0

	for i := 0; i < n; i++ {
		ch := int(s[i] - 'a')
		best := 0
		for j := ch + 1; j < 26; j++ {
			if dp[j] > best {
				best = dp[j]
			}
		}
		col := best + 1
		colors[i] = col
		if dp[ch] < col {
			dp[ch] = col
		}
		if col > maxColor {
			maxColor = col
		}
	}

	fmt.Fprintln(out, maxColor)
	for i, c := range colors {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, c)
	}
	out.WriteByte('\n')
}
